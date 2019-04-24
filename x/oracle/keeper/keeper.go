package keeper

import (
	"math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	codespace sdk.CodespaceType

	consensusNeeded float64
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, codespace sdk.CodespaceType, consensusNeeded float64) (Keeper, sdk.Error) {
	if consensusNeeded <= 0 || consensusNeeded > 1 {
		return Keeper{}, types.ErrMinimumConsensusNeededInvalid(codespace)
	}
	return Keeper{
		coinKeeper:      coinKeeper,
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

func (k Keeper) ProcessClaim(ctx sdk.Context, id string, validator sdk.AccAddress, claim string) (types.Status, sdk.Error) {
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
		//	//check if claimthreshold is passed
		//	//if does
		//		//check enough claims match and are valid
		//		//update prophecy to be successful
		//		//trigger minting
		//		//save finalized prophecy to db
		//		//return
		//	//if doesnt
		//		//check if failure threshold passed
		//		//if does, fail and update and return
		//		//else, save updated prophecy to db and return
		return types.Status{}, sdk.ErrInternal("Not yet implemented")
	} else {
		if err.Code() != types.CodeProphecyNotFound {
			return types.Status{}, err
		}
		newProphecy := types.NewProphecy(id)
		newProphecy.AddClaim(validator, claim)
		err := k.saveProphecy(ctx, newProphecy)
		if err != nil {
			return types.Status{}, err
		}
		//return result
		return prophecy.Status, nil
	}
}

func (k Keeper) getPowerThreshold() int {
	minimumPower := float64(getTotalPower()) * k.consensusNeeded
	return int(math.Ceil(minimumPower))

}

func getTotalPower() int {
	//TODO: Get from Tendermint/last block/staking module?
	return 10
}
