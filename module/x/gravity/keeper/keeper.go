package keeper

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper types.StakingKeeper

	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace paramtypes.Subspace

	cdc            codec.BinaryMarshaler // The wire codec for binary encoding/decoding.
	bankKeeper     types.BankKeeper
	SlashingKeeper types.SlashingKeeper

	EthereumEventProcessor interface {
		Handle(sdk.Context, types.EthereumEvent) error
	}
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, slashingKeeper types.SlashingKeeper) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		cdc:            cdc,
		paramSpace:     paramSpace,
		storeKey:       storeKey,
		StakingKeeper:  stakingKeeper,
		bankKeeper:     bankKeeper,
		SlashingKeeper: slashingKeeper,
	}
	k.EthereumEventProcessor = EthereumEventProcessor{
		keeper:     k,
		bankKeeper: bankKeeper,
	}

	return k
}

/////////////////////////////
//     VALSET REQUESTS     //
/////////////////////////////

// SetSignerSetTxRequest returns a new instance of the Gravity BridgeValidatorSet
// i.e. {"nonce": 1, "memebers": [{"eth_addr": "foo", "power": 11223}]}
func (k Keeper) SetSignerSetTxRequest(ctx sdk.Context) *types.SignerSetTx {
	valset := k.GetCurrentSignerSetTx(ctx)
	k.StoreSignerSetTx(ctx, valset)

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

// StoreSignerSetTx is for storing a validator set at a given height
func (k Keeper) StoreSignerSetTx(ctx sdk.Context, valset *types.SignerSetTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSignerSetTxKey(valset.Nonce), k.cdc.MustMarshalBinaryBare(valset))
	k.SetLatestSignerSetTxNonce(ctx, valset.Nonce)
}

// SetLatestSignerSetTxNonce sets the latest valset nonce
func (k Keeper) SetLatestSignerSetTxNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestSignerSetTxNonce}, types.UInt64Bytes(nonce))
}

// StoreSignerSetTxUnsafe is for storing a valiator set at a given height
func (k Keeper) StoreSignerSetTxUnsafe(ctx sdk.Context, ussTx *types.SignerSetTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSignerSetTxKey(ussTx.Nonce), k.cdc.MustMarshalBinaryBare(ussTx))
	k.SetLatestSignerSetTxNonce(ctx, ussTx.Nonce)
}

// HasValsetRequest returns true if a valset defined by a nonce exists
func (k Keeper) HasValsetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSignerSetTxKey(nonce))
}

// DeleteValset deletes the valset at a given nonce from state
func (k Keeper) DeleteValset(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetSignerSetTxKey(nonce))
}

// GetLatestValsetNonce returns the latest valset nonce
func (k Keeper) GetLatestValsetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte{types.LatestSignerSetTxNonce})

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetSignerSetTx returns a valset by nonce
func (k Keeper) GetSignerSetTx(ctx sdk.Context, nonce uint64) *types.SignerSetTx {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSignerSetTxKey(nonce))
	if bz == nil {
		return nil
	}
	var ussTx types.SignerSetTx
	k.cdc.MustUnmarshalBinaryBare(bz, &ussTx)
	return &ussTx
}

// IterateSignerSetTxs retruns all valsetRequests
func (k Keeper) IterateSignerSetTxs(ctx sdk.Context, cb func(key []byte, ussTx *types.SignerSetTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.SignerSetTxKey})
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ussTx types.SignerSetTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ussTx)
		// cb returns true to stop early
		if cb(iter.Key(), &ussTx) {
			break
		}
	}
}

// GetSignerSetTxs returns all the validator sets in state
func (k Keeper) GetSignerSetTxs(ctx sdk.Context) (out []*types.SignerSetTx) {
	k.IterateSignerSetTxs(ctx, func(_ []byte, val *types.SignerSetTx) bool {
		out = append(out, val)
		return false
	})
	sort.Sort(types.SignerSetTxs(out))
	return
}

// GetLatestSignerSetTx returns the latest validator set in state
func (k Keeper) GetLatestSignerSetTx(ctx sdk.Context) (out *types.SignerSetTx) {
	latestValsetNonce := k.GetLatestValsetNonce(ctx)
	out = k.GetSignerSetTx(ctx, latestValsetNonce)
	return
}

