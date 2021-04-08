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
	attestationHandler AttestationHandler) Keeper {
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

func (k Keeper) GetOutgoingTx(ctx sdk.Context, id uint64) (types.OutgoingTransferTx, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxPoolKey(id))
	if len(bz) == 0 {
		return types.OutgoingTransferTx{}, false
	}

	var tx types.OutgoingTransferTx
	k.cdc.UnmarshalBinaryBare(bz, &tx)
	return tx, true
}

func (k Keeper) SetOutgoingTx(ctx sdk.Context, id uint64, tx types.OutgoingTransferTx) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&tx)
	store.Set(types.GetOutgoingTxPoolKey(id), bz)
}

func (k Keeper) DeleteOutgoingTx(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxPoolKey(id))
}

// IterateOutgoingTxs
func (k Keeper) IterateOutgoingTxs(ctx sdk.Context, cb func(id uint64, tx types.OutgoingTransferTx) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxPoolKey)

	iterator := prefixStore.ReverseIterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tx types.OutgoingTransferTx
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &tx)
		id := sdk.BigEndianToUint64(iterator.Key()[:1]) // TODO: check correctness
		if cb(id, tx) {
			break // stop iteration
		}
	}
}

// GetOutgoingTxs returns all the outgoing transactions from the pool in desc order.
// TODO: create struct with ID and transferTx
func (k Keeper) GetOutgoingTxs(ctx sdk.Context) []types.OutgoingTransferTx {
	txs := make([]types.OutgoingTransferTx, 0)
	k.IterateOutgoingTxs(ctx, func(id uint64, tx types.OutgoingTransferTx) bool {
		txs = append(txs, tx)
		return false
	})

	return txs
}
