package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// GetLogicCallTx gets an  logic call
func (k Keeper) GetLogicCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) (types.LogicCallTx, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyOutgoingLogicCall)
	bz := store.Get(types.GetLogicCallTxKey(invalidationID, invalidationNonce))
	if len(bz) == 0 {
		return types.LogicCallTx{}, false
	}

	var call types.LogicCallTx
	k.cdc.MustUnmarshalBinaryBare(bz, &call)
	return call, true
}

// SetLogicCallTx sets an  logic call tx to the store
func (k Keeper) SetLogicCallTx(ctx sdk.Context, invalidationID tmbytes.HexBytes, invalidationNonce uint64, call types.LogicCallTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyOutgoingLogicCall)
	store.Set(types.GetLogicCallTxKey(invalidationID, invalidationNonce), k.cdc.MustMarshalBinaryBare(&call))
}

// DeleteLogicCallTx deletes a given logic call from the store
func (k Keeper) DeleteLogicCallTx(ctx sdk.Context, invalidationID tmbytes.HexBytes, invalidationNonce uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyOutgoingLogicCall)
	store.Delete(types.GetLogicCallTxKey(invalidationID, invalidationNonce))
}

// CancelLogicCallTx releases removes a given logic call tx from the store and
// emits events and logs
func (k Keeper) CancelLogicCallTx(ctx sdk.Context, invalidationID tmbytes.HexBytes, invalidationNonce uint64) {
	k.DeleteLogicCallTx(ctx, invalidationID, invalidationNonce)

	nonceStr := strconv.FormatUint(invalidationNonce, 64)
	k.Logger(ctx).Info(
		"logic call cancelled",
		"invalidation-id", invalidationID.String(),
		"invalidation-nonce", nonceStr,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLogicCallCanceled,
			sdk.NewAttribute(types.AttributeKeyInvalidationID, invalidationID.String()),
			sdk.NewAttribute(types.AttributeKeyInvalidationNonce, nonceStr),
		),
	)
}

// IterateLogicCallTxs iterates over outgoing logic calls
func (k Keeper) IterateLogicCallTxs(ctx sdk.Context, cb func(invalidationID tmbytes.HexBytes, invalidationNonce uint64, tx types.LogicCallTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyOutgoingLogicCall)
	iterator := prefixStore.Iterator(nil, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tx types.LogicCallTx
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &tx)

		invalidationID := tmbytes.HexBytes{}
		invalidationNonce := uint64(0)
		if cb(invalidationID, invalidationNonce, tx) {
			break //  stop iteration
		}
	}
}

// GetIdentifiedLogicCalls returns the all the outgoing logic txs with they corresponding
// store key (invalidation id and nonce).
func (k Keeper) GetIdentifiedLogicCalls(ctx sdk.Context) []types.IdentifiedLogicCall {
	txs := []types.IdentifiedLogicCall{}

	k.IterateLogicCallTxs(ctx, func(invalidationID tmbytes.HexBytes, invalidationNonce uint64, tx types.LogicCallTx) bool {
		call := types.IdentifiedLogicCall{
			InvalidationID:    invalidationID,
			InvalidationNonce: invalidationNonce,
			LogicCall:         tx,
		}
		txs = append(txs, call)
		return false
	})

	return txs
}
