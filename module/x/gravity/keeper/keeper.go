package keeper

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

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
	StakingKeeper          types.StakingKeeper
	EthereumEventProcessor interface {
		Handle(sdk.Context, types.EthereumEvent) error
	}

	storeKey       sdk.StoreKey
	paramSpace     paramtypes.Subspace
	cdc            codec.BinaryMarshaler
	bankKeeper     types.BankKeeper
	SlashingKeeper types.SlashingKeeper
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
//     SignerSetTxNonce    //
/////////////////////////////

// SetLatestSignerSetTxNonce sets the latest valset nonce
func (k Keeper) SetLatestSignerSetTxNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestSignerSetTxNonceKey}, types.UInt64Bytes(nonce))
}

// GetLatestSignerSetTxNonce returns the latest valset nonce
func (k Keeper) GetLatestSignerSetTxNonce(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LatestSignerSetTxNonceKey}); len(bz) == 0 {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}

// GetLatestSignerSetTx returns the latest validator set in state
func (k Keeper) GetLatestSignerSetTx(ctx sdk.Context) (out *types.SignerSetTx) {
	otx := k.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(k.GetLatestSignerSetTxNonce(ctx)))
	out, _ = otx.(*types.SignerSetTx)
	return
}

////////////////////////////////////////
//     LastSlashedSignerSetTxNonce    //
////////////////////////////////////////

// SetLastSlashedSignerSetTxNonce sets the latest slashed valset nonce
func (k Keeper) SetLastSlashedSignerSetTxNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LastSlashedValsetNonceKey}, types.UInt64Bytes(nonce))
}

// GetLastSlashedValsetNonce returns the latest slashed valset nonce
func (k Keeper) GetLastSlashedValsetNonce(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastSlashedValsetNonceKey}); len(bz) == 0 {
		return 0
	} else {
		return types.UInt64FromBytes(bz)
	}
}

//////////////////////////////
// LastUnbondingBlockHeight //
//////////////////////////////

// SetLastUnBondingBlockHeight sets the last unbonding block height
func (k Keeper) SetLastUnBondingBlockHeight(ctx sdk.Context, unbondingBlockHeight uint64) {
	ctx.KVStore(k.storeKey).Set([]byte{types.LastUnBondingBlockHeight}, types.UInt64Bytes(unbondingBlockHeight))
}

// GetLastUnBondingBlockHeight returns the last unbonding block height
func (k Keeper) GetLastUnBondingBlockHeight(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastUnBondingBlockHeight}); len(bz) == 0 {
		return 0
	} else {
		return types.UInt64FromBytes(bz)
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
	key := types.GetEthereumSignatureKey(sig.GetStoreIndex(), val)
	ctx.KVStore(k.storeKey).Set(key, sig.GetSignature())
	return key
}

func (k Keeper) DeleteEthereumSignature(ctx sdk.Context, storeIndex []byte, validator sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetEthereumSignatureKey(storeIndex, validator))
}

func (k Keeper) HasEthereumSignature(ctx sdk.Context, storeIndex []byte, validator sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetEthereumSignatureKey(storeIndex, validator))
}

// GetEthereumSignatures returns all etherum signatures for a given outgoing tx by store index
func (k Keeper) GetEthereumSignatures(ctx sdk.Context, storeIndex []byte) (signatures map[string]hexutil.Bytes) {
	k.IterateEthereumSignatures(ctx, storeIndex, func(val sdk.ValAddress, h hexutil.Bytes) bool {
		signatures[val.String()] = h
		return false
	})
	return
}

// IterateEthereumSignatures iterates through all valset confirms by nonce in ASC order
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
func (k Keeper) IterateEthereumSignatures(ctx sdk.Context, storeIndex []byte, cb func(sdk.ValAddress, hexutil.Bytes) bool) {
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
	ctx.KVStore(k.storeKey).Set(types.GetOrchestratorValidatorAddressKey(orch), val.Bytes())
}

// GetOrchestratorValidatorAddress returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidatorAddress(ctx sdk.Context, orch sdk.AccAddress) sdk.ValAddress {
	return sdk.ValAddress(ctx.KVStore(k.storeKey).Get(types.GetOrchestratorValidatorAddressKey(orch)))
}

////////////////////////
// VAL -> ETH ADDRESS //
////////////////////////

