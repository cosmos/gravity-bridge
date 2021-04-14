package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// FIXME: clean up and write docs

// SetEthSignerSetRequest returns a new instance of the bridge ethereum signer set
func (k Keeper) SetEthSignerSetRequest(ctx sdk.Context) types.EthSignerSet {
	signerSet := k.GetCurrentEthSignerSet(ctx)
	k.StoreEthSignerSet(ctx, signerSet)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMultisigUpdateRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyMultisigID, fmt.Sprint(signerSet.Nonce)),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(signerSet.Nonce)),
		),
	)

	return signerSet
}

// StoreEthSignerSet is for storing a validator set at a given height
func (k Keeper) StoreEthSignerSet(ctx sdk.Context, signerSet types.EthSignerSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthSignerSetKey(signerSet.Nonce), k.cdc.MustMarshalBinaryBare(signerSet))
	k.SetLatestEthSignerSetNonce(ctx, signerSet.Nonce)
}

// SetLatestEthSignerSetNonce sets the latest signerSet nonce
func (k Keeper) SetLatestEthSignerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestEthSignerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// DeleteEthSignerSet deletes the signerSet at a given nonce from state
func (k Keeper) DeleteEthSignerSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetEthSignerSetKey(nonce))
}

// GetLatestEthSignerSetNonce returns the latest signerSet nonce
func (k Keeper) GetLatestEthSignerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LatestEthSignerSetNonce)
	if len(bytes) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bytes)
}

// GetEthSignerSet returns a signerSet by nonce
func (k Keeper) GetEthSignerSet(ctx sdk.Context, nonce uint64) (types.EthSignerSet, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetEthSignerSetKey(nonce))
	if bz == nil {
		return types.EthSignerSet{}, false
	}
	var signerSet types.EthSignerSet
	k.cdc.MustUnmarshalBinaryBare(bz, &signerSet)

	return signerSet, true
}

// GetLatestEthSignerSet returns the latest validator set in state
func (k Keeper) GetLatestEthSignerSet(ctx sdk.Context) (out *types.EthSignerSet) {
	latestEthSignerSetNonce := k.GetLatestEthSignerSetNonce(ctx)
	out = k.GetEthSignerSet(ctx, latestEthSignerSetNonce)
	return
}

// setLastSlashedEthSignerSetNonce sets the latest slashed signerSet nonce
func (k Keeper) SetLastSlashedEthSignerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedEthSignerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedEthSignerSetNonce returns the latest slashed signerSet nonce
func (k Keeper) GetLastSlashedEthSignerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedEthSignerSetNonce)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// SetLastUnBondingBlockHeight sets the last unbonding block height
func (k Keeper) SetLastUnBondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastUnBondingBlockHeight, sdk.Uint64ToBigEndian(unbondingBlockHeight))
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

// GetUnSlashedEthSignerSets returns all the unslashed validator sets in state
func (k Keeper) GetUnSlashedEthSignerSets(ctx sdk.Context, maxHeight uint64) (out []*types.EthSignerSet) {
	lastSlashedEthSignerSetNonce := k.GetLastSlashedEthSignerSetNonce(ctx)
	k.IterateEthSignerSetBySlashedEthSignerSetNonce(ctx, lastSlashedEthSignerSetNonce, maxHeight, func(_ []byte, signerSet *types.EthSignerSet) bool {
		if signerSet.Nonce > lastSlashedEthSignerSetNonce {
			out = append(out, signerSet)
		}
		return false
	})
	return
}

// IterateEthSignerSetBySlashedEthSignerSetNonce iterates through all signerSet by last slashed signerSet nonce in ASC order
func (k Keeper) IterateEthSignerSetBySlashedEthSignerSetNonce(ctx sdk.Context, lastSlashedEthSignerSetNonce uint64, maxHeight uint64, cb func([]byte, *types.EthSignerSet) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.EthSignerSetRequestKey)
	iter := prefixStore.Iterator(sdk.Uint64ToBigEndian(lastSlashedEthSignerSetNonce), sdk.Uint64ToBigEndian(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var signerSet types.EthSignerSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &signerSet)
		// cb returns true to stop early
		if cb(iter.Key(), &signerSet) {
			break
		}
	}
}

// TODO: fix doc

// GetCurrentEthSignerSet gets powers from the store and normalizes them
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
func (k Keeper) GetCurrentEthSignerSet(ctx sdk.Context) types.EthSignerSet {
	signers := make([]types.EthSigner, 0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		// TODO: Remove this query. It doesn't make any sense to store the address separated from the power
		ethereumAddr := k.GetEthAddress(ctx, validator.GetOperator())

		signer := types.EthSigner{
			EthereumAddress: ethereumAddr.String(),
			Power:           validator.GetConsensusPower(), // TODO: be explicit that this is just the value not the %
		}

		signers = append(signers, signer)
		return false
	})

	return types.NewSignerSet(uint64(ctx.BlockHeight()), signers...)
}
