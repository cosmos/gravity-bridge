package oracle

import (
	"strconv"

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

// Gets the entire prophecy data struct for a given nonce
func (k Keeper) GetProphecy(ctx sdk.Context, nonce int) (BridgeProphecy, sdk.Error) {
	if nonce < 0 {
		return NewEmptyBridgeProphecy(), ErrInvalidNonce(DefaultCodespace)
	}
	nonceKey := strconv.Itoa(nonce)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(nonceKey)) {
		return NewEmptyBridgeProphecy(), ErrNotFound(DefaultCodespace)
	}
	bz := store.Get([]byte(nonceKey))
	var prophecy BridgeProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &prophecy)
	return prophecy, nil
}

// Creates a new prophecy with an initial claim
func (k Keeper) createProphecy(ctx sdk.Context, prophecy BridgeProphecy) sdk.Error {
	if prophecy.Nonce < 0 {
		return ErrInvalidNonce(DefaultCodespace)
	}
	if prophecy.MinimumClaims < 2 {
		return ErrMinimumTooLow(DefaultCodespace)
	}
	nonceKey := strconv.Itoa(prophecy.Nonce)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(nonceKey), k.cdc.MustMarshalBinaryBare(prophecy))
	return nil
}
