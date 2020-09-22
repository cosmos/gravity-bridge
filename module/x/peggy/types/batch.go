package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

type BatchStatus uint8

const (
	BatchStatusUnknown   BatchStatus = 0
	BatchStatusPending   BatchStatus = 1 // initial status
	BatchStatusSubmitted BatchStatus = 2 // in flight to ETH
	BatchStatusProcessed BatchStatus = 3 // observed - end state
	BatchStatusCancelled BatchStatus = 4 // end state
)

func (b BatchStatus) String() string {
	return []string{"unknown", "pending", "submitted", "observed", "processed", "cancelled"}[b]
}

type OutgoingTxBatch struct {
	Nonce              Nonce
	Elements           []OutgoingTransferTx
	CreatedAt          time.Time
	TotalFee           ERC20Token
	BridgedDenominator BridgedDenominator
	BatchStatus        BatchStatus
}

func (b *OutgoingTxBatch) Cancel() error {
	if b.BatchStatus != BatchStatusPending {
		return sdkerrors.Wrap(ErrInvalid, "status - batch not pending")
	}
	b.BatchStatus = BatchStatusCancelled
	return nil
}

func (v OutgoingTxBatch) GetCheckpoint() []byte {
	//// bytes32 encoding of "transactionBatch"
	//bytes32 methodName = 0x7472616e73616374696f6e426174636800000000000000000000000000000000;
	//
	//// Get hash of the transaction batch
	//bytes32 transactionsHash = keccak256(
	//	abi.encode(state_peggyId, methodName, _amounts, _destinations, _fees, _nonces)
	//);

	hash := crypto.Keccak256Hash([]byte(`todo`))
	return hash.Bytes()
}

func (b *OutgoingTxBatch) Observed() error {
	if b.BatchStatus != BatchStatusPending && b.BatchStatus != BatchStatusSubmitted {
		return sdkerrors.Wrap(ErrInvalid, "status")
	}
	b.BatchStatus = BatchStatusProcessed
	return nil
}

type OutgoingTransferTx struct {
	ID          uint64          `json:"txid"`
	Sender      sdk.AccAddress  `json:"sender"`
	DestAddress EthereumAddress `json:"dest_address"`
	Amount      ERC20Token      `json:"send"`
	BridgeFee   ERC20Token      `json:"bridge_fee"`
}
