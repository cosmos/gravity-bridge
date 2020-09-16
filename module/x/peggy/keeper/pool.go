package keeper

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
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

const voucherPrefixLen = len(types.VoucherDenomPrefix + types.DenomSeparator)

func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	store := ctx.KVStore(k.storeKey)

	// check and burn vouchers

	// todo: ensure in msg amounts are not negative and of type voucher
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	if !strings.HasPrefix(totalAmount.Denom, types.VoucherDenomPrefix) || len(totalAmount.Denom) != types.VoucherDenomLen {
		return 0, sdkerrors.Wrapf(types.ErrInvalid, "not a peggy denominator: %d", len(totalAmount.Denom))
	}
	// no unique separator in this sdk version possible :-(
	//parts := strings.Split(totalAmount.Denom, types.DenomSeparator)
	//if len(parts) != 2 || parts[0] != types.VoucherDenomPrefix {
	//	return 0, sdkerrors.Wrap(types.ErrInvalid, "not a peggy denominator")
	//}
	unprefixedVoucherDenom := totalAmount.Denom[voucherPrefixLen:]
	if !k.HasCounterpartDenominator(ctx, unprefixedVoucherDenom) {
		return 0, sdkerrors.Wrap(types.ErrUnknown, "denominator")
	}

	counterpartDenom, err := k.GetCounterpartDenominator(ctx, unprefixedVoucherDenom)
	if err != nil {
		return 0, err
	}

	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, totalInVouchers)
	if err != nil {
		return 0, err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, totalInVouchers); err != nil {
		panic(err)
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastTXPoolID)
	// add to pool
	outgoing := types.OutgoingTx{
		//BridgeContractAddress: , // TODO: do we need to store this?
		Sender:      sender,
		DestAddress: counterpartReceiver,
		Amount:      types.AsTransferCoin(*counterpartDenom, amount),
		BridgeFee:   types.AsTransferCoin(*counterpartDenom, fee),
	}
	err = k.setPoolEntry(ctx, nextID, outgoing)
	if err != nil {
		return 0, err
	}

	// add a second index with the fee
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	if store.Has(idxKey) {
		bz := store.Get(idxKey)
		k.cdc.UnmarshalBinaryBare(bz, &idSet)
	}
	idSet = append(idSet, nextID)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(idSet))
	return nextID, nil
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

func (k Keeper) GetCounterpartDenominator(ctx sdk.Context, voucherDenom string) (*types.BridgedDenominator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenominatorKey(voucherDenom))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "denominator")
	}
	var bridgedDenominator types.BridgedDenominator
	return &bridgedDenominator, k.cdc.UnmarshalBinaryBare(bz, &bridgedDenominator)
}

func (k Keeper) SetCounterpartDenominator(ctx sdk.Context, bridgeContractAddr, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	voucherDenominator := toVoucherDenominator(bridgeContractAddr, tokenID)
	bridgedDenominator := types.BridgedDenominator{
		BridgeContractAddress: bridgeContractAddr,
		TokenID:               tokenID,
	}
	store.Set(types.GetDenominatorKey(voucherDenominator[voucherPrefixLen:]), k.cdc.MustMarshalBinaryBare(bridgedDenominator))
}

func (k Keeper) HasCounterpartDenominator(ctx sdk.Context, voucherDenominator string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDenominatorKey(voucherDenominator))
}

func (k Keeper) IterateOutgoingPoolByFee(ctx sdk.Context, cb func(uint64, types.OutgoingTx) bool) error {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.ReverseIterator(nil, nil)
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

func toVoucherDenominator(contractAddr, token string) string {
	denomTrace := fmt.Sprintf("%s/%s/", contractAddr, token)
	var hash tmbytes.HexBytes = tmhash.Sum([]byte(denomTrace))
	simpleVoucherDenum := types.VoucherDenomPrefix + types.DenomSeparator + hash.String()
	sdkVersionHackDenum := strings.ToLower(simpleVoucherDenum[0:15]) // todo: up to 15 chars (lowercase) allowed in this sdk version only
	return sdkVersionHackDenum
}
