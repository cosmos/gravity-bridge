package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

type WitnessTxMapper struct {
    cdc *wire.Coded
    key sdk.StoreKey
}

func NewWitnessTxMapper(key sdk.StoreKey) WitnessTxMapper {
    cdc := wire.NewCodec()
    cdc.RegisterConcrete(WitnessData{}, "com.cosmos.peggy.WitnessData", nil)
    cdc.RegisterConcrete(WitnessTx{}, "com.cosmos.peggy.WitnessTx", nil)

    return WitnessTxMapper {
        cdc: cdc,
        key: key,
    }
}
