package keeper

import (
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
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
	batch := types.OutgoingTxBatch{
		Elements:                    selectedTx,
		CreatedAt:                   ctx.BlockTime(),
		CosmosDenom:                 voucherDenom,
		BridgedTokenSymbol:          bridgedDenom.Symbol,
		BridgedTokenContractAddress: bridgedDenom.TokenContractAddress,
		TotalFee:                    totalFee,
		BatchStatus:                 types.BatchStatusPending,
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	k.storeBatch(ctx, nextID, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, types.BridgeContractAddress.String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, types.BridgeContractChainID),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, strconv.Itoa(int(nextID))),
		sdk.NewAttribute(types.AttributeKeyNonce, types.NonceFromUint64(nextID).String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nextID, nil
}

func (k Keeper) storeBatch(ctx sdk.Context, nextID uint64, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(sdk.Uint64ToBigEndian(nextID), k.cdc.MustMarshalBinaryBare(batch))
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
			Amount:      types.AsTransferCoin(bridgedDenom, tx.Amount),
			BridgeFee:   types.AsTransferCoin(bridgedDenom, tx.BridgeFee),
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
	bz := store.Get(sdk.Uint64ToBigEndian(batchID))
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
		k.prependToUnbatchedTXIndex(ctx, sdk.NewInt64Coin(string(batch.CosmosDenom), int64(tx.BridgeFee.Amount)), tx.ID)
	}
	k.storeBatch(ctx, batchID, *batch)
	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, types.BridgeContractAddress.String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, types.BridgeContractChainID),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, strconv.Itoa(int(batchID))),
		sdk.NewAttribute(types.AttributeKeyNonce, types.NonceFromUint64(batchID).String()),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}
