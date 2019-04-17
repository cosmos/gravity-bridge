package keeper

import (
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
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
		codespace:  codespace,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetProphecy gets the entire prophecy data struct for a given id
func (k Keeper) GetProphecy(ctx sdk.Context, id string) (types.BridgeProphecy, sdk.Error) {
	if id == "" {
		return types.NewEmptyBridgeProphecy(), types.ErrInvalidIdentifier(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return types.NewEmptyBridgeProphecy(), types.ErrNotFound(k.Codespace())
	}
	bz := store.Get([]byte(id))
	var prophecy types.BridgeProphecy
	k.cdc.MustUnmarshalBinaryBare(bz, &prophecy)
	return prophecy, nil
}

// CreateProphecy creates a new prophecy with an initial claim
func (k Keeper) CreateProphecy(ctx sdk.Context, prophecy types.BridgeProphecy) sdk.Error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier(k.Codespace())
	}
	if prophecy.MinimumPower < 2 {
		return types.ErrMinimumPowerTooLow(k.Codespace())
	}
	if len(prophecy.BridgeClaims) <= 0 {
		return types.ErrNoClaims(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(prophecy))
	return nil
}
