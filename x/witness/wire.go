package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
    cdc.RegisterConcrete(WitnessTx{},
        "com.cosmos.peggy.WitnessTx", nil)
}
