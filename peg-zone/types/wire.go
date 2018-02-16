package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

var cdc = MakeTxCodec()

// RegisterWire is the functions that registers application's
// messages types to a wire.Codec.
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(WitnessTx{},
		"com.cosmos.peggy.WitnessTx", nil)
	cdc.RegisterConcrete(SendTx{},
		"com.cosmos.peggy.SendTx", nil)
	cdc.RegisterConcrete(WithdrawTx{},
		"com.cosmos.peggy.WithdrawTx", nil)
	cdc.RegisterConcrete(SignTx{},
		"com.cosmos.peggy.SignTx", nil)
}

// MakeTxCodec instantiate a wire.Codec and register
// all application's types; it returns the new codec.
func MakeTxCodec() (cdc *wire.Codec) {
	cdc = wire.NewCodec()

	// Register crypto.[PubKey,PrivKey,Signature] types.
	crypto.RegisterWire(cdc)

	// Register clearchain types.
	RegisterWire(cdc)

	// Must register message interface to parse sdk.StdTx
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	return
}
