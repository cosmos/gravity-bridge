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

// AddToOutgoingPool
// - checks a counterpart denomintor exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	// Ensure that the coin is a peggy voucher
	if _, err := types.ValidatePeggyCoin(totalAmount); err != nil {
		return 0, fmt.Errorf("amount not a peggy voucher coin: %s", err)
	}

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
	k.appendToUnbatchedTXIndex(ctx, fee, nextID)

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

// appendToUnbatchedTXIndex add at the end when tx with same fee exists
func (k Keeper) appendToUnbatchedTXIndex(ctx sdk.Context, fee sdk.Coin, txID uint64) {
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

// appendToUnbatchedTXIndex add at the top when tx with same fee exists
func (k Keeper) prependToUnbatchedTXIndex(ctx sdk.Context, fee sdk.Coin, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	if store.Has(idxKey) {
		bz := store.Get(idxKey)
		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	}
	idSet.Ids = append([]uint64{txID}, idSet.Ids...)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
}

// removeFromUnbatchedTXIndex removes the tx from the index and makes it implicit no available anymore
func (k Keeper) removeFromUnbatchedTXIndex(ctx sdk.Context, fee sdk.Coin, txID uint64) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	bz := store.Get(idxKey)
	if bz == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "fee")
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	for i := range idSet.Ids {
		if idSet.Ids[i] == txID {
			idSet.Ids = append(idSet.Ids[0:i], idSet.Ids[i+1:]...)
			if len(idSet.Ids) != 0 {
				store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
			} else {
				store.Delete(idxKey)
			}
			return nil
		}
	}
	return sdkerrors.Wrap(types.ErrUnknown, "tx id")
}

func (k Keeper) setPoolEntry(ctx sdk.Context, id uint64, val *types.OutgoingTx) error {
	bz, err := k.cdc.MarshalBinaryBare(val)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTxPoolKey(id), bz)
	return nil
}

func (k Keeper) getPoolEntry(ctx sdk.Context, id uint64) (*types.OutgoingTx, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxPoolKey(id))
	if bz == nil {
		return nil, types.ErrUnknown
	}
	var r types.OutgoingTx
	k.cdc.UnmarshalBinaryBare(bz, &r)
	return &r, nil
}

func (k Keeper) removePoolEntry(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxPoolKey(id))
}

// IterateOutgoingPoolByFee itetates over the outgoing pool which is sorted by fee
func (k Keeper) IterateOutgoingPoolByFee(ctx sdk.Context, contract string, cb func(uint64, *types.OutgoingTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.ReverseIterator(prefixRange([]byte(contract)))
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

func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(idKey, bz)
	return id
}
