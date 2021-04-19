package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName
)

var (
	// EthAddressKey indexes cosmos validator account addresses
	// i.e. cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	EthAddressKey = []byte{0x1}

	// SignersetRequestKey indexes valset requests by nonce
	SignersetRequestKey = []byte{0x2}

	// SignersetConfirmKey indexes valset confirmations by nonce and the validator account address
	// i.e cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	SignersetConfirmKey = []byte{0x3}

	// OracleEventKey Event event by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// A event is named more intuitively than an Attestation, it is literally
	// a validator making a event to have seen something happen. Events are
	// attached to attestations which can be thought of as 'the event' that
	// will eventually be executed.
	OracleEventKey = []byte{0x4}

	// AttestationKeyPrefix attestation event by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// An attestation can be thought of as the 'event to be executed' while
	// the Events are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple events vote on and
	// eventually executes
	AttestationKeyPrefix = []byte{0x5}

	// TransferTxPoolKey indexes the last nonce for the outgoing tx pool
	TransferTxPoolKey = []byte{0x6}

	// SecondIndexTransferTxFeeKey indexes fee amounts by token contract address
	SecondIndexTransferTxFeeKey = []byte{0x9}

	// BatchTxKey indexes outgoing tx batches under a nonce and token address
	BatchTxKey = []byte{0xa}

	// BatchTxBlockKey indexes outgoing tx batches under a block height and token address
	BatchTxBlockKey = []byte{0xb}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xe1}

	// SecondIndexNonceByEventKey indexes latest nonce for a given event type
	SecondIndexNonceByEventKey = []byte{0xf}

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xf1}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xf2}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0x7}

	// KeyLastTxPoolID indexes the lastTxPoolID
	KeyLastTxPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastBatchTxNonce indexes the lastBatchID
	KeyLastBatchTxNonce = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress = []byte{0xe8}

	// KeyOutgoingLogicCall indexes the outgoing logic calls
	KeyOutgoingLogicCall = []byte{0xde}

	// KeyOutgoingLogicConfirm indexes the outgoing logic confirms
	KeyOutgoingLogicConfirm = []byte{0xae}

	// LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	LastObservedEthereumBlockHeightKey = []byte{0xf9}

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key = []byte{0xf3}

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey = []byte{0xf4}

	// LastSlashedSignersetNonce indexes the latest slashed valset nonce
	LastSlashedSignersetNonce = []byte{0xf5}

	// LatestSignersetNonce indexes the latest valset nonce
	LatestSignersetNonce = []byte{0xf6}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0xf7}

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0xf8}
)

// GetOrchestratorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return append(KeyOrchestratorAddress, orc.Bytes()...)
}

// GetEthAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthAddressKey(validator sdk.ValAddress) []byte {
	return append(EthAddressKey, validator.Bytes()...)
}

// GetSignersetKey returns the following key format
// prefix    nonce
// [0x0][0 0 0 0 0 0 0 1]
func GetSignersetKey(nonce uint64) []byte {
	return append(SignersetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetSignersetConfirmKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetSignersetConfirmKey(nonce uint64, validator sdk.AccAddress) []byte {
	return append(SignersetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), validator.Bytes()...)...)
}

// GetEventKey returns the following key format
// prefix type               cosmos-validator-address                       nonce                             attestation-event-hash
// [0x0][0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// The Event hash identifies a unique event, for example it would have a event nonce, a sender and a receiver. Or an event nonce and a batch nonce. But
// the Event is stored indexed with the eventer key to make sure that it is unique.
func GetEventKey(event EthereumEvent) []byte {
	var eventHash []byte
	if event != nil {
		eventHash = event.Hash()
	} else {
		panic("No event without event!")
	}
	eventTypeLen := len([]byte(event.GetType()))
	nonceBz := sdk.Uint64ToBigEndian(event.GetNonce())
	key := make([]byte, len(OracleEventKey)+eventTypeLen+sdk.AddrLen+len(nonceBz)+len(eventHash))
	copy(key[0:], OracleEventKey)
	copy(key[len(OracleEventKey):], []byte(event.GetType()))
	// TODO this is the delegate address, should be stored by the valaddress
	copy(key[len(OracleEventKey)+eventTypeLen:], event.GetEventer())
	copy(key[len(OracleEventKey)+eventTypeLen+sdk.AddrLen:], nonceBz)
	copy(key[len(OracleEventKey)+eventTypeLen+sdk.AddrLen+len(nonceBz):], eventHash)
	return key
}

// GetTransferTxPoolKey returns the following key format
// prefix     id
// [0x6][0 0 0 0 0 0 0 1]
func GetTransferTxPoolKey(id string) []byte {
	// TODO: use hex hash .Bytes() instead ?
	return append(TransferTxPoolKey, []byte(id)...)
}

// GetBatchTxKey returns the following key format
// prefix  eth-contract-address
// [0xa][0xc783df8a850f42e7F7e57013759C285caa701eB6][HASH]
// TODO: use bytes
func GetBatchTxKey(tokenContract string, txID tmbytes.HexBytes) []byte {
	return append([]byte(tokenContract), txID.Bytes()...)
}

// GetBatchTxBlockKey returns the following key format
// prefix     blockheight
// [0xb][0 0 0 0 2 1 4 3]
func GetBatchTxBlockKey(block uint64) []byte {
	return sdk.Uint64ToBigEndian(block)
}

// GetBatchConfirmKey returns the following key format
// prefix           eth-contract-address                BatchNonce                       Validator-address
// [0xe1][0xc783df8a850f42e7F7e57013759C285caa701eB6][0 0 0 0 0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// TODO this should be a sdk.ValAddress
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, validator sdk.AccAddress) []byte {
	a := append(sdk.Uint64ToBigEndian(batchNonce), validator.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetLastEventNonceByValidatorKey indexes lateset event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

func GetDenomToERC20Key(denom string) []byte {
	return append(DenomToERC20Key, []byte(denom)...)
}

func GetERC20ToDenomKey(erc20 string) []byte {
	return append(ERC20ToDenomKey, []byte(erc20)...)
}

func GetOutgoingLogicCallKey(invalidationId []byte, invalidationNonce uint64) []byte {
	a := append(KeyOutgoingLogicCall, invalidationId...)
	return append(a, sdk.Uint64ToBigEndian(invalidationNonce)...)
}

func GetLogicCallTxKey(invalidationId []byte, invalidationNonce uint64, validator sdk.AccAddress) []byte {
	interm := append(KeyOutgoingLogicConfirm, invalidationId...)
	interm = append(interm, sdk.Uint64ToBigEndian(invalidationNonce)...)
	return append(interm, validator.Bytes()...)
}
