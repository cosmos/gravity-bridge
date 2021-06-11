package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func (k Keeper) getCosmosOriginatedDenom(ctx sdk.Context, tokenContract string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MakeERC20ToDenomKey(tokenContract))

	if bz != nil {
		return string(bz), true
	}
	return "", false
}

func (k Keeper) getCosmosOriginatedERC20(ctx sdk.Context, denom string) (common.Address, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MakeDenomToERC20Key(denom))

	if bz != nil {
		return common.BytesToAddress(bz), true
	}
	return common.BytesToAddress([]byte{}), false
}

func (k Keeper) setCosmosOriginatedDenomToERC20(ctx sdk.Context, denom string, tokenContract string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MakeDenomToERC20Key(denom), common.HexToAddress(tokenContract).Bytes())
	store.Set(types.MakeERC20ToDenomKey(tokenContract), []byte(denom))
}

// DenomToERC20 returns (bool isCosmosOriginated, string ERC20, err)
// Using this information, you can see if an asset is native to Cosmos or Ethereum,
// and get its corresponding ERC20 address.
// This will return an error if it cant parse the denom as a gravity denom, and then also can't find the denom
// in an index of ERC20 contracts deployed on Ethereum to serve as synthetic Cosmos assets.
func (k Keeper) DenomToERC20Lookup(ctx sdk.Context, denom string) (bool, common.Address, error) {
	// First try parsing the ERC20 out of the denom
	tc1, err := types.GravityDenomToERC20(denom)
	if err != nil {
		// Look up ERC20 contract in index and error if it's not in there.
		tc2, exists := k.getCosmosOriginatedERC20(ctx, denom)
		if !exists {
			return false, common.Address{},
				fmt.Errorf("denom not a gravity voucher coin: %s, and also not in cosmos-originated ERC20 index", err)
		}
		// This is a cosmos-originated asset
		return true, tc2, nil
	}

	// This is an ethereum-originated asset
	return false, common.HexToAddress(tc1), nil
}

// ERC20ToDenom returns (bool isCosmosOriginated, string denom, err)
// Using this information, you can see if an ERC20 address represents an asset is native to Cosmos or Ethereum,
// and get its corresponding denom
func (k Keeper) ERC20ToDenomLookup(ctx sdk.Context, tokenContract string) (bool, string) {
	// First try looking up tokenContract in index
	dn1, exists := k.getCosmosOriginatedDenom(ctx, tokenContract)
	if exists {
		// It is a cosmos originated asset
		return true, dn1
	}

	// If it is not in there, it is not a cosmos originated token, turn the ERC20 into a gravity denom

	return false, types.NewERC20Token(0, tokenContract).GravityCoin().Denom
}

// iterateERC20ToDenom iterates over erc20 to denom relations
func (k Keeper) iterateERC20ToDenom(ctx sdk.Context, cb func([]byte, *types.ERC20ToDenom) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.ERC20ToDenomKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		erc20ToDenom := types.ERC20ToDenom{
			Erc20: string(iter.Key()),
			Denom: string(iter.Value()),
		}
		// cb returns true to stop early
		if cb(iter.Key(), &erc20ToDenom) {
			break
		}
	}
}
