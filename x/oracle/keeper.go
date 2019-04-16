package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// GetProphecy gets the entire prophecy data struct for a given id
func (k Keeper) GetProphecy(ctx sdk.Context, id string) (BridgeProphecy, sdk.Error) {
	if id == "" {
		return NewEmptyBridgeProphecy(), ErrInvalidIdentifier(DefaultCodespace)
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return NewEmptyBridgeProphecy(), ErrNotFound(DefaultCodespace)
	}
	bz := store.Get([]byte(id))
	var prophecy BridgeProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &prophecy)
	return prophecy, nil
}

// Creates a new prophecy with an initial claim
func (k Keeper) createProphecy(ctx sdk.Context, prophecy BridgeProphecy) sdk.Error {
	if prophecy.ID == "" {
		return ErrInvalidIdentifier(DefaultCodespace)
	}
	if prophecy.MinimumClaims < 2 {
		return ErrMinimumTooLow(DefaultCodespace)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(prophecy))
	return nil
}
