package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    crypto "github.com/tendermint/go-crypto"
)

type WithdrawTx struct {
    address crypto.Address
    amount  int64
}

var _ sdk.Msg = (*WithdrawTx)(nil)

func (wtx WithdrawTx) ValidateBasic() sdk.Error {
    return nil
}

func (wtx WithdrawTx) Type() string {
    return "WithdrawTx"
}

type WithdrawData struct {
    Amount         int64
    Destination    crypto.Address
    SignedWithdraw []SignTx
}

type SignTx struct {
}
