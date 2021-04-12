package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// GetLogicCallTx gets an  logic call
func (k Keeper) GetLogicCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) (types.LogicCallTx, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLogicCallTxKey(invalidationID, invalidationNonce))
	if len(bz) == 0 {
		return types.LogicCallTx{}, false
	}

	var call types.LogicCallTx
	k.cdc.MustUnmarshalBinaryBare(bz, &call)
	return call, true
}

// SetLogicCallTx sets an  logic call tx to the store
func (k Keeper) SetLogicCallTx(ctx sdk.Context, call types.LogicCallTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLogicCallTxKey(call.InvalidationId, call.InvalidationNonce), k.cdc.MustMarshalBinaryBare(&call))
}

// DeleteLogicCallTx deletes a given logic call from the store
func (k Keeper) DeleteLogicCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetLogicCallTxKey(invalidationID, invalidationNonce))
}

// CancelLogicCallTx releases removes a given logic call tx from the store and
// emits events and logs
func (k Keeper) CancelLogicCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) {
	k.DeleteLogicCallTx(ctx, invalidationID, invalidationNonce)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingLogicCallCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyInvalidationID, fmt.Sprint(invalidationID)),
			sdk.NewAttribute(types.AttributeKeyInvalidationNonce, fmt.Sprint(invalidationNonce)),
		),
	)

	k.Logger(ctx).Debug("logic call tx cancelled")
}

// IterateLogicCallTxs iterates over outgoing logic calls
func (k Keeper) IterateLogicCallTxs(ctx sdk.Context, cb func(types.LogicCallTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyOutgoingLogicCall)
	iterator := prefixStore.Iterator(nil, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tx types.LogicCallTx
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &tx)

		if cb(tx) {
			break //  stop iteration
		}
	}
}

// GetOutgoingLogicCalls returns the all the outgoing logic txs
func (k Keeper) GetOutgoingLogicCalls(ctx sdk.Context) []types.LogicCallTx {
	txs := []types.LogicCallTx{}
	k.IterateLogicCallTxs(ctx, func(tx types.LogicCallTx) bool {
		txs = append(txs, tx)
		return false
	})

	return txs
}