// SetLastSlashedValsetNonce sets the latest slashed valset nonce
func (k Keeper) SetLastSlashedValsetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LastSlashedValsetNonce}, types.UInt64Bytes(nonce))
}

// GetLastSlashedValsetNonce returns the latest slashed valset nonce
func (k Keeper) GetLastSlashedValsetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte{types.LastSlashedValsetNonce})

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// SetLastUnBondingBlockHeight sets the last unbonding block height
func (k Keeper) SetLastUnBondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LastUnBondingBlockHeight}, types.UInt64Bytes(unbondingBlockHeight))
}

// GetLastUnBondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnBondingBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte{types.LastUnBondingBlockHeight})

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetUnSlashedValsets returns all the unslashed validator sets in state
func (k Keeper) GetUnSlashedValsets(ctx sdk.Context, maxHeight uint64) (out []*types.SignerSetTx) {
	lastSlashedValsetNonce := k.GetLastSlashedValsetNonce(ctx)
	k.IterateValsetBySlashedValsetNonce(ctx, lastSlashedValsetNonce, maxHeight, func(_ []byte, valset *types.SignerSetTx) bool {
		if valset.Nonce > lastSlashedValsetNonce {
			out = append(out, valset)
		}
		return false
	})
	return
}

// IterateValsetBySlashedValsetNonce iterates through all valset by last slashed valset nonce in ASC order
func (k Keeper) IterateValsetBySlashedValsetNonce(ctx sdk.Context, lastSlashedValsetNonce uint64, maxHeight uint64, cb func([]byte, *types.SignerSetTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.SignerSetTxKey})
	iter := prefixStore.Iterator(types.UInt64Bytes(lastSlashedValsetNonce), types.UInt64Bytes(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var valset types.SignerSetTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &valset)
		// cb returns true to stop early
		if cb(iter.Key(), &valset) {
			break
		}
	}
}

///////////////////////////////
//     ETHEREUM SIGNATURES   //
///////////////////////////////

// GetEthereumSignature returns a valset confirmation by a nonce and validator address
func (k Keeper) GetEthereumSignature(ctx sdk.Context, storeIndex []byte, validator sdk.ValAddress) hexutil.Bytes {
	return ctx.KVStore(k.storeKey).Get(types.GetEthereumSignatureKey(storeIndex, validator))
}

// SetEthereumSignature sets a valset confirmation
func (k Keeper) SetEthereumSignature(ctx sdk.Context, sig types.EthereumSignature, val sdk.ValAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	key := append(sig.GetStoreIndex(), val.Bytes()...)
	store.Set(key, sig.GetSignature())
	return key
}

// GetSignerSetTxSignatures returns all validator set confirmations by nonce
func (k Keeper) GetSignerSetTxSignatures(ctx sdk.Context, nonce uint64) (confirms []*types.SignerSetTxSignature) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.EthereumSignatureKey})
	start, end := prefixRange(types.UInt64Bytes(nonce))
	iterator := prefixStore.Iterator(start, end)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		confirm := types.SignerSetTxSignature{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &confirm)
		confirms = append(confirms, &confirm)
	}

	return confirms
}

// IterateValsetConfirmByNonce iterates through all valset confirms by nonce in ASC order
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
// TODO: specify which nonce this is
func (k Keeper) IterateValsetConfirmByNonce(ctx sdk.Context, nonce uint64, cb func([]byte, types.MsgSubmitEthereumSignature) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.EthereumSignatureKey})
	iter := prefixStore.Iterator(prefixRange(types.UInt64Bytes(nonce)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.MsgSubmitEthereumSignature{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

/////////////////////////////
//      BATCH CONFIRMS     //
/////////////////////////////

// GetBatchConfirm returns a batch confirmation given its nonce, the token contract, and a validator address
func (k Keeper) GetBatchConfirm(ctx sdk.Context, nonce uint64, tokenContract string, validator sdk.ValAddress) *types.BatchTxSignature {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetBatchTxSignatureKey(tokenContract, nonce, validator))
	if entity == nil {
		return nil
	}
	confirm := types.BatchTxSignature{}
	k.cdc.MustUnmarshalBinaryBare(entity, &confirm)
	return &confirm
}

// SetBatchConfirm sets a batch confirmation by a validator
func (k Keeper) SetBatchConfirm(ctx sdk.Context, batch *types.BatchTxSignature) []byte {
	store := ctx.KVStore(k.storeKey)
	acc, err := sdk.AccAddressFromBech32(batch.Orchestrator)
	if err != nil {
		panic(err)
	}
	key := types.GetBatchTxSignatureKey(batch.TokenContract, batch.Nonce, acc)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))
	return key
}

