package keeper

import (
	"encoding/binary"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	//SetModuleAccount(sdk.Context, supply.ModuleAccountI)
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authexported.Account
}

func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	// todo: ensure in msg amounts are not negative and of type voucher
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	voucherDenom, err := types.AsVoucherDenom(totalAmount.Denom)
	if err != nil {
		return 0, err
	}

	if !k.HasCounterpartDenominator(ctx, voucherDenom) {
		return 0, sdkerrors.Wrap(types.ErrUnknown, "denominator")
	}

	// burn vouchers
	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, totalInVouchers)
	if err != nil {
		return 0, err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, totalInVouchers); err != nil {
		panic(err)
	}

	// persist TX in pool
	nextID := k.autoIncrementID(ctx, types.KeyLastTXPoolID)
	outgoing := types.OutgoingTx{
		//BridgeContractAddress: , // TODO: do we need to store this?
		Sender:      sender,
		DestAddress: counterpartReceiver,
		Amount:      amount,
		BridgeFee:   fee,
	}
	err = k.setPoolEntry(ctx, nextID, outgoing)
	if err != nil {
		return 0, err
	}

	// add a second index with the fee
	k.appendToUnbatchedTXIndex(ctx, fee, nextID)

	// todo: add second index for sender so that we can easily query: give pending Tx by sender
	// todo: what about a second index for receiver?
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
	idSet = append(idSet, txID)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(idSet))
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
	idSet = append([]uint64{txID}, idSet...)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(idSet))
}

func (k Keeper) RemoveFromUnbatchedTXIndex(ctx sdk.Context, fee sdk.Coin, txID uint64) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	bz := store.Get(idxKey)
	if bz == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "fee")
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	for i := range idSet {
		if idSet[i] == txID {
			idSet = append(idSet[0:i], idSet[i+1:]...)
			if len(idSet) != 0 {
				store.Set(idxKey, k.cdc.MustMarshalBinaryBare(idSet))
			} else {
				store.Delete(idxKey)
			}
			return nil
		}
	}
	return sdkerrors.Wrap(types.ErrUnknown, "tx id")
}

func (k Keeper) setPoolEntry(ctx sdk.Context, id uint64, val types.OutgoingTx) error {
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

func (k Keeper) GetCounterpartDenominator(ctx sdk.Context, voucherDenom types.VoucherDenom) (*types.BridgedDenominator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenominatorKey(voucherDenom.Unprefixed()))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "denominator")
	}
	var bridgedDenominator types.BridgedDenominator
	return &bridgedDenominator, k.cdc.UnmarshalBinaryBare(bz, &bridgedDenominator)
}

func (k Keeper) SetCounterpartDenominator(ctx sdk.Context, bridgeContractAddr, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	voucherDenominator := types.NewVoucherDenom(bridgeContractAddr, tokenID)
	bridgedDenominator := types.BridgedDenominator{
		BridgeContractAddress: bridgeContractAddr,
		TokenID:               tokenID,
	}
	store.Set(types.GetDenominatorKey(voucherDenominator.Unprefixed()), k.cdc.MustMarshalBinaryBare(bridgedDenominator))
}

func (k Keeper) HasCounterpartDenominator(ctx sdk.Context, voucherDenominator types.VoucherDenom) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDenominatorKey(voucherDenominator.Unprefixed()))
}

func (k Keeper) IterateOutgoingPoolByFee(ctx sdk.Context, voucherDenom types.VoucherDenom, cb func(uint64, types.OutgoingTx) bool) error {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.ReverseIterator(prefixRange([]byte(voucherDenom.Unprefixed())))
	for ; iter.Valid(); iter.Next() {
		var ids types.IDSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)
		// cb returns true to stop early
		for _, id := range ids {
			tx, err := k.getPoolEntry(ctx, id)
			if err != nil {
				return err
			}
			if cb(id, *tx) {
				return nil
			}
		}
	}
	return nil
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
