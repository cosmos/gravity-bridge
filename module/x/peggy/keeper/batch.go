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
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, voucherDenom types.VoucherDenom, maxElements int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	bridgedDenom := k.GetCounterpartDenominator(ctx, voucherDenom)
	if bridgedDenom == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "bridged denominator")
	}
	selectedTx, err := k.pickUnbatchedTX(ctx, voucherDenom, *bridgedDenom, maxElements)
	if len(selectedTx) == 0 || err != nil {
		return nil, err
	}
	totalFee := selectedTx[0].BridgeFee
	for _, tx := range selectedTx[1:] {
		totalFee = totalFee.Add(tx.BridgeFee)
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	nonce := types.NewUInt64Nonce(nextID)
	batch := types.OutgoingTxBatch{
		Nonce:              nonce,
		Elements:           selectedTx,
		CreatedAt:          ctx.BlockTime(),
		BridgedDenominator: *bridgedDenom,
		TotalFee:           totalFee,
		BatchStatus:        types.BatchStatusPending,
	}
	k.storeBatch(ctx, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, nonce.String()),
		sdk.NewAttribute(types.AttributeKeyNonce, nonce.String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return &batch, nil
}

func (k Keeper) storeBatch(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTxBatchKey(batch.Nonce), k.cdc.MustMarshalBinaryBare(batch))
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
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, nonce types.UInt64Nonce) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxBatchKey(nonce))
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)
	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch to the "available" second index. BatchStatus is set to canceled
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, nonce types.UInt64Nonce) error {
	batch := k.GetOutgoingTXBatch(ctx, nonce)
	if batch == nil {
		return types.ErrUnknown
	}
	if err := batch.Cancel(); err != nil {
		return err
	}
	for _, tx := range batch.Elements {
		k.prependToUnbatchedTXIndex(ctx, batch.BridgedDenominator.ToVoucherCoin(tx.BridgeFee.Amount), tx.ID)
	}
	k.storeBatch(ctx, *batch)
	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, nonce.String()),
		sdk.NewAttribute(types.AttributeKeyNonce, nonce.String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// SetOutgoingTXBatchConfirm stores the signature an orchestrator has submitted for an outgoing batch
func (k Keeper) SetOutgoingTXBatchConfirm(ctx sdk.Context, nonce types.UInt64Nonce, validator sdk.ValAddress, signature []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTXBatchConfirmKey(nonce, validator), signature)
}

// Iterate through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(key []byte, batch types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), batch) {
			break
		}
	}
}