// IterateBatchConfirmByNonceAndTokenContract iterates through all batch confirmations
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
// TODO: specify which nonce this is
func (k Keeper) IterateBatchConfirmByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string, cb func([]byte, types.MsgConfirmBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.BatchTxSignatureKey})
	prfx := append(common.HexToAddress(tokenContract).Bytes(), types.UInt64Bytes(nonce)...)
	iter := prefixStore.Iterator(prefixRange(prfx))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.BatchTxSignature{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

// GetBatchTxSignatureByNonceAndTokenContract returns the batch confirms
func (k Keeper) GetBatchTxSignatureByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string) (out []types.BatchTxSignature) {
	k.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, msg types.BatchTxSignature) bool {
		out = append(out, msg)
		return false
	})
	return
}

/////////////////////////////
//    ADDRESS DELEGATION   //
/////////////////////////////

// SetOrchestratorValidator sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidator(ctx sdk.Context, val sdk.ValAddress, orch sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOrchestratorAddressKey(orch), val.Bytes())
}

// GetOrchestratorValidator returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidator(ctx sdk.Context, orch sdk.AccAddress) sdk.ValAddress {
	store := ctx.KVStore(k.storeKey)
	return sdk.ValAddress(store.Get(types.GetOrchestratorAddressKey(orch)))
}

/////////////////////////////
//       ETH ADDRESS       //
/////////////////////////////

// SetEthAddress sets the ethereum address for a given validator
func (k Keeper) SetEthAddress(ctx sdk.Context, validator sdk.ValAddress, ethAddr string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthereumAddressKey(validator), []byte(ethAddr))
}

// GetEthAddress returns the eth address for a given gravity validator
func (k Keeper) GetEthAddress(ctx sdk.Context, validator sdk.ValAddress) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.GetEthereumAddressKey(validator)))
}

// GetCurrentSignerSetTx gets powers from the store and normalizes them
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
func (k Keeper) GetCurrentSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethereumSigners := make([]*types.EthereumSigner, len(validators))
	var totalPower uint64
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for i, validator := range validators {
		val := validator.GetOperator()

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, val))
		totalPower += p

		ethereumSigners[i] = &types.EthereumSigner{Power: p}
		if ethAddr := k.GetEthAddress(ctx, val); ethAddr != "" {
			ethereumSigners[i].EthereumAddress = ethAddr
		}
	}
	// normalize power values
	for i := range ethereumSigners {
		ethereumSigners[i].Power = sdk.NewUint(ethereumSigners[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	// TODO: make the nonce an incrementing one (i.e. fetch last nonce from state, increment, set here)
	return types.NewValset(uint64(ctx.BlockHeight()), uint64(ctx.BlockHeight()), ethereumSigners)
}

/////////////////////////////
//       LOGICCALLS        //
/////////////////////////////

// GetContractCallTx gets an outgoing logic call
func (k Keeper) GetContractCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) *types.ContractCallTx {
	store := ctx.KVStore(k.storeKey)
	call := types.ContractCallTx{}
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.GetContractCallTxKey(invalidationID, invalidationNonce)), &call)
	return &call
}

// SetOutogingLogicCall sets an outgoing logic call
func (k Keeper) SetContractCallTx(ctx sdk.Context, call *types.ContractCallTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractCallTxKey(call.InvalidationScope, call.InvalidationNonce),
		k.cdc.MustMarshalBinaryBare(call))
}

// DeleteContractCallTx deletes outgoing logic calls
func (k Keeper) DeleteContractCallTx(ctx sdk.Context, invalidationID []byte, invalidationNonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetContractCallTxKey(invalidationID, invalidationNonce))
}

