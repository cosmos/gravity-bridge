package keeper

import (
	"crypto/sha256"
	"strconv"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// CreateBatchTx starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit withtout creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) CreateBatchTx(ctx sdk.Context, contractAddress common.Address) (tmbytes.HexBytes, error) {
	// select transfer txs from outgoing transfer tx pool sorted by fee in desc order
	// TODO: mark the transfer txs as batched
	txs := k.SelectTransferTxs(ctx, contractAddress)
	if len(txs) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrEmpty, "batch tx failed for address %s", contractAddress)
	}

	timeoutHeight, err := k.GetBatchTimeoutHeight(ctx)
	if err != nil {
		return nil, err
	}

	// get the latest batch nonce for replay protection
	nonce := k.GetLastBatchNonce(ctx)
	nonce++

	batchTx := types.BatchTx{
		Nonce:         nonce,
		Timeout:       timeoutHeight,
		Transactions:  txs,
		TokenContract: contractAddress.String(),
		Block:         uint64(ctx.BlockHeight()),
	}

	// set batch tx by transaction ID and block height
	txID := k.SetBatchTx(ctx, batchTx)
	k.SetBatchTxByBlock(ctx, batchTx.Block, txID)

	// set the updated batch nonce for replay protection
	k.setLastBatchNonce(ctx, nonce)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatch,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTxID, txID.String()),
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
	timestampDiff := sdk.NewDec(int64(ctx.BlockTime().Sub(info.Timestamp)))

	newBlocks := timestampDiff.QuoInt64(int64(params.AverageBlockTime)).TruncateInt64()
	currentEthereumHeight := uint64(newBlocks) + info.Height

	// TODO: [IMPORTANT] ensure timeout is in blocks and not time/duration on the contract
	timeout := currentEthereumHeight + params.TargetBatchTimeout
	return timeout, nil
}

// OnBatchTxExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OnBatchTxExecuted(ctx sdk.Context, tokenContract common.Address, batchTxID tmbytes.HexBytes) error {
	batchTx, found := k.GetBatchTx(ctx, tokenContract, batchTxID)
	if !found {
		// TODO: fix error msg
		return sdkerrors.Wrapf(types.ErrOutgoingTxNotFound, "contract %s and tx id %s")
	}

	// cleanup outgoing transfer tx pool
	for _, transferTxID := range batchTx.Transactions {
		k.DeleteTransferTx(ctx, transferTxID)
	}

	// Iterate through remaining batches with the same token contract
	// TODO: still not getting why we need to cancel the other batch txs?
	k.IterateBatchTxsByToken(ctx, tokenContract, func(txID tmbytes.HexBytes, batch types.BatchTx) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if batch.Nonce < batchTx.Nonce {
			k.CancelBatchTx(ctx, tokenContract, txID, batchTx)
		}

		return false
	})

	k.DeleteBatchTx(ctx, tokenContract, batchTxID, batchTx.Block)
	return nil
}

func (k Keeper) GetLastBatchNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastBatchTxNonce)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) setLastBatchNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastBatchTxNonce, sdk.Uint64ToBigEndian(nonce))
}

// SelectTransferTxs selects the transfer transactions with the highest fee for a given contract
// until the batch is full or all the txs have been selected
func (k Keeper) SelectTransferTxs(ctx sdk.Context, tokenContract common.Address) []tmbytes.HexBytes {
	var txs []tmbytes.HexBytes

	params := k.GetParams(ctx)

	k.IterateTransferPoolByFee(ctx, tokenContract, func(_ uint64, txIDs []tmbytes.HexBytes) bool {
		size := len(txs) + len(txIDs)
		switch {
		case size < int(params.BatchSize):
			txs = append(txs, txIDs...)
			// TODO: delete transactions from second index
			return false
		case size == int(params.BatchSize):
			txs = append(txs, txIDs...)
			// TODO: delete transactions from second index
			return true // stop if we've reached the limit
		default: // case size > int(params.BatchSize):
			diff := size - int(params.BatchSize)
			txs = append(txs, txIDs[:diff-1]...)
			// TODO: set the diff to store
			return true // stop if we've reached the limit
		}
	})

	return txs
}

// GetBatchTx loads a batch object
func (k Keeper) GetBatchTx(ctx sdk.Context, tokenContract common.Address, txID tmbytes.HexBytes) (types.BatchTx, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxKey)
	key := types.GetBatchTxKey(tokenContract.String(), txID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return types.BatchTx{}, false
	}

	var batchTx types.BatchTx
	k.cdc.MustUnmarshalBinaryBare(bz, &batchTx)
	return batchTx, true
}

