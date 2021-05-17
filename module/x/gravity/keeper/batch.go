package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

const BatchTxSize = 100

// BuildBatchTx starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit withtout creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildBatchTx(
	ctx sdk.Context,
	contractAddress common.Address,
	maxElements int) (*types.BatchTx, error) {

	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}

	if lastBatch := k.GetLastOutgoingBatchByTokenType(ctx, contractAddress); lastBatch != nil {
		currentFees := k.GetBatchFeesByTokenType(ctx, contractAddress)
		if currentFees == nil {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "error getting fees from tx pool")
		}

		lastFees := lastBatch.GetFees()
		if lastFees.GT(currentFees.Amount) {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "new batch would not be more profitable")
		}
	}

	selectedTx, err := k.PickUnbatchedTX(ctx, contractAddress, maxElements)
	if err != nil {
		return nil, err
	}
	if len(selectedTx) == 0 {
		return nil, nil
	}

	nextNonce := k.IncrementLastOutgoingBatchNonce(ctx)
	batch := &types.BatchTx{
		Nonce:         nextNonce,
		Timeout:       k.getBatchTimeoutHeight(ctx),
		Transactions:  selectedTx,
		TokenContract: contractAddress.Hex(),
	}
	k.SetOutgoingTx(ctx, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nextNonce)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nextNonce)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return batch, nil
}

// This gets the batch timeout height in Ethereum blocks.
func (k Keeper) getBatchTimeoutHeight(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	currentCosmosHeight := ctx.BlockHeight()
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values are zero because
	// no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedEthereumBlockHeight(ctx)
	if heights.CosmosHeight == 0 || heights.EthereumHeight == 0 {
		return 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentCosmosHeight) - heights.CosmosHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projectedCurrentEthereumHeight := (projectedMillis / params.AverageEthereumBlockTime) + heights.EthereumHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	blocksToAdd := params.TargetBatchTimeout / params.AverageEthereumBlockTime
	return projectedCurrentEthereumHeight + blocksToAdd
}

// BatchTxExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It deletes all the transactions in the batch, then cancels all earlier batches
func (k Keeper) BatchTxExecuted(ctx sdk.Context, tokenContract common.Address, nonce uint64) error {
	otx := k.GetOutgoingTx(ctx, types.MakeBatchTxKey(tokenContract, nonce))
	if otx == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	batchTx, _ := otx.(*types.BatchTx)

	// cleanup outgoing TX pool
	for _, tx := range batchTx.Transactions {
		k.DeletePoolEntry(ctx, tx.Id)
	}
	var err error
	// Iterate through remaining batches
	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		iterBatch, _ := otx.(*types.BatchTx)
		// todo: reformat the store to additionally prefix by token type,
		// and also sort by nonce, to save on needless iteration
		if (iterBatch.Nonce < batchTx.Nonce) && (batchTx.TokenContract == tokenContract.Hex()) {
			err = k.CancelBatchTx(ctx, tokenContract, iterBatch.Nonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteOutgoingTx(ctx, batchTx.GetStoreIndex())

	return err
}

// PickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) PickUnbatchedTX(
	ctx sdk.Context,
	contractAddress common.Address,
	maxElements int) ([]*types.SendToEthereum, error) {
	var selectedTx []*types.SendToEthereum
	var err error
	k.IterateOutgoingPoolByFee(ctx, contractAddress, func(txID uint64, tx *types.SendToEthereum) bool {
		if tx != nil && tx.Erc20Fee.GravityCoin().IsZero() {
			selectedTx = append(selectedTx, tx)
			err = k.removeFromUnbatchedTXIndex(ctx, tx.Erc20Fee, txID)
			return err != nil || len(selectedTx) == maxElements
		}

		return true
	})
	return selectedTx, err
}

// CancelBatchTx releases all TX in the batch and deletes the batch
func (k Keeper) CancelBatchTx(ctx sdk.Context, tokenContract common.Address, nonce uint64) error {
	otx := k.GetOutgoingTx(ctx, types.MakeBatchTxKey(tokenContract, nonce))
	if otx == nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "no batch tx found for token: %s nonce: %d", tokenContract, nonce)
	}
	batch, _ := otx.(*types.BatchTx)
	for _, tx := range batch.Transactions {
		k.AppendToUnbatchedTXIndex(ctx, tx.Erc20Fee, tx.Id)
	}

	// Delete batch since it is finished
	k.DeleteOutgoingTx(ctx, batch.GetStoreIndex())

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nonce)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nonce)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// IterateBatchTxs iterates through all outgoing batches in DESC order.
func (k Keeper) IterateBatchTxs(ctx sdk.Context, cb func(key []byte, batch *types.BatchTx) bool) {
	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)
		return cb(key, btx)
	})
}

// GetBatchTxes returns the outgoing tx batches
func (k Keeper) GetBatchTxes(ctx sdk.Context) (out []*types.BatchTx) {
	k.IterateBatchTxs(ctx, func(_ []byte, batch *types.BatchTx) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token common.Address) *types.BatchTx {
	batches := k.GetBatchTxes(ctx)
	var lastBatch *types.BatchTx = nil
	lastNonce := uint64(0)
	for _, batch := range batches {
		if common.HexToAddress(batch.TokenContract) == token && batch.Nonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.Nonce
		}
	}
	return lastBatch
}

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	ctx.KVStore(k.storeKey).Set([]byte{types.LastSlashedBatchBlockKey}, types.UInt64Bytes(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastSlashedBatchBlockKey}); bz == nil {
		return 0
	} else {
		return types.UInt64FromBytes(bz)
	}
}

// GetUnSlashedBatches returns all the unslashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (out []*types.BatchTx) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)
	k.IterateBatchBySlashedBatchBlock(ctx,
		lastSlashedBatchBlock,
		maxHeight,
		func(_ []byte, batch *types.BatchTx) bool {
			if batch.Height > lastSlashedBatchBlock {
				out = append(out, batch)
			}
			return false
		})
	return
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(
	ctx sdk.Context,
	lastSlashedBatchBlock uint64,
	maxHeight uint64,
	cb func([]byte, *types.BatchTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.BatchTxBlockKey})
	iter := prefixStore.Iterator(types.UInt64Bytes(lastSlashedBatchBlock), types.UInt64Bytes(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var Batch types.BatchTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &Batch)
		// cb returns true to stop early
		if cb(iter.Key(), &Batch) {
			break
		}
	}
}

func (k Keeper) IncrementLastOutgoingBatchNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LastOutgoingBatchNonceKey})
	var id uint64 = 0
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	newId := id + 1
	bz = sdk.Uint64ToBigEndian(newId)
	store.Set([]byte{types.LastOutgoingBatchNonceKey}, bz)
	return newId
}
