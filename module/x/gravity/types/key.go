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
	// Key Delegation
	ValidatorEthereumAddressKey
	OrchestratorValidatorAddressKey
	EthereumOrchestratorAddressKey

	// Core types
	EthereumSignatureKey
	EthereumEventVoteRecordKey
	OutgoingTxKey
	SendToEthereumKey

	// Latest nonce indexes
	LastEventNonceByValidatorKey
	LastObservedEventNonceKey
	LatestSignerSetTxNonceKey
	LastSlashedBatchBlockKey
	LastSlashedSignerSetTxNonceKey
	LastOutgoingBatchNonceKey

	// SecondIndexSendToEthereumFeeKey indexes fee amounts by token contract address
	SecondIndexSendToEthereumFeeKey

	// BatchTxBlockKey indexes outgoing tx batches under a block height and token address
	BatchTxBlockKey

	// LastSendToEthereumIDKey indexes the lastTxPoolID
	LastSendToEthereumIDKey

	// LastEthereumBlockHeightKey indexes the latest Ethereum block height
	LastEthereumBlockHeightKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight

	LastObservedValsetKey
)

////////////////////
// Key Delegation //
////////////////////

// GetOrchestratorValidatorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOrchestratorValidatorAddressKey(orc sdk.AccAddress) []byte {
	return append([]byte{OrchestratorValidatorAddressKey}, orc.Bytes()...)
}

// GetValidatorEthereumAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetValidatorEthereumAddressKey(validator sdk.ValAddress) []byte {
	return append([]byte{ValidatorEthereumAddressKey}, validator.Bytes()...)
}

// GetEthereumOrchestratorAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthereumOrchestratorAddressKey(eth common.Address) []byte {
	return append([]byte{EthereumOrchestratorAddressKey}, eth.Bytes()...)
}

/////////////////////////
// Etheruem Signatures //
/////////////////////////

// GetEthereumSignatureKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthereumSignatureKey(storeIndex []byte, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{EthereumSignatureKey}, storeIndex, validator.Bytes()}, []byte{})
}

/////////////////////////////////
// Etheruem Event Vote Records //
/////////////////////////////////

// GetEthereumEventVoteRecordKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
func GetEthereumEventVoteRecordKey(eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{{EthereumEventVoteRecordKey}, sdk.Uint64ToBigEndian(eventNonce), claimHash}, []byte{})
}

//////////////////
// Outgoing Txs //
//////////////////

// GetOutgoingTxKey returns the store index passed with a prefix
func GetOutgoingTxKey(storeIndex []byte) []byte {
	return append([]byte{OutgoingTxKey}, storeIndex...)
}

//////////////////////
// Send To Etheruem //
//////////////////////

// GetSendToEthereumKey returns the following key format
// prefix     id
// [0x6][0 0 0 0 0 0 0 1]
func GetSendToEthereumKey(id uint64) []byte {
	return append([]byte{SendToEthereumKey}, sdk.Uint64ToBigEndian(id)...)
}

// GetFeeSecondIndexKey returns the following key format
// prefix            eth-contract-address            fee_amount
// [0x9][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000]
func GetFeeSecondIndexKey(fee sdk.Coin) []byte {
	amount := make([]byte, 32)
	return bytes.Join([][]byte{{SecondIndexSendToEthereumFeeKey}, common.HexToAddress(NewERC20TokenFromCoin(fee).Contract).Bytes(), fee.Amount.BigInt().FillBytes(amount)}, []byte{})
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
