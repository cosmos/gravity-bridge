package keeper

import (
	"sort"

	"github.com/althea-net/peggy/module/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper types.StakingKeeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, stakingKeeper types.StakingKeeper) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		StakingKeeper: stakingKeeper,
	}
}

func (k Keeper) MakeValsetRequest(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	valset := k.GetValset(ctx)
	nonce := ctx.BlockHeight()
	store.Set(types.GetValsetRequestKey(nonce), k.cdc.MustMarshalBinaryBare(valset))
}

func (k Keeper) GetValsetRequest(ctx sdk.Context, nonce int64) types.Valset {
	store := ctx.KVStore(k.storeKey)

	valset := types.Valset{}
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.GetValsetRequestKey(nonce)), valset)
	return valset
}

func (k Keeper) SetValsetConfirm(ctx sdk.Context, valsetConf types.MsgValsetConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValsetConfirmKey(valsetConf.Nonce, valsetConf.Validator), k.cdc.MustMarshalBinaryBare(valsetConf))
}

func (k Keeper) SetEthAddress(ctx sdk.Context, validator sdk.AccAddress, ethAddr string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthAddressKey(validator), []byte(ethAddr))
}

func (k Keeper) GetEthAddress(ctx sdk.Context, validator sdk.AccAddress) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.GetEthAddressKey(validator)))
}

type valsetSort types.Valset

func (a valsetSort) Len() int { return len(a.EthAdresses) }
func (a valsetSort) Swap(i, j int) {
	a.EthAdresses[i], a.EthAdresses[j] = a.EthAdresses[j], a.EthAdresses[i]
	a.Powers[i], a.Powers[j] = a.Powers[j], a.Powers[i]
}
func (a valsetSort) Less(i, j int) bool {
	// Secondary sort on eth address in case powers are equal
	if a.Powers[i] == a.Powers[j] {
		return a.EthAdresses[i] < a.EthAdresses[j]
	}
	return a.Powers[i] < a.Powers[j]
}

func (k Keeper) GetValset(ctx sdk.Context) types.Valset {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethAddrs := make([]string, len(validators))
	powers := make([]int64, len(validators))
	for i, validator := range validators {
		validatorAddress := validator.GetOperator()
		p := k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress)
		powers[i] = p
		ethAddrs[i] = k.GetEthAddress(ctx, sdk.AccAddress(validatorAddress))
	}
	valset := types.Valset{EthAdresses: ethAddrs, Powers: powers}
	sort.Sort(valsetSort(valset))
	return valset
}

// // Gets the entire Whois metadata struct for a name
// func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
// 	store := ctx.KVStore(k.storeKey)

// 	if !k.IsNamePresent(ctx, name) {
// 		return types.NewWhois()
// 	}

// 	bz := store.Get([]byte(name))

// 	var whois types.Whois

// 	k.cdc.MustUnmarshalBinaryBare(bz, &whois)

// 	return whois
// }

// // Sets the entire Whois metadata struct for a name
// func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
// 	if whois.Owner.Empty() {
// 		return
// 	}

// 	store := ctx.KVStore(k.storeKey)

// 	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
// }

// // Deletes the entire Whois metadata struct for a name
// func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Delete([]byte(name))
// }

// // ResolveName - returns the string that the name resolves to
// func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
// 	return k.GetWhois(ctx, name).Value
// }

// // SetName - sets the value string that a name resolves to
// func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
// 	whois := k.GetWhois(ctx, name)
// 	whois.Value = value
// 	k.SetWhois(ctx, name, whois)
// }

// // SetName - sets the value string that a name resolves to
// func (k Keeper) SetEthAddress(ctx sdk.Context, name string, value string) {
// 	whois := k.GetWhois(ctx, name)
// 	whois.Value = value
// 	k.SetWhois(ctx, name, whois)
// }

// // HasOwner - returns whether or not the name already has an owner
// func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
// 	return !k.GetWhois(ctx, name).Owner.Empty()
// }

// // GetOwner - get the current owner of a name
// func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
// 	return k.GetWhois(ctx, name).Owner
// }

// // SetOwner - sets the current owner of a name
// func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
// 	whois := k.GetWhois(ctx, name)
// 	whois.Owner = owner
// 	k.SetWhois(ctx, name, whois)
// }

// // GetPrice - gets the current price of a name
// func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
// 	return k.GetWhois(ctx, name).Price
// }

// // SetPrice - sets the current price of a name
// func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
// 	whois := k.GetWhois(ctx, name)
// 	whois.Price = price
// 	k.SetWhois(ctx, name, whois)
// }

// // Get an iterator over all names in which the keys are the names and the values are the whois
// func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
// 	store := ctx.KVStore(k.storeKey)
// 	return sdk.KVStorePrefixIterator(store, nil)
// }

// // Check if the name is present in the store or not
// func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
// 	store := ctx.KVStore(k.storeKey)
// 	return store.Has([]byte(name))
// }
