package etgate

import (
    sdk "github.com/cosmos/cosmos-sdk" // dev branch
    "github.com/tendermint/iavl" // dev branch
)

const (
    ByteInit  = byte(0xe0)
    ByteUpdate    = byte(0xe1)
//    ByteValChange = byte(0xe2)
    ByteDeposit   = byte(0xe3)
//    ByteWithdraw  = byte(0xe4)
//    ByteTransfer  = byte(0xe5)

    TypeInit  = NameETGate + "/register"
    TypeUpdate    = NameETGate + "/update"
//    TypeValChange = NameETGate + "/valchange"
    TypeDeposit   = NameETGate + "/deposit"
//    TypeWithdraw  = NameETGate + "/withdraw"
//    TypeTransfer  = NameETGate + "/transfer"
)

type InitTx struct {
    Header []byte
}

func (tx InitTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx InitTx) ValidateBasic() error {
    var header eth.Header
    if err := rlp.DecodeBytes(tx.Header, &header); err != nil {
        return err
    }
    return nil
}

type UpdateTx struct {
    Headers [][]byte   
}

func (tx UpdateTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx UpdateTx) ValidateBasic() error {
    if len(tx.Headers) == 0 {
        return errors.New("Empty header list submission")
    }
    // we will pass the actual validation to validateHeaders
    // since UpdateTx will be called frequently
    return nil
}

type DepositTx struct {
    Proof LogProof
}

type Deposit struct {
    To [20]byte
    Value uint64 // extend this, 2^64 wei is not enough
    Token common.Address
    Chain []byte
    Seq uint64
}

func (tx DepositTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx DepositTx) ValidateBasic() error {
    // remove this part for efficiency?
    log, err := tx.Proof.Log()
    if err != nil {
        return err
    }

    deposit := new(Deposit)
    return depositabi.Unpack(deposit, "Deposit", log)
}
/*
// this will be replaced by IBC packetpost
type WithdrawTx struct {
    Height uint
    To [20]byte
    Value uint64
    Token string
    ChainID string
    Sequence uint64
}

func (tx WithdrawTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx WithdrawTx) ValidateBasic() error {
    return nil
}*/
