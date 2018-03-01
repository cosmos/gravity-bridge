package witness

import (
    wire "github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
    cdc.RegisterConcrete(WitnessMsg{},
        "com.cosmos.peggy.WitnessMsg", nil)
}
