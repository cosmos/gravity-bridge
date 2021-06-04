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
	LastSlashedOutgoingTxBlockKey
	LastSlashedSignerSetTxNonceKey
	LastOutgoingBatchNonceKey

	// LastSendToEthereumIDKey indexes the lastTxPoolID
	LastSendToEthereumIDKey

	// LastEthereumBlockHeightKey indexes the latest Ethereum block height
	LastEthereumBlockHeightKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey

	// LastUnBondingBlockHeightKey indexes the last validator unbonding block height
	LastUnBondingBlockHeightKey

	LastObservedSignerSetKey
)

////////////////////
// Key Delegation //
////////////////////

// MakeOrchestratorValidatorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeOrchestratorValidatorAddressKey(orc sdk.AccAddress) []byte {
	return append([]byte{OrchestratorValidatorAddressKey}, orc.Bytes()...)
}

// MakeValidatorEthereumAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeValidatorEthereumAddressKey(validator sdk.ValAddress) []byte {
	return append([]byte{ValidatorEthereumAddressKey}, validator.Bytes()...)
}

// MakeEthereumOrchestratorAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEthereumOrchestratorAddressKey(eth common.Address) []byte {
	return append([]byte{EthereumOrchestratorAddressKey}, eth.Bytes()...)
}

/////////////////////////
// Etheruem Signatures //
/////////////////////////

// MakeEthereumSignatureKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEthereumSignatureKey(storeIndex []byte, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{EthereumSignatureKey}, storeIndex, validator.Bytes()}, []byte{})
}

/////////////////////////////////
// Etheruem Event Vote Records //
/////////////////////////////////

// MakeEthereumEventVoteRecordKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
func MakeEthereumEventVoteRecordKey(eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{{EthereumEventVoteRecordKey}, sdk.Uint64ToBigEndian(eventNonce), claimHash}, []byte{})
}

//////////////////
// Outgoing Txs //
//////////////////

// MakeOutgoingTxKey returns the store index passed with a prefix
func MakeOutgoingTxKey(storeIndex []byte) []byte {
	return append([]byte{OutgoingTxKey}, storeIndex...)
}

//////////////////////
// Send To Etheruem //
//////////////////////

// MakeSendToEthereumKey returns the following key format
// prefix            eth-contract-address            fee_amount        id
// [0x9][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000][0 0 0 0 0 0 0 1]
func MakeSendToEthereumKey(id uint64, fee ERC20Token) []byte {
	amount := make([]byte, 32)
	return bytes.Join([][]byte{{SendToEthereumKey}, common.HexToAddress(fee.Contract).Bytes(), fee.Amount.BigInt().FillBytes(amount), sdk.Uint64ToBigEndian(id)}, []byte{})
}

// MakeLastEventNonceByValidatorKey indexes lateset event nonce by validator
// MakeLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append([]byte{LastEventNonceByValidatorKey}, validator.Bytes()...)
}

func MakeDenomToERC20Key(denom string) []byte {
	return append([]byte{DenomToERC20Key}, []byte(denom)...)
}

func MakeERC20ToDenomKey(erc20 string) []byte {
	return append([]byte{ERC20ToDenomKey}, []byte(erc20)...)
}

func MakeSignerSetTxKey(nonce uint64) []byte {
	return append([]byte{SignerSetTxPrefixByte}, sdk.Uint64ToBigEndian(nonce)...)
}

func MakeBatchTxKey(addr common.Address, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxPrefixByte}, addr.Bytes(), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeContractCallTxKey(invalscope []byte, invalnonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxPrefixByte}, invalscope, sdk.Uint64ToBigEndian(invalnonce)}, []byte{})
}
