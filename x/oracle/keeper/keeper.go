package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/peggy/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	stakeKeeper types.StakingKeeper
	codespace   sdk.CodespaceType

	// TODO: use this as param instead
	consensusNeeded float64 // The minimum % of stake needed to sign claims in order for consensus to occur
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, stakeKeeper types.StakingKeeper, codespace sdk.CodespaceType, consensusNeeded float64) Keeper {
	if consensusNeeded <= 0 || consensusNeeded > 1 {
		panic(types.ErrMinimumConsensusNeededInvalid(codespace).Error())
	}
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		stakeKeeper:     stakeKeeper,
		codespace:       codespace,
		consensusNeeded: consensusNeeded,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetProphecy gets the entire prophecy data struct for a given id
func (k Keeper) GetProphecy(ctx sdk.Context, id string) (types.Prophecy, sdk.Error) {
	if id == "" {
		return types.NewEmptyProphecy(), types.ErrInvalidIdentifier(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))
	if bz == nil {
		return types.NewEmptyProphecy(), types.ErrProphecyNotFound(k.Codespace())
	}
	var dbProphecy types.DBProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &dbProphecy)

	deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return types.NewEmptyProphecy(), types.ErrInternalDB(k.Codespace(), err)
	}
	return deSerializedProphecy, nil
}

// setProphecy saves a prophecy with an initial claim
func (k Keeper) setProphecy(ctx sdk.Context, prophecy types.Prophecy) sdk.Error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier(k.Codespace())
	}
	if len(prophecy.ClaimValidators) == 0 {
		return types.ErrNoClaims(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		return types.ErrInternalDB(k.Codespace(), err)
	}
	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(serializedProphecy))
	return nil
}

// ProcessClaim TODO: write description
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.Claim) (types.Status, sdk.Error) {
	activeValidator := k.checkActiveValidator(ctx, claim.ValidatorAddress)
	if !activeValidator {
		return types.Status{}, types.ErrInvalidValidator(k.Codespace())
	}
	if strings.TrimSpace(claim.Content) == "" {
		return types.Status{}, types.ErrInvalidClaim(k.Codespace())
	}
	prophecy, err := k.GetProphecy(ctx, claim.ID)
	if err != nil {
		if err.Code() != types.CodeProphecyNotFound {
			return types.Status{}, err
		}
		prophecy = types.NewProphecy(claim.ID)
	} else {
		if prophecy.Status.Text == types.SuccessStatusText || prophecy.Status.Text == types.FailedStatusText {
			return types.Status{}, types.ErrProphecyFinalized(k.Codespace())
		}
		if prophecy.ValidatorClaims[claim.ValidatorAddress.String()] != "" {
			return types.Status{}, types.ErrDuplicateMessage(k.Codespace())
		}
	}
	prophecy.AddClaim(claim.ValidatorAddress, claim.Content)
	prophecy = k.processCompletion(ctx, prophecy)
	err = k.setProphecy(ctx, prophecy)
	if err != nil {
		return types.Status{}, err
	}
	return prophecy.Status, nil
}

func (k Keeper) checkActiveValidator(ctx sdk.Context, validatorAddress sdk.ValAddress) bool {
	validator, found := k.stakeKeeper.GetValidator(ctx, validatorAddress)
	if !found {
		return false
	}
	bondStatus := validator.GetStatus()
	return bondStatus == sdk.Bonded
}

// processCompletion looks at a given prophecy an assesses whether the claim with the highest power on that prophecy has enough
// power to be considered successful, or alternatively, will never be able to become successful due to not enough validation power being
// left to push it over the threshold required for consensus.
func (k Keeper) processCompletion(ctx sdk.Context, prophecy types.Prophecy) types.Prophecy {
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(ctx, k.stakeKeeper)
	totalPower := k.stakeKeeper.GetLastTotalPower(ctx)
	highestConsensusRatio := float64(highestClaimPower) / float64(totalPower.Int64())
	remainingPossibleClaimPower := totalPower.Int64() - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := float64(highestPossibleClaimPower) / float64(totalPower.Int64())
	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.Text = types.SuccessStatusText
		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.consensusNeeded {
		prophecy.Status.Text = types.FailedStatusText
	}
	return prophecy
}
