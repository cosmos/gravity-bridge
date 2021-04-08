package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO: this should be a parameter
const OutgoingTxBatchSize = 100

// CreateBatchTx starts the following process chain:
// - find bridged denominator for given voucher type
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) CreateBatchTx(ctx sdk.Context, contractAddress string) (uint64, error) {
	// select transfer txs from outgoing pool sorted by fee in desc order
	// TODO: use parameter for batch
	txs, err := k.pickUnbatchedTx(ctx, contractAddress, OutgoingTxBatchSize)
	if len(txs) == 0 || err != nil {
		k.Logger(ctx).Debug("batch tx failed: outgoing tx pool is empty", "address", contractAddress)
		return 0, err
	}

	timeoutHeight, err := k.GetBatchTimeoutHeight(ctx)
	if err != nil {
		return 0, err
	}

	// TODO: txID ≠ nonce
	// TODO: use hash for tx id
	txID := uint64(0)
	nonce := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)

	batchTx := types.BatchTx{
		Nonce:         nonce,
		Timeout:       timeoutHeight,
		Transactions:  txs,
		TokenContract: contractAddress,
	}

	// TODO: pass tx id as key
	k.SetBatchTx(ctx, batchTx)
	// TODO: set nonce

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatch,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, strconv.FormatUint(txID, 64)),
			sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(nonce, 64)),
		),
	)
	return txID, nil
}

// GetBatchTimeoutHeight returns the timeout block height on Ethereum based on the current bridge parameters.
//
func (k Keeper) GetBatchTimeoutHeight(ctx sdk.Context) (uint64, error) {
	params := k.GetParams(ctx)
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	ethereumInfo := k.GetLastObservedEthereumBlockHeight(ctx)
	if ethereumInfo.EthereumBlockHeight == 0 {
		// TODO: check error msg
		return 0, sdkerrors.Wrap(
			sdkerrors.ErrInvalidHeight,
			"tracked ethereum height is 0. Track an populate the heights through a deposit event",
		)
	}

	// calculate the time duration difference between the current block timestamp and the timestamp
	// when the last Ethereum block height was observed on the bridge
	timestampDiff := sdk.NewDec(int64(ctx.BlockTime().Sub(ethereumInfo.Timestamp)))

	newBlocks := timestampDiff.QuoInt64(int64(params.AverageBlockTime)).TruncateInt64()
	currentEthereumHeight := newBlocks + ethereumInfo.EthereumBlockHeight

	// TODO: ensure timeout is in blocks
	timeout := currentEthereumHeight + params.TargetBatchTimeout
	return timeout, nil
}

// OnBatchTxExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OnBatchTxExecuted(ctx sdk.Context, tokenContract string, nonce uint64) error {
	b := k.GetOutgoingTxBatch(ctx, tokenContract, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}

	// cleanup outgoing TX pool
	for _, tx := range b.Transactions {
		k.removePoolEntry(ctx, tx.Id)
	}

	// Iterate through remaining batches
	k.IterateOutgoingTxBatches(ctx, func(key []byte, iter_batch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if iter_batch.BatchNonce < b.BatchNonce {
			k.CancelOutgoingTxBatch(ctx, tokenContract, iter_batch.BatchNonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatchTx(ctx, *b)
	return nil
}

// StoreBatch stores a transaction batch
func (k Keeper) SetBatchTx(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	// set the current block height when storing the batch
	batch.Block = uint64(ctx.BlockHeight())
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)
	store.Set(blockKey, k.cdc.MustMarshalBinaryBare(batch))
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatchTx(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTx(ctx sdk.Context, contractAddress string, maxElements int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateOutgoingPoolByFee(ctx, contractAddress, func(txID uint64, tx *types.OutgoingTransferTx) bool {
		if tx != nil && tx.Erc20Fee != nil {
			selectedTx = append(selectedTx, tx)
			err = k.removeFromUnbatchedTxIndex(ctx, *tx.Erc20Fee, txID)
			return err != nil || len(selectedTx) == maxElements
		}
		// we found a nil, exit
		return true
	})
	return selectedTx, err
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTxBatch(ctx sdk.Context, tokenContract string, nonce uint64) (types.OutgoingTxBatch, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return types.OutgoingTxBatch{}, false
	}

	var batchTx types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &batchTx)

	// TODO: why are we populating these fields???
	// for _, tx := range batchTx.Transactions {
	// 	tx.Erc20Token.Contract = tokenContract
	// 	tx.Erc20Fee.Contract = tokenContract
	// }
	return batchTx, true
}

// CancelOutgoingTxBatch releases all txs in the batch and deletes the batch
func (k Keeper) CancelOutgoingTxBatch(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batchTx, found := k.GetOutgoingTxBatch(ctx, tokenContract, nonce)
	if !found {
		// TODO: fix error msg
		return sdkerrors.Wrap(types.ErrEmpty, "outgoing batch tx not found")
	}

	for _, tx := range batchTx.Transactions {
		tx.Erc20Fee.Contract = tokenContract // ?? why ?
		k.prependToUnbatchedTxIndex(ctx, tokenContract, *tx.Erc20Fee, tx.Id)
	}

	// Delete batch since it is finished
	k.DeleteBatchTx(ctx, batchTx)

	// TODO: fix events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatchCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nonce)),
			sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(nonce, 64)),
		),
	)
	return nil
}

// IterateOutgoingTXBatches iterates through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(key []byte, batch *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), &batch) {
			break
		}
	}
}

// GetOutgoingTxBatches returns the outgoing tx batches
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedBatchBlock)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// GetUnSlashedBatches returns all the unslashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (out []*types.OutgoingTxBatch) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)
	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock, maxHeight, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		if batch.Block > lastSlashedBatchBlock {
			out = append(out, batch)
		}
		return false
	})
	return
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(ctx sdk.Context, lastSlashedBatchBlock uint64, maxHeight uint64, cb func([]byte, *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxBatchBlockKey)
	iter := prefixStore.Iterator(sdk.Uint64ToBigEndian(lastSlashedBatchBlock), sdk.Uint64ToBigEndian(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var Batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &Batch)
		// cb returns true to stop early
		if cb(iter.Key(), &Batch) {
			break
		}
	}
}
