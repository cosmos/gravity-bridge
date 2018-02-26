package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    crypto "github.com/tendermint/go-crypto"
)
/*
type WitnessMsg struct {
    amount      int64
    destination crypto.Address
    token       crypto.Address
    signer      crypto.Address
}

var _ sdk.Msg = (*WitnessMsg)(nil)

func (msg WitnessMsg) ValidateBasic() sdk.Error {
    return nil
}

func (msg WitnessMsg) Type() string {
    return "WitnessTx"
}

func (msg WitnessMsg) Get(key interface{}) interface{} {
    return nil
}

func (msg WitnessMsg) GetSignBytes() []byte {
    b, err := proto.Marshal(msg)
}
*/

// Using LockMsg directly because of GetSignBytes

type WitnessMsg interface {
    isWitnessMsg()
}

var _ WitnessMsg = (*LockMsg)(nil)

func (msg LockMsg) isWitnessMsg() {}

var _ sdk.Msg = (*LockMsg)(nil)

func (msg LockMsg) ValidateBasic() sdk.Error {
    return nil
}

func (msg LockMsg) Type() string {
    return "LockMsg"
}

func (msg LockMsg) Get(key interface{}) interface{} {
    return nil
}

func (msg LockMsg) GetSignBytes() []byte {
    data, err := proto.Marshal(msg)
    if err != nil {
        panic(err)
    }
    return data
}

func (msg LockMsg) GetSigners() []crypto.Address {
    return []crypto.Address{ msg.Signer }
}

type WitnessData struct {
    Witnesses      []crypto.Address
    Amount         int64
    Destination    crypto.Address
    credited       bool
}

