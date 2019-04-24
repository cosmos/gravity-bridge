package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper  bank.Keeper
	stakeKeeper staking.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	codespace sdk.CodespaceType

	consensusNeeded float64
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(coinKeeper bank.Keeper, stakeKeeper staking.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, codespace sdk.CodespaceType, consensusNeeded float64) (Keeper, sdk.Error) {
	if consensusNeeded <= 0 || consensusNeeded > 1 {
		return Keeper{}, types.ErrMinimumConsensusNeededInvalid(codespace)
	}
	return Keeper{
		coinKeeper:      coinKeeper,
		stakeKeeper:     stakeKeeper,
		storeKey:        storeKey,
		cdc:             cdc,
		codespace:       codespace,
		consensusNeeded: consensusNeeded,
	}, nil
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
	if !store.Has([]byte(id)) {
		return types.NewEmptyProphecy(), types.ErrProphecyNotFound(k.Codespace())
	}
	bz := store.Get([]byte(id))
	var dbProphecy types.DBProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &dbProphecy)

	deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return types.NewEmptyProphecy(), types.ErrInternalDB(k.Codespace(), err)
	}
	return deSerializedProphecy, nil
}

// saveProphecy saves a prophecy with an initial claim
func (k Keeper) saveProphecy(ctx sdk.Context, prophecy types.Prophecy) sdk.Error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier(k.Codespace())
	}
	if len(prophecy.ClaimValidators) <= 0 {
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

func (k Keeper) ProcessClaim(ctx sdk.Context, id string, validator sdk.ValAddress, claim string) (types.Status, sdk.Error) {
	if claim == "" {
		return types.Status{}, types.ErrInvalidClaim(k.Codespace())
	}
	prophecy, err := k.GetProphecy(ctx, id)
	if err == nil {
		if prophecy.Status.StatusText == types.SuccessStatusText || prophecy.Status.StatusText == types.FailedStatusText {
			return types.Status{}, types.ErrProphecyFinalized(k.Codespace())
		}
		if prophecy.ValidatorClaims[validator.String()] != "" {
			return types.Status{}, types.ErrDuplicateMessage(k.Codespace())
		}
		prophecy.AddClaim(validator, claim)
	} else {
		if err.Code() != types.CodeProphecyNotFound {
			return types.Status{}, err
		}
		prophecy = types.NewProphecy(id)
		prophecy.AddClaim(validator, claim)
	}
	prophecy = k.processCompletion(ctx, prophecy)
	err = k.saveProphecy(ctx, prophecy)
	if err != nil {
		return types.Status{}, err
	}
	return prophecy.Status, nil
}

func (k Keeper) processCompletion(ctx sdk.Context, prophecy types.Prophecy) types.Prophecy {
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(ctx, k.stakeKeeper)
	totalPower := k.stakeKeeper.GetLastTotalPower(ctx)
	highestConsensusRatio := float64(highestClaimPower) / float64(totalPower.Int64())
	remainingPossibleClaimPower := totalPower.Int64() - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := float64(highestPossibleClaimPower) / float64(totalPower.Int64())
	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.StatusText = types.SuccessStatusText
		prophecy.Status.FinalClaim = highestClaim
		//TODO: trigger minting
	} else if highestPossibleConsensusRatio <= k.consensusNeeded {
		prophecy.Status.StatusText = types.FailedStatusText
	}
	return prophecy
}
