package keeper

import (
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const OutgoingTxBatchSize = 100

// BuildOutgoingTXBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, voucherDenom types.VoucherDenom, maxElements int) (uint64, error) {
	if maxElements == 0 {
		return 0, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	bridgedDenom := k.GetCounterpartDenominator(ctx, voucherDenom)
	if bridgedDenom == nil {
		return 0, sdkerrors.Wrap(types.ErrUnknown, "bridged denominator")
	}
	selectedTx, err := k.pickUnbatchedTX(ctx, voucherDenom, *bridgedDenom, maxElements)
	if len(selectedTx) == 0 || err != nil {
		return 0, err
	}
	totalFee := selectedTx[0].BridgeFee
	for _, tx := range selectedTx[1:] {
		totalFee = totalFee.Add(tx.BridgeFee)
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := types.OutgoingTxBatch{
		Nonce:              types.NonceFromUint64(nextID),
		Elements:           selectedTx,
		CreatedAt:          ctx.BlockTime(),
		BridgedDenominator: *bridgedDenom,
		TotalFee:           totalFee,
		BatchStatus:        types.BatchStatusPending,
	}
	k.storeBatch(ctx, nextID, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, strconv.Itoa(int(nextID))),
		sdk.NewAttribute(types.AttributeKeyNonce, types.NonceFromUint64(nextID).String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nextID, nil
}

func (k Keeper) storeBatch(ctx sdk.Context, batchID uint64, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTxBatchKey(batchID), k.cdc.MustMarshalBinaryBare(batch))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTX(ctx sdk.Context, denom types.VoucherDenom, bridgedDenom types.BridgedDenominator, maxElements int) ([]types.OutgoingTransferTx, error) {
	var selectedTx []types.OutgoingTransferTx
	var err error
	k.IterateOutgoingPoolByFee(ctx, denom, func(txID uint64, tx types.OutgoingTx) bool {
		txOut := types.OutgoingTransferTx{
			ID:          txID,
			Sender:      tx.Sender,
			DestAddress: tx.DestAddress,
			Amount:      bridgedDenom.ToERC20Token(tx.Amount),
			BridgeFee:   bridgedDenom.ToERC20Token(tx.BridgeFee),
		}
		selectedTx = append(selectedTx, txOut)
		err = k.removeFromUnbatchedTXIndex(ctx, tx.BridgeFee, txID)
		return err != nil || len(selectedTx) == maxElements
	})
	return selectedTx, err
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, batchID uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxBatchKey(batchID))
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)
	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch to the "available" second index. BatchStatus is set to canceled
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, batchID uint64) error {
	batch := k.GetOutgoingTXBatch(ctx, batchID)
	if batch == nil {
		return types.ErrUnknown
	}
	if err := batch.Cancel(); err != nil {
		return err
	}
	for _, tx := range batch.Elements {
		k.prependToUnbatchedTXIndex(ctx, batch.BridgedDenominator.ToVoucherCoin(tx.BridgeFee.Amount), tx.ID)
	}
	k.storeBatch(ctx, batchID, *batch)
	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, strconv.Itoa(int(batchID))),
		sdk.NewAttribute(types.AttributeKeyNonce, types.NonceFromUint64(batchID).String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// SetOutgoingTXBatchConfirm stores the signature an orchestrator has submitted for an outgoing batch
func (k Keeper) SetOutgoingTXBatchConfirm(ctx sdk.Context, batchID uint64, validator sdk.ValAddress, signature []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTXBatchConfirmKey(batchID, validator), signature)
}

// HasOutgoingTXBatchConfirm returns true when a signature was persisted for the given batch and validator address
func (k Keeper) HasOutgoingTXBatchConfirm(ctx sdk.Context, batchID uint64, validatorAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOutgoingTXBatchConfirmKey(batchID, validatorAddr))
}

// Iterate through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(batchID uint64, batch types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(types.DecodeUin64(iter.Key()), batch) {
			break
		}
	}
}
