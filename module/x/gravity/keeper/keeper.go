package keeper

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/query"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/ethereum/go-ethereum/common"

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
	StakingKeeper          types.StakingKeeper
	EthereumEventProcessor interface {
		Handle(sdk.Context, types.EthereumEvent) error
	}

	storeKey       sdk.StoreKey
	paramSpace     paramtypes.Subspace
	cdc            codec.BinaryCodec
	bankKeeper     types.BankKeeper
	SlashingKeeper types.SlashingKeeper
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, slashingKeeper types.SlashingKeeper) Keeper {
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
//     SignerSetTxNonce    //
/////////////////////////////

// incrementLatestSignerSetTxNonce sets the latest valset nonce
func (k Keeper) incrementLatestSignerSetTxNonce(ctx sdk.Context) uint64 {
	current := k.GetLatestSignerSetTxNonce(ctx)
	new := current + 1
	ctx.KVStore(k.storeKey).Set([]byte{types.LatestSignerSetTxNonceKey}, sdk.Uint64ToBigEndian(new))
	return new
}

// GetLatestSignerSetTxNonce returns the latest valset nonce
func (k Keeper) GetLatestSignerSetTxNonce(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LatestSignerSetTxNonceKey}); bz != nil {
		return binary.BigEndian.Uint64(bz)
	}
	return 0
}

// GetLatestSignerSetTx returns the latest validator set in state
func (k Keeper) GetLatestSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	key := types.MakeSignerSetTxKey(k.GetLatestSignerSetTxNonce(ctx))
	otx := k.GetOutgoingTx(ctx, key)
	out, _ := otx.(*types.SignerSetTx)
	return out
}

//////////////////////////////
// LastUnbondingBlockHeight //
//////////////////////////////

// setLastUnbondingBlockHeight sets the last unbonding block height
func (k Keeper) setLastUnbondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	ctx.KVStore(k.storeKey).Set([]byte{types.LastUnBondingBlockHeightKey}, sdk.Uint64ToBigEndian(unbondingBlockHeight))
}

// GetLastUnbondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnbondingBlockHeight(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastUnBondingBlockHeightKey}); len(bz) == 0 {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}

///////////////////////////////
//     ETHEREUM SIGNATURES   //
///////////////////////////////

// getEthereumSignature returns a valset confirmation by a nonce and validator address
func (k Keeper) getEthereumSignature(ctx sdk.Context, storeIndex []byte, validator sdk.ValAddress) []byte {
	return ctx.KVStore(k.storeKey).Get(types.MakeEthereumSignatureKey(storeIndex, validator))
}

// SetEthereumSignature sets a valset confirmation
func (k Keeper) SetEthereumSignature(ctx sdk.Context, sig types.EthereumTxConfirmation, val sdk.ValAddress) []byte {
	key := types.MakeEthereumSignatureKey(sig.GetStoreIndex(), val)
	ctx.KVStore(k.storeKey).Set(key, sig.GetSignature())
	return key
}

// GetEthereumSignatures returns all etherum signatures for a given outgoing tx by store index
func (k Keeper) GetEthereumSignatures(ctx sdk.Context, storeIndex []byte) map[string][]byte {
	var signatures = make(map[string][]byte)
	k.iterateEthereumSignatures(ctx, storeIndex, func(val sdk.ValAddress, h []byte) bool {
		signatures[val.String()] = h
		return false
	})
	return signatures
}

// iterateEthereumSignatures iterates through all valset confirms by nonce in ASC order
func (k Keeper) iterateEthereumSignatures(ctx sdk.Context, storeIndex []byte, cb func(sdk.ValAddress, []byte) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.EthereumSignatureKey}, storeIndex...))
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(iter.Key(), iter.Value()) {
			break
		}
	}
}

/////////////////////////
//  ORC -> VAL ADDRESS //
/////////////////////////

// SetOrchestratorValidatorAddress sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidatorAddress(ctx sdk.Context, val sdk.ValAddress, orch sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.MakeOrchestratorValidatorAddressKey(orch), val.Bytes())
}

// GetOrchestratorValidatorAddress returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidatorAddress(ctx sdk.Context, orch sdk.AccAddress) sdk.ValAddress {
	return sdk.ValAddress(ctx.KVStore(k.storeKey).Get(types.MakeOrchestratorValidatorAddressKey(orch)))
}

////////////////////////
// VAL -> ETH ADDRESS //
////////////////////////

