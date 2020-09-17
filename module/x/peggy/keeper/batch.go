package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const OutgoingTxBatchSize = 100

func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, voucherDenom types.VoucherDenom, maxElements int) (uint64, error) {
	if maxElements == 0 {
		return 0, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	bridgedDenom, err := k.GetCounterpartDenominator(ctx, voucherDenom)
	if err != nil {
		return 0, err
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
		Elements:              selectedTx,
		CreatedAt:             ctx.BlockTime(),
		CosmosDenom:           voucherDenom,
		BridgedTokenID:        bridgedDenom.TokenID,
		BridgeContractAddress: bridgedDenom.BridgeContractAddress,
		TotalFee:              totalFee,
		BatchStatus:           types.BatchStatusPending,
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	k.storeBatch(ctx, nextID, batch)
	return nextID, nil
}

func (k Keeper) storeBatch(ctx sdk.Context, nextID uint64, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(sdk.Uint64ToBigEndian(nextID), k.cdc.MustMarshalBinaryBare(batch))
}

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
		err = k.RemoveFromUnbatchedTXIndex(ctx, tx.BridgeFee, txID)
		return err != nil || len(selectedTx) == maxElements
	})
	return selectedTx, err
}

func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, batchID uint64) (*types.OutgoingTxBatch, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sdk.Uint64ToBigEndian(batchID))
	if len(bz) == 0 {
		return nil, types.ErrUnknown
	}
	var b types.OutgoingTxBatch
	return &b, k.cdc.UnmarshalBinaryBare(bz, &b)
}

func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, batchID uint64) error {
	batch, err := k.GetOutgoingTXBatch(ctx, batchID)
	if err != nil {
		return err
	}
	if err := batch.Cancel(); err != nil {
		return err
	}
	for _, tx := range batch.Elements {
		k.prependToUnbatchedTXIndex(ctx, sdk.NewInt64Coin(string(batch.CosmosDenom), int64(tx.BridgeFee.Amount)), tx.ID)
	}
	k.storeBatch(ctx, batchID, *batch)
	return nil
}