// IterateContractCallTxs iterates over outgoing logic calls
func (k Keeper) IterateContractCallTxs(ctx sdk.Context, cb func([]byte, *types.ContractCallTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.ContractCallTxKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var call types.ContractCallTx
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &call)
		// cb returns true to stop early
		if cb(iter.Key(), &call) {
			break
		}
	}
}

// GetContractCallTxs returns the outgoing tx batches
func (k Keeper) GetContractCallTxs(ctx sdk.Context) (out []*types.ContractCallTx) {
	k.IterateContractCallTxs(ctx, func(_ []byte, call *types.ContractCallTx) bool {
		out = append(out, call)
		return false
	})
	return
}

// CancelContractCallTxs releases all TX in the batch and deletes the batch
func (k Keeper) CancelContractCallTx(ctx sdk.Context, invalidationId []byte, invalidationNonce uint64) error {
	call := k.GetContractCallTx(ctx, invalidationId, invalidationNonce)
	if call == nil {
		return types.ErrUnknown
	}
	// Delete batch since it is finished
	k.DeleteContractCallTx(ctx, call.InvalidationScope, call.InvalidationNonce)

	// a consuming application will have to watch for this event and act on it
	batchEvent := sdk.NewEvent(
		types.EventTypeContractCallTxCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyInvalidationID, fmt.Sprint(call.InvalidationScope)),
		sdk.NewAttribute(types.AttributeKeyInvalidationNonce, fmt.Sprint(call.InvalidationNonce)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

/////////////////////////////
//       LOGICCONFIRMS     //
/////////////////////////////

// SetContractCallTxSignature sets a logic confirm in the store
func (k Keeper) SetContractCallTxSignature(ctx sdk.Context, msg *types.ContractCallTxSignature) {
	val, err := sdk.ValAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(k.storeKey).
		Set(types.GetContractCallTxSignatureKey(msg.InvalidationScope, msg.InvalidationNonce, val), k.cdc.MustMarshalBinaryBare(msg))
}

// GetContractCallTxSignature gets a logic confirm from the store
func (k Keeper) GetContractCallTxSignature(ctx sdk.Context, invalidationId []byte, invalidationNonce uint64, val sdk.AccAddress) *types.ContractCallTxSignature {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetContractCallTxSignatureKey(invalidationId, invalidationNonce, val))
	if data == nil {
		return nil
	}
	out := types.ContractCallTxSignature{}
	k.cdc.MustUnmarshalBinaryBare(data, &out)
	return &out
}

// DeleteContractCallTxSignature deletes a logic confirm from the store
func (k Keeper) DeleteContractCallTxSignature(
	ctx sdk.Context,
	invalidationID []byte,
	invalidationNonce uint64,
	val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetContractCallTxSignatureKey(invalidationID, invalidationNonce, val))
}

// IterateContractCallTxSignatureByInvalidationIDAndNonce iterates over all logic confirms stored by nonce
func (k Keeper) IterateContractCallTxSignatureByInvalidationIDAndNonce(
	ctx sdk.Context,
	invalidationID []byte,
	invalidationNonce uint64,
	cb func([]byte, *types.ContractCallTxSignature) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.ContractCallTxSignatureKey})
	iter := prefixStore.Iterator(prefixRange(append(invalidationID, types.UInt64Bytes(invalidationNonce)...)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.ContractCallTxSignature{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), &confirm) {
			break
		}
	}
}

// GetLogicConfirmsByInvalidationScopeAndNonce returns the logic call confirms
func (k Keeper) GetContractCallTxSignatureByInvalidationIDAndNonce(ctx sdk.Context, invalidationId []byte, invalidationNonce uint64) (out []types.ContractCallTxSignature) {
	k.IterateContractCallTxSignatureByInvalidationIDAndNonce(ctx, invalidationId, invalidationNonce, func(_ []byte, msg *types.ContractCallTxSignature) bool {
		out = append(out, *msg)
		return false
	})
	return
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

// GetBridgeContractAddress returns the bridge contract address on ETH
func (k Keeper) GetBridgeContractAddress(ctx sdk.Context) string {
	var a string
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractAddress, &a)
	return a
}

