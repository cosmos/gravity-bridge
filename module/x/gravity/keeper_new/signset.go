package keeper

import (
	"math"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

// FIXME: rename to EthSignSet, clean up and write docs

// SetEthSignSetRequest returns a new instance of the Gravity BridgeValidatorSet
// i.e. {"nonce": 1, "memebers": [{"eth_addr": "foo", "power": 11223}]}
func (k Keeper) SetEthSignSetRequest(ctx sdk.Context) *types.EthSignSet {
	valset := k.GetCurrentEthSignSet(ctx)
	k.StoreEthSignSet(ctx, valset)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMultisigUpdateRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
			sdk.NewAttribute(types.AttributeKeyMultisigID, fmt.Sprint(valset.Nonce)),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(valset.Nonce)),
		),
	)

	return valset
}

// StoreEthSignSet is for storing a validator set at a given height
func (k Keeper) StoreEthSignSet(ctx sdk.Context, valset *types.EthSignSet) {
	store := ctx.KVStore(k.storeKey)
	valset.Height = uint64(ctx.BlockHeight())
	store.Set(types.GetEthSignSetKey(valset.Nonce), k.cdc.MustMarshalBinaryBare(valset))
	k.SetLatestEthSignSetNonce(ctx, valset.Nonce)
}

//  SetLatestEthSignSetNonce sets the latest valset nonce
func (k Keeper) SetLatestEthSignSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestEthSignSetNonce, types.UInt64Bytes(nonce))
}

// StoreEthSignSetUnsafe is for storing a validator set at a given height
func (k Keeper) StoreEthSignSetUnsafe(ctx sdk.Context, valset *types.EthSignSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthSignSetKey(valset.Nonce), k.cdc.MustMarshalBinaryBare(valset))
	k.SetLatestEthSignSetNonce(ctx, valset.Nonce)
}

// HasEthSignSetRequest returns true if a valset defined by a nonce exists
func (k Keeper) HasEthSignSetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetEthSignSetKey(nonce))
}

// DeleteEthSignSet deletes the valset at a given nonce from state
func (k Keeper) DeleteEthSignSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetEthSignSetKey(nonce))
}

// GetLatestEthSignSetNonce returns the latest valset nonce
func (k Keeper) GetLatestEthSignSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LatestEthSignSetNonce)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// GetEthSignSet returns a valset by nonce
func (k Keeper) GetEthSignSet(ctx sdk.Context, nonce uint64) *types.EthSignSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetEthSignSetKey(nonce))
	if bz == nil {
		return nil
	}
	var valset types.EthSignSet
	k.cdc.MustUnmarshalBinaryBare(bz, &valset)
	return &valset
}

// IterateEthSignSets retruns all valsetRequests
func (k Keeper) IterateEthSignSets(ctx sdk.Context, cb func(key []byte, val *types.EthSignSet) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.EthSignSetRequestKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var valset types.EthSignSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &valset)
		// cb returns true to stop early
		if cb(iter.Key(), &valset) {
			break
		}
	}
}

// GetEthSignSets returns all the validator sets in state
func (k Keeper) GetEthSignSets(ctx sdk.Context) (out []*types.EthSignSet) {
	k.IterateEthSignSets(ctx, func(_ []byte, val *types.EthSignSet) bool {
		out = append(out, val)
		return false
	})
	sort.Sort(types.EthSignSets(out))
	return
}

// GetLatestEthSignSet returns the latest validator set in state
func (k Keeper) GetLatestEthSignSet(ctx sdk.Context) (out *types.EthSignSet) {
	latestEthSignSetNonce := k.GetLatestEthSignSetNonce(ctx)
	out = k.GetEthSignSet(ctx, latestEthSignSetNonce)
	return
}

// setLastSlashedEthSignSetNonce sets the latest slashed valset nonce
func (k Keeper) SetLastSlashedEthSignSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedEthSignSetNonce, types.UInt64Bytes(nonce))
}

// GetLastSlashedEthSignSetNonce returns the latest slashed valset nonce
func (k Keeper) GetLastSlashedEthSignSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedEthSignSetNonce)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// SetLastUnBondingBlockHeight sets the last unbonding block height
func (k Keeper) SetLastUnBondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastUnBondingBlockHeight, types.UInt64Bytes(unbondingBlockHeight))
}

// GetLastUnBondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnBondingBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastUnBondingBlockHeight)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// GetUnSlashedEthSignSets returns all the unslashed validator sets in state
func (k Keeper) GetUnSlashedEthSignSets(ctx sdk.Context, maxHeight uint64) (out []*types.EthSignSet) {
	lastSlashedEthSignSetNonce := k.GetLastSlashedEthSignSetNonce(ctx)
	k.IterateEthSignSetBySlashedEthSignSetNonce(ctx, lastSlashedEthSignSetNonce, maxHeight, func(_ []byte, valset *types.EthSignSet) bool {
		if valset.Nonce > lastSlashedEthSignSetNonce {
			out = append(out, valset)
		}
		return false
	})
	return
}

// IterateEthSignSetBySlashedEthSignSetNonce iterates through all valset by last slashed valset nonce in ASC order
func (k Keeper) IterateEthSignSetBySlashedEthSignSetNonce(ctx sdk.Context, lastSlashedEthSignSetNonce uint64, maxHeight uint64, cb func([]byte, *types.EthSignSet) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.EthSignSetRequestKey)
	iter := prefixStore.Iterator(types.UInt64Bytes(lastSlashedEthSignSetNonce), types.UInt64Bytes(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var valset types.EthSignSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &valset)
		// cb returns true to stop early
		if cb(iter.Key(), &valset) {
			break
		}
	}
}

// TODO: fix

// GetCurrentEthSignSet gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'gravity power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the gravity contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) GetCurrentEthSignSet(ctx sdk.Context) *types.EthSignSet {
	validators := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	bridgeValidators := make([]*types.BridgeValidator, len(validators))
	var totalPower uint64
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for i, validator := range validators {
		val := validator.GetOperator()

		p := uint64(k.stakingKeeper.GetLastValidatorPower(ctx, val))
		totalPower += p

		bridgeValidators[i] = &types.BridgeValidator{Power: p}
		if ethAddr := k.GetEthAddress(ctx, val); ethAddr != common.Address{} {
			bridgeValidators[i].EthereumAddress = ethAddr
		}
	}
	// normalize power values
	for i := range bridgeValidators {
		bridgeValidators[i].Power = sdk.NewUint(bridgeValidators[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	// TODO: make the nonce an incrementing one (i.e. fetch last nonce from state, increment, set here)
	return types.NewEthSignSet(uint64(ctx.BlockHeight()), uint64(ctx.BlockHeight()), bridgeValidators)
}