// SetBatchTx stores a batch transaction
func (k Keeper) SetBatchTx(ctx sdk.Context, batchTx types.BatchTx) tmbytes.HexBytes {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxKey)
	bz := k.cdc.MustMarshalBinaryBare(&batchTx)

	hash := sha256.Sum256(bz)
	txID := tmbytes.HexBytes(hash[:])

	key := types.GetBatchTxKey(batchTx.TokenContract, txID)
	store.Set(key, bz)
	return txID
}

// DeleteBatchTx deletes an outgoing transaction batch
func (k Keeper) DeleteBatchTx(ctx sdk.Context, tokenContract common.Address, txID tmbytes.HexBytes, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.BatchTxKey, types.GetBatchTxKey(tokenContract.String(), txID)...))
	store.Delete(append(types.BatchTxBlockKey, types.GetBatchTxBlockKey(blockHeight)...))
}

func (k Keeper) SetBatchTxByBlock(ctx sdk.Context, blockHeight uint64, txID tmbytes.HexBytes) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxBlockKey)
	// second index to store the block height at which the batch tx was stored
	blockKey := types.GetBatchTxBlockKey(blockHeight)

	store.Set(blockKey, txID.Bytes())
}

// CancelBatchTx releases all txs in the batch and deletes the batch
func (k Keeper) CancelBatchTx(ctx sdk.Context, tokenContract common.Address, txID tmbytes.HexBytes, batchTx types.BatchTx) error {
	// since the transaction is now canceled, we need to readd each of the
	// transactions to the second index by fee
	for _, transferID := range batchTx.Transactions {
		_, found := k.GetTransferTx(ctx, transferID)
		if !found {
			// TODO: panic? this shouldn't happen regardless because we only delete transfer txs from
			// the store once the OnBatchTxExecuted callback is executed
			sdkerrors.Wrapf(
				types.ErrOutgoingTxNotFound, // TODO: fix
				"transfer tx not found for contract %s and tx id %s", tokenContract, transferID,
			)
		}

		// TODO: only set once!!
		// k.prependToUnbatchedTxIndex(ctx, txID, tokenContract, tx.Erc20Fee)
	}

	// Delete batch since it is finished
	k.DeleteBatchTx(ctx, tokenContract, txID, batchTx.Block)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatchCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTxID, txID.String()),
			sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(batchTx.Nonce, 64)),
		),
	)
	return nil
}

// IterateBatchTxs iterates through all outgoing batches in DESC order.
func (k Keeper) IterateBatchTxs(ctx sdk.Context, cb func(tokenContract common.Address, txID tmbytes.HexBytes, batch types.BatchTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchTxKey)
	// FIXME: use iterator of fix order iteration
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.BatchTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)

		// FIXME: split key
		tokenContract := common.Address{}
		txID := tmbytes.HexBytes{}
		// cb returns true to stop early
		if cb(tokenContract, txID, batch) {
			break
		}
	}
}

// IterateBatchTxsByToken iterates over all the outgoing batches with a given token address
func (k Keeper) IterateBatchTxsByToken(ctx sdk.Context, tokenContract common.Address, cb func(txID tmbytes.HexBytes, batch types.BatchTx) bool) {
	key := append(types.BatchTxKey, []byte(tokenContract.String())...)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), key)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.BatchTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)

		txID := tmbytes.HexBytes{}
		// cb returns true to stop early
		if cb(txID, batch) {
			break
		}
	}
}

// GetBatchTxs returns the outgoing tx batches
func (k Keeper) GetBatchTxs(ctx sdk.Context) []types.BatchTx {
	txs := []types.BatchTx{}
	k.IterateBatchTxs(ctx, func(tokenContract common.Address, txID tmbytes.HexBytes, batchTx types.BatchTx) bool {
		txs = append(txs, batchTx)
		return false
	})

	return txs
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token common.Address) types.BatchTx {
	lastNonce := uint64(0)
	lastBatch := types.BatchTx{}

	k.IterateBatchTxsByToken(ctx, token, func(txID tmbytes.HexBytes, batchTx types.BatchTx) bool {
		if batchTx.Nonce > lastNonce {
			lastBatch = batchTx
			lastNonce = batchTx.Nonce
		}
		return false
	})

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
