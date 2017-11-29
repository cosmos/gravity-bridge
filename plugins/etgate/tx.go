package etgate

import (
    "math/big"

    sdk "github.com/cosmos/cosmos-sdk" // dev branch
    "github.com/cosmos/cosmos-sdk/modules/coin"
    
    //"github.com/tendermint/iavl" // dev branch

    "github.com/ethereum/go-ethereum/common"
)

const (
    ByteUpdate    = byte(0xe0)
//    ByteValChange = byte(0xe1)
    ByteDeposit   = byte(0xe2)
    ByteWithdraw  = byte(0xe3)
    ByteTransfer  = byte(0xe4)

    TypeUpdate    = NameETGate + "/update"
//    TypeValChange = NameETGate + "/valchange"
    TypeDeposit   = NameETGate + "/deposit"
    TypeWithdraw  = NameETGate + "/withdraw"
    TypeTransfer  = NameETGate + "/transfer"
)

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
    Value *big.Int
    Token common.Address
    Chain []byte
    Sequence uint64
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

type WithdrawTx struct {
    To [20]byte
    Value coin.Coins
    Token string
}

func (tx WithdrawTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx WithdrawTx) ValidateBasic() error {
    return nil
}

type TransferTx struct {
    DestChain string
    To [20]byte
    Value coin.Coins
}

func (tx TransferTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx TransferTx) ValidateBasic() error {
    if !tx.Value.IsValid() {
        return
    }
}