// GetBridgeChainID returns the chain id of the ETH chain we are running against
func (k Keeper) GetBridgeChainID(ctx sdk.Context) uint64 {
	var a uint64
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractChainID, &a)
	return a
}

// GetGravityID returns the GravityID the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a validator set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) GetGravityID(ctx sdk.Context) string {
	var a string
	k.paramSpace.Get(ctx, types.ParamsStoreKeyGravityID, &a)
	return a
}

// Set GravityID sets the GravityID the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a validator set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) SetGravityID(ctx sdk.Context, v string) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyGravityID, v)
}

// logger returns a module-specific logger.
func (k Keeper) logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UnpackEthereumEventVoteRecordEvent(att *types.EthereumEventVoteRecord) (types.EthereumEvent, error) {
	var msg types.EthereumEvent
	err := k.cdc.UnpackAny(att.Event, &msg)
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

// GetDelegateKeys iterates both the EthAddress and Orchestrator address indexes to produce
// a vector of MsgSetOrchestratorAddress entires containing all the delgate keys for state
// export / import. This may seem at first glance to be excessively complicated, why not combine
// the EthAddress and Orchestrator address indexes and simply iterate one thing? The answer is that
// even though we set the Eth and Orchestrator address in the same place we use them differently we
// always go from Orchestrator address to Validator address and from validator address to Ethereum address
// we want to keep looking up the validator address for various reasons, so a direct Orchestrator to Ethereum
// address mapping will mean having to keep two of the same data around just to provide lookups.
//
// For the time being this will serve
func (k Keeper) GetDelegateKeys(ctx sdk.Context) []*types.MsgDelegateKeys {
	store := ctx.KVStore(k.storeKey)
	prfx := []byte{types.EthereumAddressKey}

	iter := prefix.NewStore(store, prfx).Iterator(nil, nil)
	defer iter.Close()

	ethAddresses := make(map[string]string)

	for ; iter.Valid(); iter.Next() {
		// the 'key' contains both the prfx and the value, so we need
		// to cut off the starting bytes, if you don't do this a valid
		// cosmos key will be made out of EthereumAddressKey + the startin bytes
		// of the actual key
		key := iter.Key()[1:]
		value := iter.Value()
		ethAddress := string(value)
		valAddress := sdk.ValAddress(key)
		ethAddresses[valAddress.String()] = ethAddress
	}

	store = ctx.KVStore(k.storeKey)
	prfx = []byte{types.KeyOrchestratorAddress}
	iter = store.Iterator(prefixRange(prfx))
	defer iter.Close()

	orchAddresses := make(map[string]string)

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()[1:]
		value := iter.Value()
		orchAddress := sdk.AccAddress(key).String()
		valAddress := sdk.ValAddress(value)
		orchAddresses[valAddress.String()] = orchAddress
	}

	var result []*types.MsgDelegateKeys

	for valAddr, ethAddr := range ethAddresses {
		orch, ok := orchAddresses[valAddr]
		if !ok {
			// this should never happen unless the store
			// is somehow inconsistent
			panic("Can't find address")
		}
		result = append(result, &types.MsgDelegateKeys{
			OrchestratorAddress: orch,
			ValidatorAddress:    valAddr,
			EthAddress:          ethAddr,
		})

	}

	// we iterated over a map, so now we have to sort to ensure the
	// output here is deterministic, eth address chosen for no particular
	// reason
	sort.Slice(result[:], func(i, j int) bool {
		return result[i].EthAddress < result[j].EthAddress
	})

	return result
}

// GetUnbondingvalidators returns UnbondingValidators.
// Adding here in gravity keeper as cdc is available inside endblocker.
func (k Keeper) GetUnbondingvalidators(unbondingVals []byte) stakingtypes.ValAddresses {
	unbondingValidators := stakingtypes.ValAddresses{}
	k.cdc.MustUnmarshalBinaryBare(unbondingVals, &unbondingValidators)
	return unbondingValidators
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// 		Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// 				 []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
//		Example: []byte{255, 255, 255, 255} becomes nil
// MARK finish-batches: this is where some crazy shit happens
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
