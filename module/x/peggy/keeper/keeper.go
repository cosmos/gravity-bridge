package keeper

import (
	"encoding/binary"
	"math"
	"sort"
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper types.StakingKeeper

	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace params.Subspace

	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	supplyKeeper types.SupplyKeeper

	AttestationHandler interface {
		Handle(sdk.Context, types.Attestation) error
	}
}

// NewKeeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace, stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper) Keeper {
	k := Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeKey:      storeKey,
		StakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
	}
	k.AttestationHandler = AttestationHandler{
		keeper:       k,
		supplyKeeper: supplyKeeper,
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return k
}

func (k Keeper) SetValsetRequest(ctx sdk.Context) types.Valset {
	store := ctx.KVStore(k.storeKey)
	valset := k.GetCurrentValset(ctx)
	nonce := ctx.BlockHeight()
	valset.Nonce = nonce
	store.Set(types.GetValsetRequestKey(nonce), k.cdc.MustMarshalBinaryBare(valset))

	event := sdk.NewEvent(
		types.EventTypeMultisigUpdateRequest,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyMultisigID, strconv.Itoa(int(nonce))),
		sdk.NewAttribute(types.AttributeKeyNonce, types.NonceFromUint64(uint64(nonce)).String()),
	)
	ctx.EventManager().EmitEvent(event)
	return valset
}

func (k Keeper) HasValsetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValsetRequestKey(int64(nonce))) // todo: revisit type where nonce is used as int64
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

// Iterate through all valset set requests in DESC order.
func (k Keeper) IterateValsetRequest(ctx sdk.Context, cb func(key []byte, val types.Valset) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValsetRequestKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var valset types.Valset
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &valset)
		// cb returns true to stop early
		if cb(iter.Key(), valset) {
			break
		}
	}
}

func (k Keeper) SetValsetConfirm(ctx sdk.Context, valsetConf types.MsgValsetConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValsetConfirmKey(valsetConf.Nonce, valsetConf.Validator), k.cdc.MustMarshalBinaryBare(valsetConf))
}

func (k Keeper) HasValsetConfirm(ctx sdk.Context, nonce int64, validatorAddr sdk.AccAddress) bool {
	// todo: param should be sdk.ValAddress instead
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValsetConfirmKey(nonce, validatorAddr))
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

// Iterate through all valset confirms for a nonce in ASC order
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

func (k Keeper) SetEthAddress(ctx sdk.Context, validator sdk.ValAddress, ethAddr types.EthereumAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthAddressKey(validator), ethAddr.Bytes())
}

func (k Keeper) GetEthAddress(ctx sdk.Context, validator sdk.ValAddress) *types.EthereumAddress {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.GetEthAddressKey(validator))
	if len(val) == 0 {
		return nil
	}
	addr := types.NewEthereumAddress(string(val))
	return &addr
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

// GetCurrentValset gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'Peggy power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the Peggy contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) GetCurrentValset(ctx sdk.Context) types.Valset {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethAddrs := make([]string, len(validators))
	powers := make([]int64, len(validators))
	var totalPower int64
	const maxUint32 int64 = math.MaxUint32
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for _, validator := range validators {
		validatorAddress := validator.GetOperator()
		p := k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress)
		totalPower += p
	}
	for i, validator := range validators {
		validatorAddress := validator.GetOperator()
		p := k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress) // TODO: avoid this Gas costs
		p = maxUint32 * p / totalPower
		//	safe math version:	p = sdk.NewInt(p).MulRaw(maxUint32).QuoRaw(totalPower).Int64()
		powers[i] = p
		ethAddr := k.GetEthAddress(ctx, validatorAddress)
		if ethAddr != nil {
			ethAddrs[i] = ethAddr.String()
		}
	}
	valset := types.Valset{EthAddresses: ethAddrs, Powers: powers}
	sort.Sort(valsetSort(valset))
	return valset
}

// GetParams returns the total set of wasm parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) setParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

func (k Keeper) GetBridgeContractAddress(ctx sdk.Context) types.EthereumAddress {
	var a types.EthereumAddress
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractAddress, &a)
	return a

}

func (k Keeper) GetBridgeChainID(ctx sdk.Context) uint64 {
	var a uint64
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractChainID, &a)
	return a
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
