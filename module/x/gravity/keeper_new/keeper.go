package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type AttestationHandler interface {
	Handle(sdk.Context, types.Attestation, types.EthereumClaim) error
}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace

	cdc            codec.BinaryMarshaler
	bankKeeper     types.BankKeeper
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper

	attestationHandler AttestationHandler
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, slashingKeeper types.SlashingKeeper,
	attestationHandler AttestationHandler,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                cdc,
		paramSpace:         paramSpace,
		storeKey:           storeKey,
		bankKeeper:         bankKeeper,
		slashingKeeper:     slashingKeeper,
		stakingKeeper:      stakingKeeper,
		attestationHandler: attestationHandler,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetEthAddress returns the eth address for a given gravity validatorAddr
func (k Keeper) GetEthAddress(ctx sdk.Context, validatorAddr sdk.ValAddress) common.Address {
	// TODO: use prefix store
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetEthAddressKey(validatorAddr))
	if len(bz) == 0 {
		// return zero address
		return common.Address{}
	}

	return common.BytesToAddress(bz)
}

// SetEthAddress sets the ethereum address for a given validator
func (k Keeper) SetEthAddress(ctx sdk.Context, validatorAddr sdk.ValAddress, ethereumAddr common.Address) {
	// TODO: use prefix store
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthAddressKey(validatorAddr), ethereumAddr.Bytes())
}

// GetOrchestratorValidator returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidator(ctx sdk.Context, orch sdk.AccAddress) sdk.ValAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOrchestratorAddressKey(orch))
	if len(bz) == 0 {
		return nil
	}

	return sdk.ValAddress(bz)
}

// SetOrchestratorValidator sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidator(ctx sdk.Context, val sdk.ValAddress, orch sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOrchestratorAddressKey(orch), val.Bytes())
}

// GetEthereumInfo returns the ethereum block height and timestamp of the last
// observed attestation.
func (k Keeper) GetEthereumInfo(ctx sdk.Context) (types.EthereumInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastObservedEthereumBlockHeightKey)
	if len(bz) == 0 {
		return types.EthereumInfo{}, false
	}

	ctx.BlockTime().UnixNano()

	var info types.EthereumInfo
	k.cdc.UnmarshalBinaryBare(bz, &info)
	return info, false
}

// SetEthereumInfo sets an observed ethereum block height and timestamp to the store.
func (k Keeper) SetEthereumInfo(ctx sdk.Context, info types.EthereumInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEthereumBlockHeightKey, k.cdc.MustMarshalBinaryBare(&info))
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)
	if len(bytes) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bytes)
}

// setLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(nonce))
}

func (k Keeper) GetTransferTx(ctx sdk.Context, id uint64) (types.TransferTx, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTransferTxPoolKey(id))
	if len(bz) == 0 {
		return types.TransferTx{}, false
	}

	var tx types.TransferTx
	k.cdc.UnmarshalBinaryBare(bz, &tx)
	return tx, true
}

func (k Keeper) SetTransferTx(ctx sdk.Context, id uint64, tx types.TransferTx) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&tx)
	store.Set(types.GetTransferTxPoolKey(id), bz)
}

func (k Keeper) DeleteTransferTx(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTransferTxPoolKey(id))
}

// IterateTransferTxs
func (k Keeper) IterateTransferTxs(ctx sdk.Context, cb func(id uint64, tx types.TransferTx) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransferTxPoolKey)

	iterator := prefixStore.ReverseIterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tx types.TransferTx
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &tx)
		id := sdk.BigEndianToUint64(iterator.Key()[:1]) // TODO: check correctness
		if cb(id, tx) {
			break // stop iteration
		}
	}
}

// GetTransferTxs returns all the outgoing transactions from the pool in desc order.
// TODO: create struct with ID and transferTx
func (k Keeper) GetTransferTxs(ctx sdk.Context) []types.TransferTx {
	txs := make([]types.TransferTx, 0)
	k.IterateTransferTxs(ctx, func(id uint64, tx types.TransferTx) bool {
		txs = append(txs, tx)
		return false
	})

	return txs
}
