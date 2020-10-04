package keeper

import (
	"math"
	"sort"
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	store.Set(types.GetValsetRequestKey(valset.Nonce), k.cdc.MustMarshalBinaryBare(valset))

	event := sdk.NewEvent(
		types.EventTypeMultisigUpdateRequest,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyMultisigID, valset.Nonce.String()),
		sdk.NewAttribute(types.AttributeKeyNonce, valset.Nonce.String()),
	)
	ctx.EventManager().EmitEvent(event)
	return valset
}

func (k Keeper) SetBootstrapValset(ctx sdk.Context, nonce types.UInt64Nonce, valset types.Valset) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValsetRequestKey(nonce)
	if store.Has(key) {
		return sdkerrors.Wrap(types.ErrDuplicate, "nonce")
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(valset))

	event := sdk.NewEvent(
		types.EventTypeMultisigUpdateRequest,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyMultisigID, nonce.String()),
		sdk.NewAttribute(types.AttributeKeyNonce, nonce.String()),
	)
	ctx.EventManager().EmitEvent(event)
	return nil
}

func (k Keeper) HasValsetRequest(ctx sdk.Context, nonce types.UInt64Nonce) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValsetRequestKey(nonce))
}

func (k Keeper) GetValsetRequest(ctx sdk.Context, nonce types.UInt64Nonce) *types.Valset {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValsetRequestKey(nonce))
	if bz == nil {
		return nil
	}
	var valset types.Valset
	k.cdc.MustUnmarshalBinaryBare(bz, &valset)
	return &valset
}

// Iterate through all valset set requests in DESC order.
func (k Keeper) IterateValsetRequest(ctx sdk.Context, cb func(key []byte, val types.Valset) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValsetRequestKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var valset types.Valset
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &valset)
		// cb returns true to stop early
		if cb(iter.Key(), valset) {
			break
		}
	}
}

// deprecated use SetBridgeApprovalSignature instead
func (k Keeper) SetValsetConfirm(ctx sdk.Context, valsetConf types.MsgValsetConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValsetConfirmKey(valsetConf.Nonce, valsetConf.Validator), k.cdc.MustMarshalBinaryBare(valsetConf))
}

func (k Keeper) SetBridgeApprovalSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, signature []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeApprovalSignatureKey(claimType, nonce, validator), signature)
}

func (k Keeper) GetBridgeApprovalSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetBridgeApprovalSignatureKey(claimType, nonce, validator))
}

func (k Keeper) HasBridgeApprovalSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeApprovalSignatureKey(claimType, nonce, validator))
}

func (k Keeper) IterateBridgeApprovalSignatures(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, cb func(_ []byte, sig []byte) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetBridgeApprovalSignatureKeyPrefix(claimType))
	iter := prefixStore.Iterator(prefixRange(nonce.Bytes()))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(iter.Key(), iter.Value()) {
			break
		}
	}
}

// deprecated use GetBridgeApprovalSignature
func (k Keeper) GetValsetConfirm(ctx sdk.Context, nonce types.UInt64Nonce, validator sdk.AccAddress) *types.MsgValsetConfirm {
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
// deprecated
func (k Keeper) IterateValsetConfirmByNonce(ctx sdk.Context, nonce types.UInt64Nonce, cb func([]byte, types.MsgValsetConfirm) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValsetConfirmKey)
	iter := prefixStore.Iterator(prefixRange(nonce.Bytes()))
	defer iter.Close()

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
		return a.EthAddresses[i].LessThan(a.EthAddresses[j])
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
	ethAddrs := make([]types.EthereumAddress, len(validators))
	powers := make([]uint64, len(validators))
	var totalPower uint64
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for i, validator := range validators {
		validatorAddress := validator.GetOperator()
		ethAddr := k.GetEthAddress(ctx, validatorAddress)
		if ethAddr != nil {
			ethAddrs[i] = *ethAddr
		}

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress))
		totalPower += p
		powers[i] = p
	}
	// normalize power values
	for i := range powers {
		powers[i] = sdk.NewUint(powers[i]).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}
	valset := types.Valset{
		EthAddresses: ethAddrs,
		Powers:       powers,
		Nonce:        types.NewUInt64Nonce(uint64(ctx.BlockHeight())),
	}
	sort.Sort(valsetSort(valset))
	return valset
}

// GetParams returns the total set of wasm parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
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

func (k Keeper) GetPeggyID(ctx sdk.Context) []byte {
	var a []byte
	k.paramSpace.Get(ctx, types.ParamsStoreKeyPeggyID, &a)
	return a
}
func (k Keeper) setPeggyID(ctx sdk.Context, v []byte) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyPeggyID, v)
}

func (k Keeper) GetStartThreshold(ctx sdk.Context) uint64 {
	var a uint64
	k.paramSpace.Get(ctx, types.ParamsStoreKeyStartThreshold, &a)
	return a
}

func (k Keeper) setStartThreshold(ctx sdk.Context, v uint64) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyStartThreshold, v)
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
