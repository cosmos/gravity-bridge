package keeper

import (
	"encoding/binary"
	"sort"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k Keeper) SetValsetRequest(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	valset := k.GetCurrentValset(ctx)
	nonce := ctx.BlockHeight()
	valset.Nonce = nonce
	store.Set(types.GetValsetRequestKey(nonce), k.cdc.MustMarshalBinaryBare(valset))
}

func (k Keeper) GetValsetRequest(ctx sdk.Context, nonce int64) *types.Valset {
	store := ctx.KVStore(k.storeKey)

	store_bytes := store.Get(types.GetValsetRequestKey(nonce))
	if store_bytes == nil {
		return nil
	}
	var valset types.Valset
	k.cdc.MustUnmarshalBinaryBare(store_bytes, &valset)
	return &valset
}

func (k Keeper) SetValsetConfirm(ctx sdk.Context, valsetConf types.MsgValsetConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValsetConfirmKey(valsetConf.Nonce, valsetConf.Validator), k.cdc.MustMarshalBinaryBare(valsetConf))
}

func (k Keeper) GetValsetConfirm(ctx sdk.Context, nonce int64, validator sdk.AccAddress) *types.MsgValsetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetValsetConfirmKey(nonce, validator))
	if entity == nil {
		return nil
	}
	confirm := types.MsgValsetConfirm{}
	k.cdc.MustUnmarshalBinaryBare(entity, &confirm)
	return &confirm
}

func (k Keeper) IterateValsetConfirmByNonce(ctx sdk.Context, nonce int64, cb func([]byte, types.MsgValsetConfirm) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValsetConfirmKey)

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	iter := prefixStore.Iterator(prefixRange(nonceBytes))
	for ; iter.Valid(); iter.Next() {
		confirm := types.MsgValsetConfirm{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

func (k Keeper) SetEthAddress(ctx sdk.Context, validator sdk.AccAddress, ethAddr string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthAddressKey(validator), []byte(ethAddr))
}

func (k Keeper) GetEthAddress(ctx sdk.Context, validator sdk.AccAddress) string {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.GetEthAddressKey(validator))
	return string(val)
}

type valsetSort types.Valset

func (a valsetSort) Len() int { return len(a.EthAddresses) }
func (a valsetSort) Swap(i, j int) {
	a.EthAddresses[i], a.EthAddresses[j] = a.EthAddresses[j], a.EthAddresses[i]
	a.Powers[i], a.Powers[j] = a.Powers[j], a.Powers[i]
}
func (a valsetSort) Less(i, j int) bool {
	// Secondary sort on eth address in case powers are equal
	if a.Powers[i] == a.Powers[j] {
		return a.EthAddresses[i] < a.EthAddresses[j]
	}
	return a.Powers[i] < a.Powers[j]
}

func (k Keeper) GetCurrentValset(ctx sdk.Context) types.Valset {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethAddrs := make([]string, len(validators))
	powers := make([]int64, len(validators))
	for i, validator := range validators {
		validatorAddress := validator.GetOperator()
		p := k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress)
		powers[i] = p
		ethAddrs[i] = k.GetEthAddress(ctx, sdk.AccAddress(validatorAddress))
	}
	valset := types.Valset{EthAddresses: ethAddrs, Powers: powers}
	sort.Sort(valsetSort(valset))
	return valset
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// 		Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// 				 []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
//		Example: []byte{255, 255, 255, 255} becomes nil
//
func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}