// setValidatorEthereumAddress sets the ethereum address for a given validator
func (k Keeper) setValidatorEthereumAddress(ctx sdk.Context, validator sdk.ValAddress, ethAddr common.Address) {
	ctx.KVStore(k.storeKey).Set(types.MakeValidatorEthereumAddressKey(validator), ethAddr.Bytes())
}

// GetValidatorEthereumAddress returns the eth address for a given gravity validator
func (k Keeper) GetValidatorEthereumAddress(ctx sdk.Context, validator sdk.ValAddress) common.Address {
	return common.BytesToAddress(ctx.KVStore(k.storeKey).Get(types.MakeValidatorEthereumAddressKey(validator)))
}

////////////////////////
// ETH -> ORC ADDRESS //
////////////////////////

// setEthereumOrchestratorAddress sets the eth orch addr mapping
func (k Keeper) setEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address, orch sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.MakeEthereumOrchestratorAddressKey(ethAddr), orch.Bytes())
}

// GetEthereumOrchestratorAddress gets the orch address for a given eth address
func (k Keeper) GetEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address) sdk.AccAddress {
	return sdk.AccAddress(ctx.KVStore(k.storeKey).Get(types.MakeEthereumOrchestratorAddressKey(ethAddr)))
}

// CreateSignerSetTx gets the current signer set from the staking keeper, increments the nonce,
// creates the signer set tx object, emits an event and sets the signer set in state
func (k Keeper) CreateSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	nonce := k.incrementLatestSignerSetTxNonce(ctx)
	currSignerSet := k.CurrentSignerSet(ctx)
	newSignerSetTx := types.NewSignerSetTx(nonce, uint64(ctx.BlockHeight()), currSignerSet)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMultisigUpdateRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
			sdk.NewAttribute(types.AttributeKeySignerSetNonce, fmt.Sprint(nonce)),
		),
	)
	k.SetOutgoingTx(ctx, newSignerSetTx)
	ctx.Logger().Info("SignerSetTx created",
		"nonce", newSignerSetTx.Nonce,
		"height", newSignerSetTx.Height,
		"signers", len(newSignerSetTx.Signers),
	)
	return newSignerSetTx
}

// CurrentSignerSetTx gets powers from the store and normalizes them
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
func (k Keeper) CurrentSignerSet(ctx sdk.Context) types.EthereumSigners {
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
		if ethAddr := k.GetValidatorEthereumAddress(ctx, val); ethAddr.Hex() != "0x0000000000000000000000000000000000000000" {
			ethereumSigners[i].EthereumAddress = ethAddr.Hex()
		}
	}
	// normalize power values
	for i := range ethereumSigners {
		ethereumSigners[i].Power = sdk.NewUint(ethereumSigners[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	return (types.EthereumSigners)(ethereumSigners)
}

// GetSignerSetTxs returns all the signer set txs from the store
func (k Keeper) GetSignerSetTxs(ctx sdk.Context) (out []*types.SignerSetTx) {
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sstx, _ := otx.(*types.SignerSetTx)
		out = append(out, sstx)
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

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

// getBridgeContractAddress returns the bridge contract address on ETH
func (k Keeper) getBridgeContractAddress(ctx sdk.Context) string {
	var a string
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractAddress, &a)
	return a
}

// getBridgeChainID returns the chain id of the ETH chain we are running against
func (k Keeper) getBridgeChainID(ctx sdk.Context) uint64 {
	var a uint64
	k.paramSpace.Get(ctx, types.ParamsStoreKeyBridgeContractChainID, &a)
	return a
}

// getGravityID returns the GravityID the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a validator set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) getGravityID(ctx sdk.Context) string {
	var a string
	k.paramSpace.Get(ctx, types.ParamsStoreKeyGravityID, &a)
	return a
}

// logger returns a module-specific logger.
func (k Keeper) logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// getDelegateKeys iterates both the EthAddress and Orchestrator address indexes to produce
// a vector of MsgDelegateKeys entries containing all the delgate keys for state
// export / import. This may seem at first glance to be excessively complicated, why not combine
// the EthAddress and Orchestrator address indexes and simply iterate one thing? The answer is that
// even though we set the Eth and Orchestrator address in the same place we use them differently we
// always go from Orchestrator address to Validator address and from validator address to Ethereum address
// we want to keep looking up the validator address for various reasons, so a direct Orchestrator to Ethereum
// address mapping will mean having to keep two of the same data around just to provide lookups.
//
// For the time being this will serve
func (k Keeper) getDelegateKeys(ctx sdk.Context) (out []*types.MsgDelegateKeys) {
	store := ctx.KVStore(k.storeKey)
	iter := prefix.NewStore(store, []byte{types.ValidatorEthereumAddressKey}).Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		out = append(out, &types.MsgDelegateKeys{
			ValidatorAddress: sdk.ValAddress(iter.Key()).String(),
			EthereumAddress:  common.BytesToAddress(iter.Value()).Hex(),
		})
	}
	iter.Close()

	for _, msg := range out {
		msg.OrchestratorAddress = sdk.AccAddress(k.GetEthereumOrchestratorAddress(ctx, common.HexToAddress(msg.EthereumAddress))).String()
	}

	// we iterated over a map, so now we have to sort to ensure the
	// output here is deterministic, eth address chosen for no particular
	// reason
	sort.Slice(out[:], func(i, j int) bool {
		return out[i].EthereumAddress < out[j].EthereumAddress
	})

	return out
}