// SetValidatorEthereumAddress sets the ethereum address for a given validator
func (k Keeper) SetValidatorEthereumAddress(ctx sdk.Context, validator sdk.ValAddress, ethAddr common.Address) {
	ctx.KVStore(k.storeKey).Set(types.GetValidatorEthereumAddressKey(validator), ethAddr.Bytes())
}

// GetValidatorEthereumAddress returns the eth address for a given gravity validator
func (k Keeper) GetValidatorEthereumAddress(ctx sdk.Context, validator sdk.ValAddress) common.Address {
	return common.BytesToAddress(ctx.KVStore(k.storeKey).Get(types.GetValidatorEthereumAddressKey(validator)))
}

////////////////////////
// ETH -> ORC ADDRESS //
////////////////////////

// SetEthereumOrchestratorAddress sets the eth orch addr mapping
func (k Keeper) SetEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address, orch sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetEthereumOrchestratorAddressKey(ethAddr), orch.Bytes())
}

// GetEthereumOrchestratorAddress gets the orch address for a given eth address
func (k Keeper) GetEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address) sdk.AccAddress {
	return sdk.AccAddress(ctx.KVStore(k.storeKey).Get(types.GetEthereumOrchestratorAddressKey(ethAddr)))
}

// NewSignerSetTx gets powers from the store and normalizes them
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
func (k Keeper) NewSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	ethereumSigners := make([]types.EthereumSigner, len(validators))
	var totalPower uint64
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for i, validator := range validators {
		val := validator.GetOperator()

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, val))
		totalPower += p

		ethereumSigners[i] = types.EthereumSigner{Power: p}
		if ethAddr := k.GetValidatorEthereumAddress(ctx, val); ethAddr.Hex() == "0x0000000000000000000000000000000000000000" {
			ethereumSigners[i].EthereumAddress = ethAddr.Hex()
		}
	}
	// normalize power values
	for i := range ethereumSigners {
		ethereumSigners[i].Power = sdk.NewUint(ethereumSigners[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	// TODO: make the nonce an incrementing one (i.e. fetch last nonce from state, increment, set here)
	k.GetLatestSignerSetTxNonce(ctx)
	return types.NewSignerSetTx(uint64(ctx.BlockHeight()), uint64(ctx.BlockHeight()), ethereumSigners)
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
	prfx := []byte{types.ValidatorEthereumAddressKey}

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
	prfx = []byte{types.OrchestratorValidatorAddressKey}
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
			EthereumAddress:     ethAddr,
		})

	}

	// we iterated over a map, so now we have to sort to ensure the
	// output here is deterministic, eth address chosen for no particular
	// reason
	sort.Slice(result[:], func(i, j int) bool {
		return result[i].EthereumAddress < result[j].EthereumAddress
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

/////////////////
// OUTGOING TX //
/////////////////

// todo: outgoingTx prefix byte
// GetOutgoingTx
func (k Keeper) GetOutgoingTx(ctx sdk.Context, storeIndex []byte) (out types.OutgoingTx) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetOutgoingTxKey(storeIndex))

	var any *cdctypes.Any
	k.cdc.MustUnmarshalBinaryBare(bz, any)

	out, err := types.UnpackOutgoingTx(any)
	if err != nil {
		panic(err)
	}
	return
}

// SetOutgoingTx
func (k Keeper) SetOutgoingTx(ctx sdk.Context, outgoing types.OutgoingTx) {
	any, err := types.PackOutgoingTx(outgoing)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(k.storeKey).Set(types.GetOutgoingTxKey(outgoing.GetStoreIndex()), k.cdc.MustMarshalBinaryBare(any))
}

// HasOutgoingTx
func (k Keeper) HasOutgoingTx(ctx sdk.Context, storeIndex []byte) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingTxKey(storeIndex))
}

// DeleteOutgoingTx deletes a given outgoingtx
func (k Keeper) DeleteOutgoingTx(ctx sdk.Context, storeIndex []byte) {
	ctx.KVStore(k.storeKey).Delete(types.GetOutgoingTxKey(storeIndex))
}

// IterateOutgoingTxs iterates over a specific type of outgoing transaction denoted by the chosen prefix byte
func (k Keeper) IterateOutgoingTxs(ctx sdk.Context, prefixByte byte, cb func(key []byte, outgoing types.OutgoingTx) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetOutgoingTxKey([]byte{prefixByte}))
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var any *cdctypes.Any
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), any)
		otx, _ := types.UnpackOutgoingTx(any)
		if cb(iter.Key(), otx) {
			break
		}
	}
}
