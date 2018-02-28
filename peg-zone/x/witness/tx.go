package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    crypto "github.com/tendermint/go-crypto"

    wire "github.com/tendermint/go-wire"
)

type WitnessMsg struct {
    Info   WitnessInfo
    Signer crypto.Address
}

var _ sdk.Msg = (*WitnessMsg)(nil)

func (msg WitnessMsg) ValidateBasic() sdk.Error {
    return msg.Info.ValidateBasic()
    if len(msg.Destination) != 20 ||
       len(msg.Token)       != 20 ||
       len(msg.Signer)      != 20 {
        return ErrInvalidWitnessMsg()    
    }
}

func (msg WitnessMsg) Type() string {
    return "WitnessMsg"
}

func (msg WitnessMsg) Get(key interface{}) interface{} {
    return nil
}

func newCodec() *wire.Codec {
    cdc := wire.NewCodec()
    cdc.RegisterConcrete(WitnessMsg{}, "com.cosmos.peggy.WitnessMsg", nil)
    cdc.RegisterInterface((*WitnessInfo)(nil), nil)
    cdc.RegisterConcrete(LockInfo{}, "com.cosmos.peggy.LockInfo", nil)
    return cdc
}

func (msg WitnessMsg) GetSignBytes() []byte {
    cdc := newCodec()
    bz, err := cdc.MarshalBinary(msg)
    if err != nil {
        panic(err)
    }
    return bz
}

func (msg WitnessMsg) GetSigners() []crypto.Address {
    return []crypto.Address{ msg.Signer }
}

type WitnessInfo interface {
    isWitnessInfo()
}

type LockInfo struct {
    Destination crypto.Address
    Amount      int64
    Token       crypto.Address
}

var _ WitnessInfo = (LockInfo)(nil)

func (info LockInfo) isWitnessInfo() {}

