package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	proto "github.com/gogo/protobuf/proto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// EthereumSignature represents one validtors signature for a given
// outgoing ethereum transaction
type EthereumSignature interface {
	proto.Message

	GetSigner() common.Address
	GetSignature() hexutil.Bytes
	GetStoreIndex() []byte
	Validate() error
}

// _ EthereumSignature = &ContractCallTxSignature{}
// _ EthereumSignature = &BatchTxSignature{}
// _ EthereumSignature = &SignerSetUpdateTxSignature{}

// EthereumEvent represents a event from the gravity contract
// on the counterparty ethereum chain
type EthereumEvent interface {
	proto.Message

	GetNonce() uint64
	GetEthereumHeight() uint64
	Hash() tmbytes.HexBytes
	Validate() error
}

// _ EthereumEvent = &DepositEvent{}
// _ EthereumEvent = &WithdrawEvent{}
// _ EthereumEvent = &ContractCallTxExecutedEvent{}
// _ EthereumEvent = &CosmosERC20DeployedEvent{}

type OutgoingTx interface {
	// NOTE: currently the function signatures here don't match, figure out how to do this properly
	// maybe add an interface arg here and typecheck in each implementation?
	GetCheckpoint([]byte) ([]byte, error)
}

// _ OutgoingTx = &EthereumSignerSet{}
// _ OutgoingTx = &BatchTx{}
// _ OutgoingTx = &ContractCallTx{}
