package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// PushToOutgoingPool add cross transaction to pool
// - checks a counterpart denomintor exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) PushToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	// Ensure that the coin is a peggy voucher
	if _, err := types.ValidatePeggyCoin(amount); err != nil {
		return 0, fmt.Errorf("amount not a peggy voucher coin: %s", err)
	}

	// TODO fee shoule be uiris
	totalInVouchers := sdk.NewCoins(amount).Add(fee)

	// send coins to module in prep for burn
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, totalInVouchers); err != nil {
		return 0, err
	}

	// burn vouchers to send them back to ETH
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, totalInVouchers); err != nil {
		panic(err)
	}

	// get next tx id from keeper
	nextID := k.autoIncrementID(ctx, types.KeyLastTXPoolID)

	// construct outgoing tx
	outgoing := &types.OutgoingTx{
		Sender:    sender.String(),
		DestAddr:  counterpartReceiver,
		Amount:    amount,
		BridgeFee: fee,
	}

	// set the outgoing tx in the pool index
	if err := k.setPoolEntry(ctx, nextID, outgoing); err != nil {
		return 0, err
	}

	// add a second index with the fee
	k.appendUnbatchedTxByFee(ctx, fee, nextID)

	// todo: add second index for sender so that we can easily query: give pending Tx by sender
	// todo: what about a second index for receiver?

	poolEvent := sdk.NewEvent(
		types.EventTypeBridgeWithdrawalReceived,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingTXID, strconv.Itoa(int(nextID))),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nextID)),
	)
	ctx.EventManager().EmitEvent(poolEvent)

	return nextID, nil
}

// BuildTxBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildTxBatch(ctx sdk.Context, maxElements int) error {
	if maxElements == 0 {
		return sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	// TODO: figure out how to check for know or unknown denoms? this might not matter anymore
	selectedTx, err := k.pickUnbatchedTx(ctx, maxElements)
	if len(selectedTx) == 0 || err != nil {
		return err
	}

	if err := k.decrPoolUnbatchedTxCnt(ctx, len(selectedTx)); err != nil {
		return sdkerrors.Wrap(types.ErrInternal, err.Error())
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)

	batch := &types.OutgoingTxBatch{
		BatchNonce:   nextID,
		Transactions: selectedTx,
		// Valset:        k.GetCurrentValset(ctx),
		//TokenContract: contractAddress,
	}
	k.storeBatch(ctx, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nextID)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nextID)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// IteratePoolTxByFee itetates over the outgoing pool which is sorted by fee
func (k Keeper) IteratePoolTxByFee(ctx sdk.Context, cb func(uint64, *types.OutgoingTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.ReverseIterator(prefixRange([]byte{}))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ids types.IDSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)
		// cb returns true to stop early
		for _, id := range ids.Ids {
			tx, err := k.getPoolEntry(ctx, id)
			if err != nil {
				return
			}
			if cb(id, tx) {
				return
			}
		}
	}
	return
}

// appendToUnbatchedTXIndex add at the end when tx with same fee exists
func (k Keeper) appendUnbatchedTxByFee(ctx sdk.Context, fee sdk.Coin, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	if store.Has(idxKey) {
		bz := store.Get(idxKey)
		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	}
	idSet.Ids = append(idSet.Ids, txID)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
}

func (k Keeper) GetUnbatchedTx(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UnbatchedTxCountKey)
	var total uint64 = 0
	if bz != nil {
		total = binary.BigEndian.Uint64(bz)
	}
	return total
}

func (k Keeper) incrPoolUnbatchedTxCnt(ctx sdk.Context) {
	var total = k.GetUnbatchedTx(ctx)
	bz := sdk.Uint64ToBigEndian(total + 1)

	store := ctx.KVStore(k.storeKey)
	store.Set(types.UnbatchedTxCountKey, bz)
}

func (k Keeper) decrPoolUnbatchedTxCnt(ctx sdk.Context, delta int) error {
	var total = k.GetUnbatchedTx(ctx)
	if total-uint64(delta) < 0 {
		return fmt.Errorf("unbatched tx count %d, but decrement %d", total, delta)
	}

	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(total - uint64(delta))
	store.Set(types.UnbatchedTxCountKey, bz)
	return nil
}