// GetUnbondingvalidators returns UnbondingValidators.
// Adding here in gravity keeper as cdc is available inside endblocker.
func (k Keeper) GetUnbondingvalidators(unbondingVals []byte) stakingtypes.ValAddresses {
	unbondingValidators := stakingtypes.ValAddresses{}
	k.cdc.MustUnmarshal(unbondingVals, &unbondingValidators)
	return unbondingValidators
}

/////////////////
// OUTGOING TX //
/////////////////

// todo: outgoingTx prefix byte
// GetOutgoingTx
func (k Keeper) GetOutgoingTx(ctx sdk.Context, storeIndex []byte) (out types.OutgoingTx) {
	if err := k.cdc.UnmarshalInterface(ctx.KVStore(k.storeKey).Get(types.MakeOutgoingTxKey(storeIndex)), &out); err != nil {
		panic(err)
	}
	return out
}

// SetOutgoingTx
func (k Keeper) SetOutgoingTx(ctx sdk.Context, outgoing types.OutgoingTx) {
	any, err := types.PackOutgoingTx(outgoing)
	if err != nil {
		panic(err)
	}
	ctx.KVStore(k.storeKey).Set(
		types.MakeOutgoingTxKey(outgoing.GetStoreIndex()),
		k.cdc.MustMarshal(any),
	)
}

// DeleteOutgoingTx deletes a given outgoingtx
func (k Keeper) DeleteOutgoingTx(ctx sdk.Context, storeIndex []byte) {
	ctx.KVStore(k.storeKey).Delete(types.MakeOutgoingTxKey(storeIndex))
}

func (k Keeper) PaginateOutgoingTxsByType(ctx sdk.Context, pageReq *query.PageRequest, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) bool) (*query.PageResponse, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.MakeOutgoingTxKey([]byte{prefixByte}))

	return query.FilteredPaginate(prefixStore, pageReq, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if !accumulate {
			return false, nil
		}

		var any cdctypes.Any
		k.cdc.MustUnmarshal(value, &any)
		var otx types.OutgoingTx
		if err := k.cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}
		if accumulate {
			return cb(key, otx), nil
		}

		return false, nil
	})
}

// IterateOutgoingTxsByType iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) IterateOutgoingTxsByType(ctx sdk.Context, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.MakeOutgoingTxKey([]byte{prefixByte}))
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var any cdctypes.Any
		k.cdc.MustUnmarshal(iter.Value(), &any)
		var otx types.OutgoingTx
		if err := k.cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}
		if cb(iter.Key(), otx) {
			break
		}
	}
}

// iterateOutgoingTxs iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) iterateOutgoingTxs(ctx sdk.Context, cb func(key []byte, outgoing types.OutgoingTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.OutgoingTxKey})
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var any cdctypes.Any
		k.cdc.MustUnmarshal(iter.Value(), &any)
		var otx types.OutgoingTx
		if err := k.cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}
		if cb(iter.Key(), otx) {
			break
		}
	}
}

// GetLastObservedSignerSetTx retrieves the last observed validator set from the store
func (k Keeper) GetLastObservedSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	key := []byte{types.LastObservedSignerSetKey}
	if val := ctx.KVStore(k.storeKey).Get(key); val != nil {
		var out types.SignerSetTx
		k.cdc.MustUnmarshal(val, &out)
		return &out
	}
	return nil
}

// setLastObservedSignerSetTx updates the last observed validator set in the stor e
func (k Keeper) setLastObservedSignerSetTx(ctx sdk.Context, signerSet types.SignerSetTx) {
	key := []byte{types.LastObservedSignerSetKey}
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshal(&signerSet))
}
