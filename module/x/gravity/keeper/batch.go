package keeper

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO: this should be a parameter
const BatchTxSize = 100

// CreateBatchTx starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit withtout creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) CreateBatchTx(ctx sdk.Context, contractAddress common.Address) (uint64, error) {
	// select transfer txs from outgoing pool sorted by fee in desc order
	// TODO: use parameter for batch size
	txs := k.pickUnbatchedTxs(ctx, contractAddress.String(), BatchTxSize)
	if len(txs) == 0 {
		// TODO: fix error
		return 0, sdkerrors.Wrapf(types.ErrEmpty, "batch tx failed for address %s", contractAddress)
	}

	timeoutHeight, err := k.GetBatchTimeoutHeight(ctx)
	if err != nil {
		return 0, err
	}

	// TODO: txID ≠ nonce
	// TODO: use hash for tx id
	txID := uint64(0)
	nonce := k.GetBatchID(ctx)

	batchTx := types.BatchTx{
		Nonce:         nonce,
		Timeout:       timeoutHeight,
		Transactions:  txs,
		TokenContract: contractAddress.String(),
		Block:         0, // TODO: add?
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
func (k Keeper) GetBatchTimeoutHeight(ctx sdk.Context) (uint64, error) {
	params := k.GetParams(ctx)
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	info, found := k.GetEthereumInfo(ctx)
	if !found {
		return 0, sdkerrors.Wrap(
			sdkerrors.ErrInvalidHeight,
			"tracked ethereum height is 0. Track an populate the heights through a deposit event",
		)
	}

	// calculate the time duration difference between the current block timestamp and the timestamp
	// when the last Ethereum block height was observed on the bridge
	timestampDiff := sdk.NewDec(int64(ctx.BlockTime().Sub(time.Unix(0, info.Timestamp))))

	newBlocks := timestampDiff.QuoInt64(int64(params.AverageBlockTime)).TruncateInt64()
	currentEthereumHeight := uint64(newBlocks) + info.Height

	// TODO: [IMPORTANT] ensure timeout is in blocks and not time/duration on the contract
	timeout := currentEthereumHeight + params.TargetBatchTimeout
	return timeout, nil
}

// OnBatchTxExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OnBatchTxExecuted(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batch, found := k.GetBatchTx(ctx, tokenContract, nonce)
	if !found {
		// TODO: fix error msg
		return sdkerrors.Wrap(types.ErrOutgoingTxNotFound, "nonce")
	}

	// cleanup outgoing TX pool
	for _, tx := range batch.Transactions {
		// TODO: get the txs ids
		k.removePoolEntry(ctx, tx.Id)
	}

	// Iterate through remaining batches
	// TODO: still not getting why we need to cancel the other batch txs
	k.IterateBatchTxs(ctx, func(batchTx types.BatchTx) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if batchTx.Nonce < batch.Nonce {
			k.CancelBatchTx(ctx, tokenContract, batchTx.Nonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatchTx(ctx, batch)
	return nil
}

func (k Keeper) GetBatchID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastOutgoingBatchID)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) setBatchID(ctx sdk.Context, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastOutgoingBatchID, sdk.Uint64ToBigEndian(txID))
}

// SetBatchTx stores a batch transaction
func (k Keeper) SetBatchTx(ctx sdk.Context, batchTx types.BatchTx) {
	store := ctx.KVStore(k.storeKey)
	// set the current block height when storing the batch
	batchTx.Block = uint64(ctx.BlockHeight())
	key := types.GetBatchTxKey(batchTx.TokenContract, batchTx.Nonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&batchTx))

	// second index to store the block height at which the batch tx was stored
	blockKey := types.GetBatchTxBlockKey(batchTx.Block)
	// TODO: only store ID?
	store.Set(blockKey, k.cdc.MustMarshalBinaryBare(&batchTx))
}

// DeleteBatchTx deletes an outgoing transaction batch
func (k Keeper) DeleteBatchTx(ctx sdk.Context, batch types.BatchTx) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBatchTxKey(batch.TokenContract, batch.Nonce))
	store.Delete(types.GetBatchTxBlockKey(batch.Block))
}

// pickUnbatchedTx find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTxs(ctx sdk.Context, contractAddress string, maxElements int) []types.TransferTx {
	var txs []types.TransferTx

	k.IterateTransferTxsByFee(ctx, contractAddress, func(txID uint64, tx types.TransferTx) bool {
		txs = append(txs, tx)
		k.removeFromUnbatchedTxIndex(ctx, txID, tx.Erc20Fee)
		return len(txs) == maxElements // stop if we've reached the limit
	})

	return txs
}

// GetBatchTx loads a batch object
func (k Keeper) GetBatchTx(ctx sdk.Context, tokenContract string, nonce uint64) (types.BatchTx, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBatchTxKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return types.BatchTx{}, false
	}

	var batchTx types.BatchTx
	k.cdc.MustUnmarshalBinaryBare(bz, &batchTx)
	return batchTx, true
}

// CancelBatchTx releases all txs in the batch and deletes the batch
func (k Keeper) CancelBatchTx(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batchTx, found := k.GetBatchTx(ctx, tokenContract, nonce)
	if !found {
		// TODO: fix error msg
		return sdkerrors.Wrap(types.ErrEmpty, "outgoing batch tx not found")
	}

	// TODO: why do we need to do this?
	for _, tx := range batchTx.Transactions {
		// store all the transactions from the batch
		k.prependToUnbatchedTxIndex(ctx, tx.Id, tokenContract, tx.Erc20Fee)
	}

	// Delete batch since it is finished
	k.DeleteBatchTx(ctx, batchTx)

	// TODO: fix events
	nonceStr := strconv.FormatUint(nonce, 64)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatchCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, nonceStr),
			sdk.NewAttribute(types.AttributeKeyNonce, nonceStr),
		),
	)
	return nil
}

// IterateBatchTxs iterates through all outgoing batches in DESC order.
// TODO: add tx id to callback
func (k Keeper) IterateBatchTxs(ctx sdk.Context, cb func(batch types.BatchTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.BatchTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}

// GetBatchTxs returns the outgoing tx batches
func (k Keeper) GetBatchTxs(ctx sdk.Context) []types.BatchTx {
	txs := []types.BatchTx{}
	k.IterateBatchTxs(ctx, func(batchTx types.BatchTx) bool {
		txs = append(txs, batchTx)
		return false
	})

	return txs
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	batches := k.GetOutgoingTxBatches(ctx)
	var lastBatch *types.OutgoingTxBatch = nil
	lastNonce := uint64(0)
	for _, batch := range batches {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
	}
	return lastBatch
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

// GetUnslashedBatches returns all the unslashed batches in state
func (k Keeper) GetUnslashedBatches(ctx sdk.Context, maxHeight uint64) []types.BatchTx {
	txs := []types.BatchTx{}
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)

	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock+1, maxHeight, func(batch types.BatchTx) bool {
		txs = append(txs, batch)
		return false
	})

	return txs
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(ctx sdk.Context, lastSlashedBatchBlock uint64, maxHeight uint64, cb func(types.BatchTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxBlockKey)
	iter := prefixStore.Iterator(sdk.Uint64ToBigEndian(lastSlashedBatchBlock), sdk.Uint64ToBigEndian(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var batch types.BatchTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}
