package etend

import (
    "math/big"

    sdk "github.com/cosmos/cosmos-sdk" // dev branch
    "github.com/tendermint/iavl" // dev branch
)

const (
    ByteDeposit = byte(0xe8)
    ByteWithdraw = byte(0xe9)
    ByteTransfer = byte(0xea)

    TypeDeposit = NameETEnd + "/deposit"
    TypeWithdraw = NameETEnd + "/withdraw"
    TypeTransfer = NameETEnd + "/transfer"
)

type DepositTx struct {
    To [20]byte
    Value []byte //big.Int.Bytes    
}

func (tx DepositTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx DepositTx) ValidateBasic() error {
    return nil
}

type WithdrawTx struct {
    To [20]byte
    Value int64
    Token string
    OriginChain string
    Sequence uint64
}

func (tx)
