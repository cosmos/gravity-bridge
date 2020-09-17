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

type OutgoingTxBatch struct {
	Elements              []OutgoingTransferTx
	CreatedAt             time.Time
	TotalFee              TransferCoin
	CosmosDenom           VoucherDenom
	BridgedTokenID        string
	BridgeContractAddress string
	BatchStatus           BatchStatus
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

type OutgoingTransferTx struct {
	ID          uint64         `json:"txid"`
	Sender      sdk.AccAddress `json:"sender"`
	DestAddress string         `json:"dest_address"`
	Amount      TransferCoin   `json:"send"`
	BridgeFee   TransferCoin   `json:"bridge_fee"`
}

// TransferCoin is an outgoing token
type TransferCoin struct {
	TokenID string
	Amount  uint64
}

func (t TransferCoin) Add(o TransferCoin) TransferCoin {
	if t.TokenID != o.TokenID {
		panic("invalid token")
	}
	sum := sdk.NewInt(int64(t.Amount)).AddRaw(int64(o.Amount))
	if !sum.IsUint64() {
		panic("invalid amount")
	}
	return NewTransferCoin(t.TokenID, sum.Uint64())
}

func NewTransferCoin(tokenID string, amount uint64) TransferCoin {
	return TransferCoin{TokenID: tokenID, Amount: amount}
}

func AsTransferCoin(denominator BridgedDenominator, voucher sdk.Coin) TransferCoin {
	assertPeggyVoucher(voucher)
	return NewTransferCoin(denominator.TokenID, voucher.Amount.Uint64())
}
