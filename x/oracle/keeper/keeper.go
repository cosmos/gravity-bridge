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
	var prophecy types.Prophecy
	k.cdc.MustUnmarshalBinaryBare(bz, &prophecy)
	return prophecy, nil
}

// CreateProphecy creates a new prophecy with an initial claim
func (k Keeper) CreateProphecy(ctx sdk.Context, prophecy types.Prophecy) sdk.Error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier(k.Codespace())
	}
	if prophecy.MinimumPower < 2 {
		return types.ErrMinimumPowerTooLow(k.Codespace())
	}
	if len(prophecy.Claims) <= 0 {
		return types.ErrNoClaims(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(prophecy))
	return nil
}

func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.Claim) (types.ProgressUpdate, sdk.Error) {
	_, err := k.GetProphecy(ctx, claim.ID)
	if err == nil {
		//check if complete or not
		return types.ProgressUpdate{}, sdk.ErrInternal("Not yet implemented")
		//	//check if claim for this validator exists or not
		//	//if does
		//		//return error
		//	//else
		//		//add claim to list
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
	} else {
		if err.Code() != types.CodeProphecyNotFound {
			return types.ProgressUpdate{}, err
		}
		claims := []types.Claim{claim}
		newProphecy := types.NewProphecy(claim.ID, types.PendingStatus, k.getPowerThreshold(), claims)
		err := k.CreateProphecy(ctx, newProphecy)
		if err != nil {
			return types.ProgressUpdate{}, err
		}
		//return result
		return types.NewProgressUpdate(types.PendingStatus, nil), nil
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
