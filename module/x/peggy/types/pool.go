package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OutgoingTx struct {
	Sender      sdk.AccAddress `json:"sender"`
	DestAddress string         `json:"dest_address"`
	Amount      TransferCoin   `json:"send"`
	BridgeFee   TransferCoin   `json:"bridge_fee"`
	//BridgeContractAddress string         `json:"bridge_contract_address"` // todo: do we need this?
}

type BridgedDenominator struct {
	//ChainID         string
	BridgeContractAddress string
	TokenID               string
}

// TransferCoin is an outgoing token
type TransferCoin struct {
	Denom  string
	Amount uint64
}

func NewTransferCoin(denom string, amount uint64) TransferCoin {
	return TransferCoin{Denom: denom, Amount: amount}
}

func AsTransferCoin(denominator BridgedDenominator, voucher sdk.Coin) TransferCoin {
	assertPeggyVoucher(voucher)
	return NewTransferCoin(denominator.TokenID, voucher.Amount.Uint64())
}

const (
	VoucherDenomPrefix = "peggy"
	DenomSeparator     = "" // todo: only a-z0-9 supported
	//DenomSeparator     = "/"  // todo: not allowed in this versions sdk coin demon
	VoucherDenomLen = 15 // todo: cut to 15 to match this versions sdk coin demon
	//hashLen            = 64
	//separatorLen       = len(DenomSeparator)
	//prefixLen          = len(VoucherDenomPrefix)
	//VoucherDenomLen    = hashLen + prefixLen + separatorLen
)

func assertPeggyVoucher(s sdk.Coin) {
	if len(s.Denom) != VoucherDenomLen || !strings.HasPrefix(s.Denom, VoucherDenomPrefix) {
		panic("invalid denom type")
	}
	if s.Amount.IsNegative() || !s.Amount.IsUint64() {
		panic("invalid amount type")
	}
}

// IDSet is a collection of DB keys
type IDSet []uint64
