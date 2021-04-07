package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func (k Keeper) GetCoinDenomFromERC20Contract(ctx sdk.Context, tokenContract common.Address) (string, bool) {
	// TODO: prefix store
	store := ctx.KVStore(k.storeKey)

	// TODO: use address bytes instead
	bz := store.Get(types.GetERC20ToDenomKey(tokenContract.String()))
	if len(bz) == 0 {
		return "", false
	}

	return string(bz), true
}

func (k Keeper) GetERC20ContractFromCoinDenom(ctx sdk.Context, denom string) (common.Address, bool) {
	// TODO: prefix store
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenomToERC20Key(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}

	tokenContract := common.HexToAddress(string(bz))

	return tokenContract, true
}

func (k Keeper) setERC20DenomMap(ctx sdk.Context, denom string, tokenContract common.Address) {
	// TODO: prefix store
	store := ctx.KVStore(k.storeKey)
	contractHex := tokenContract.String()
	// TODO: use contract address bytes
	store.Set(types.GetDenomToERC20Key(denom), []byte(contractHex))
	store.Set(types.GetERC20ToDenomKey(contractHex), []byte(denom))
}

// IterateERC20ToDenom iterates over erc20 to denom relations
func (k Keeper) IterateERC20ToDenom(ctx sdk.Context, cb func(types.ERC20ToDenom) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ERC20ToDenomKey)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		erc20ToDenom := types.ERC20ToDenom{
			Erc20: string(iter.Key()),
			Denom: string(iter.Value()),
		}

		// cb returns true to stop early
		if cb(erc20ToDenom) {
			break
		}
	}
}
