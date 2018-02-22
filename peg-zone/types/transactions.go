package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

const (
	WitnessTx  = "WitnessTx"
	SendTx     = "sendtx"
	WithdrawTx = "withdrawtx"
	SignTx     = "signtx"
)

// ------------------------------
// WitnessTx

type WitnessTx struct {
	Amount      int64
	Destination crypto.Address
    Token       crypto.Address
}

var _ sdk.Msg = (*WitnessTx)(nil)

func (wtx WitnessTx) ValidateBasic() sdk.Error {
	return nil
}

func (wtx WitnessTx) Type() string {
	return WitnessTx
}

type WitnessData struct {
    Witnesses   []crypto.Address
    Amount      int64
	Destination crypto.Address
	Credited    bool
}

// ------------------------------
// SendTx

type SendTx struct {
	from   crypto.Address
	to     crypto.Address
	amount int64
}

var _ sdk.Msg = (*SendTx)(nil)

func (sdx SendTx) ValidateBasic() sdk.Error {
	return nil
}

func (sdx SendTx) Type() string {
	return SendTx
}

// ------------------------------
// WithdrawTx

type WithdrawTx struct {
	address crypto.Address
	amount  int64
}

type WithdrawData struct {
	Amount         int64
	Destination    crypto.Address
	SignedWithdraw []SignTx
}

// --------------------------------
// SignTx

type SignTx struct {
}
