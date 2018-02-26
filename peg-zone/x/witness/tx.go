package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    crypto "github.com/tendermint/go-crypto"

    wire "github.com/tendermint/go-wire"
)

type LockMsg struct {
    Destination crypto.Address
    Amount      int64
    Token       []byte
    Signer      crypto.Address
}

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

func newCodec() *wire.Codec {
    cdc := wire.NewCodec()
    cdc.RegisterConcrete(LockMsg{}, "com.cosmos.peggy.LockMsg", nil)
    return cdc
}

func (msg LockMsg) GetSignBytes() []byte {
    cdc := newCodec()
    bz, err := cdc.MarshalBinary(msg)
    if err != nil {
        panic(err)
    }
    return bz
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

