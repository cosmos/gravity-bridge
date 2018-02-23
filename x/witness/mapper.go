package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

type WitnessMsgMapper struct {
    cdc *wire.Codec
    key sdk.StoreKey
}

func NewWitnessMsgMapper(key sdk.StoreKey) WitnessMsgMapper {
    cdc := wire.NewCodec()
    cdc.RegisterInterface((*WitnessMsg)(nil), nil)
    cdc.RegisterConcrete(LockMsg{}, "com.cosmos.peggy.LockMsg", nil)
    cdc.RegisterConcrete(WitnessData{}, "com.cosmos.peggy.WitnessData", nil)

    return WitnessMsgMapper {
        cdc: cdc,
        key: key,
    }
}

func (wmap WitnessMsgMapper) GetWitnessData(ctx sdk.Context, tx WitnessMsg) *WitnessData {
    store := ctx.KVStore(wmap.key)
    key, err := wmap.cdc.MarshalBinary(tx)
    if err != nil {
        panic(err)
    }
    bz := store.Get(key)
    if bz == nil {
        return nil
    }
    var data WitnessData
    if err := wmap.cdc.UnmarshalBinary(bz, &data); err != nil {
        panic(err)
    }
    return &data
}

func (wmap WitnessMsgMapper) SetWitnessData(ctx sdk.Context, tx WitnessMsg, data WitnessData) {
    store := ctx.KVStore(wmap.key)
    key, err := wmap.cdc.MarshalBinary(tx)
    if err != nil {
        panic(err)
    }
    bz, err := wmap.cdc.MarshalBinary(data)
    if err != nil {
        panic(err)
    }
    store.Set(key, bz)
}
