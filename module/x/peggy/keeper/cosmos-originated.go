package keeper

import (
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCosmosOriginatedDenom(ctx sdk.Context, tokenContract string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetERC20ToDenomKey(tokenContract))

	if bz != nil {
		return string(bz), true
	}
	return "", false
}

func (k Keeper) GetCosmosOriginatedERC20(ctx sdk.Context, denom string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenomToERC20Key(denom))

	if bz != nil {
		return string(bz), true
	}
	return "", false
}

func (k Keeper) setCosmosOriginatedDenomToERC20(ctx sdk.Context, denom string, tokenContract string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDenomToERC20Key(denom), []byte(tokenContract))
	store.Set(types.GetERC20ToDenomKey(tokenContract), []byte(denom))
}

// DenomToERC20 returns (bool isCosmosOriginated, string ERC20, err)
// Using this information, you can see if an asset is native to Cosmos or Ethereum, and get its corresponding ERC20 address
// This will return an error if it cant parse the denom as a peggy denom, and then also can't find the denom
// in an index of ERC20 contracts deployed on Ethereum to serve as synthetic Cosmos assets.
func (k Keeper) DenomToERC20(ctx sdk.Context, denom string) (bool, string, error) {
	// First try parsing the ERC20 out of the denom
	tc1, err := types.PeggyDenomToERC20(denom)

	if err != nil {
		// Look up ERC20 contract in index and error if it's not in there.
		tc2, exists := k.GetCosmosOriginatedERC20(ctx, denom)
		if !exists {
			return false, "", fmt.Errorf("denom not a peggy voucher coin: %s, and also not in cosmos-originated ERC20 index", err)
		}
		// This is a cosmos-originated asset
		return true, tc2, nil
	} else {
		// This is an ethereum-originated asset
		return false, tc1, nil
	}
}

// ERC20ToDenom returns (bool isCosmosOriginated, string denom, err)
// Using this information, you can see if an ERC20 address represents an asset is native to Cosmos or Ethereum,
// and get its corresponding denom
func (k Keeper) ERC20ToDenom(ctx sdk.Context, tokenContract string) (bool, string) {
	// First try looking up tokenContract in index
	dn1, exists := k.GetCosmosOriginatedDenom(ctx, tokenContract)
	if exists {
		// It is a cosmos originated asset
		return true, dn1
	} else {
		// If it is not in there, it is not a cosmos originated token, turn the ERC20 into a peggy denom
		return false, types.PeggyDenom(tokenContract)
	}
}
