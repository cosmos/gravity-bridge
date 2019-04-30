package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
    cdc.RegisterConcrete(WithdrawTx{},
        "com.cosmos.peggy.WithdrawTx", nil)
}
