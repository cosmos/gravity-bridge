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

	// SignerSetRequestKey indexes valset requests by nonce
	SignerSetRequestKey = []byte{0x2}

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

	// EventKeyPrefix
	EventKeyPrefix = []byte{0xc}

	// TransferTxKey indexes the transaction id for the outgoing transfer tx pool
	TransferTxKey = []byte{0x6}

	// TransferTxFeeKey indexes transfer txs by token contract address and fee
	TransferTxFeeKey = []byte{0x9}

	// BatchTxKey indexes outgoing tx batches under a nonce and token address
	BatchTxKey = []byte{0xa}

	// BatchTxBlockKey indexes outgoing tx batches under a block height and token address
	BatchTxBlockKey = []byte{0xb}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xe1}

	// SecondIndexNonceByEventKey indexes latest nonce for a given event type
	SecondIndexNonceByEventKey = []byte{0xf}

	// LastEventNonceByValidatorKey indexes latest event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xf1}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xf2}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0x7}

	// KeyLastTransferTxID indexes the lastTxPoolID
	KeyLastTransferTxID = append(SequenceKeyPrefix, []byte("lastTransferId")...)

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

	// LastSlashedSignerSetNonce indexes the latest slashed valset nonce
	LastSlashedSignerSetNonce = []byte{0xf5}

	// LatestSignerSetNonce indexes the latest valset nonce
	LatestSignerSetNonce = []byte{0xf6}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0xf7}

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0xf8}

	BridgeIDKey = []byte{0x19}
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
func GetSignerSetKey(nonce uint64) []byte {
	return append(SignerSetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetSignerSetConfirmKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetSignerSetConfirmKey(nonce uint64, validator sdk.ValAddress) []byte {
	return append(SignersetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), validator.Bytes()...)...)
}

// GetTransferTxKey returns the following key format
// prefix     id
// [0x6][HASH]
func GetTransferTxKey(txID tmbytes.HexBytes) []byte {
	return txID.Bytes()
}

// GetTransferTxFeeKey returns the following key format
// prefix     eth-contract-address										amount
// [0x6][0xc783df8a850f42e7F7e57013759C285caa701eB6][0 0 0 0 2 1 4 3]
func GetTransferTxFeeKey(tokenContract string, amount uint64) []byte {
	return append([]byte(tokenContract), sdk.Uint64ToBigEndian(amount)...)
}

// GetBatchTxKey returns the following key format
// prefix  eth-contract-address
// [0xa][0xc783df8a850f42e7F7e57013759C285caa701eB6][HASH]
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
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, validator sdk.ValAddress) []byte {
	a := append(sdk.Uint64ToBigEndian(batchNonce), validator.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetLastEventNonceByValidatorKey indexes lateset event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

func GetDenomToERC20Key(denom string) []byte {
	return append(DenomToERC20Key, []byte(denom)...)
}

func GetERC20ToDenomKey(erc20 string) []byte {
	return append(ERC20ToDenomKey, []byte(erc20)...)
}

func GetLogicCallTxKey(invalidationID tmbytes.HexBytes, invalidationNonce uint64) []byte {
	return append(invalidationID, sdk.Uint64ToBigEndian(invalidationNonce)...)
}
