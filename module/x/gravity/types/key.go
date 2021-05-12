package types

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

const (
	_ = byte(iota)

	// EthereumAddressKey indexes cosmos validator account addresses
	// i.e. cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	EthereumAddressKey

	// SignerSetTxKey indexes valset requests by nonce
	SignerSetTxKey

	// SignerSetTxSignatureKey indexes valset confirmations by nonce and the validator account address
	// i.e cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	SignerSetTxSignatureKey

	// EthereumEventVoteRecordKey attestation details by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	EthereumEventVoteRecordKey

	// OutgoingTXPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTXPoolKey

	// SecondIndexOutgoingTXFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTXFeeKey

	// BatchTxKey indexes outgoing tx batches under a nonce and token address
	BatchTxKey

	// BatchTxBlockKey indexes outgoing tx batches under a block height and token address
	BatchTxBlockKey

	// BatchTxSignatureKey indexes validator confirmations by token contract address
	BatchTxSignatureKey

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey

	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID

	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress

	// ContractCallTxKey indexes the outgoing logic calls
	ContractCallTxKey

	// ContractCallTxSignatureKey indexes the outgoing logic confirms
	ContractCallTxSignatureKey

	// LastEthereumBlockHeightKey indexes the latest Ethereum block height
	LastEthereumBlockHeightKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey

	// LastSlashedValsetNonce indexes the latest slashed valset nonce
	LastSlashedValsetNonce

	// LatestSignerSetTxNonce indexes the latest valset nonce
	LatestSignerSetTxNonce

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight
)

// GetOrchestratorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return append([]byte{KeyOrchestratorAddress}, orc.Bytes()...)
}

// GetEthereumAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthereumAddressKey(validator sdk.ValAddress) []byte {
	return append([]byte{EthereumAddressKey}, validator.Bytes()...)
}

// GetSignerSetTxKey returns the following key format
// prefix    nonce
// [0x0][0 0 0 0 0 0 0 1]
func GetSignerSetTxKey(nonce uint64) []byte {
	return append([]byte{SignerSetTxKey}, sdk.Uint64ToBigEndian(nonce)...)
}

// GetSignerSetTxSignatureKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetSignerSetTxSignatureKey(nonce uint64, validator sdk.ValAddress) []byte {
	return append([]byte{SignerSetTxSignatureKey}, append(sdk.Uint64ToBigEndian(nonce), validator.Bytes()...)...)
}

// GetEthereumEventVoteRecordKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetEthereumEventVoteRecordKey(eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{{EthereumEventVoteRecordKey}, sdk.Uint64ToBigEndian(eventNonce), claimHash}, []byte{})
}

// GetEthereumEventVoteRecordKeyWithHash returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetEthereumEventVoteRecordKeyWithHash(eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{{EthereumEventVoteRecordKey}, sdk.Uint64ToBigEndian(eventNonce), claimHash} , []byte{})
}

// GetOutgoingTxPoolKey returns the following key format
// prefix     id
// [0x6][0 0 0 0 0 0 0 1]
func GetOutgoingTxPoolKey(id uint64) []byte {
	return append([]byte{OutgoingTXPoolKey}, sdk.Uint64ToBigEndian(id)...)
}

// GetBatchTxKey returns the following key format
// prefix     nonce                     eth-contract-address
// [0xa][0 0 0 0 0 0 0 1][0xc783df8a850f42e7F7e57013759C285caa701eB6]
func GetBatchTxKey(tokenContract string, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxKey}, sdk.Uint64ToBigEndian(nonce), common.HexToAddress(tokenContract).Bytes()}, []byte{})
}

// GetBatchTxBlockKey returns the following key format
// prefix     blockheight
// [0xb][0 0 0 0 2 1 4 3]
func GetBatchTxBlockKey(block uint64) []byte {
	return append([]byte{BatchTxBlockKey}, sdk.Uint64ToBigEndian(block)...)
}

// GetBatchTxSignatureKey returns the following key format
// prefix           eth-contract-address                BatchNonce                       Validator-address
// [0xe1][0xc783df8a850f42e7F7e57013759C285caa701eB6][0 0 0 0 0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// TODO this should be a sdk.ValAddress
func GetBatchTxSignatureKey(tokenContract string, batchNonce uint64, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{BatchTxSignatureKey}, common.HexToAddress(tokenContract).Bytes(), sdk.Uint64ToBigEndian(batchNonce), validator.Bytes()}, []byte{})
}

// GetFeeSecondIndexKey returns the following key format
// prefix            eth-contract-address            fee_amount
// [0x9][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000]
func GetFeeSecondIndexKey(fee sdk.Coin) []byte {
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	token := NewERC20TokenFromCoin(fee)
	return bytes.Join([][]byte{{SecondIndexOutgoingTXFeeKey}, common.HexToAddress(token.Contract).Bytes(), amount}, []byte{})
}

// GetLastEventNonceByValidatorKey indexes lateset event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append([]byte{LastEventNonceByValidatorKey}, validator.Bytes()...)
}

func GetDenomToERC20Key(denom string) []byte {
	return append([]byte{DenomToERC20Key}, []byte(denom)...)
}

func GetERC20ToDenomKey(erc20 string) []byte {
	return append([]byte{ERC20ToDenomKey}, []byte(erc20)...)
}

func GetContractCallTxKey(invalidationId []byte, invalidationNonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxKey}, invalidationId, sdk.Uint64ToBigEndian(invalidationNonce)}, []byte{})
}

// prefix    invalidationID  nonce  validatorAddr
func GetContractCallTxSignatureKey(invalidationId []byte, invalidationNonce uint64, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{ContractCallTxSignatureKey}, invalidationId, sdk.Uint64ToBigEndian(invalidationNonce), validator.Bytes()}, []byte{})
}
