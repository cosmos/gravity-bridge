package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

type WithdrawTxMapper struct {
    cdc *wire.Coded
    key sdk.StoreKey
}

func NewWithdrawTxMapper(key sdk.StoreKey) WitnessTxMapper {
    cdc := wire.NewCodec()
    cdc.RegisterConcrete(WithdrawData{}, "com.cosmos.peggy.WithdrawData", nil)
    cdc.RegisterConcrete(WithdrawTx{}, "com.cosmos.peggy.WithdrawTx", nil)

    return WithdrawTxMapper {
        cdc: cdc,
        key: key,
    }
}
