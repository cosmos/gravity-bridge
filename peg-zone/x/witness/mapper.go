package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
    crypto "github.com/tendermint/go-crypto"
)

type WitnessData struct {
    Info      WitnessInfo
    Witnesses []crypto.Address
    Credited  bool
}

type WitnessMsgMapper struct {
    cdc *wire.Codec
    key sdk.StoreKey
}

func NewWitnessMsgMapper(key sdk.StoreKey) WitnessMsgMapper {
    cdc := wire.NewCodec()
    cdc.RegisterInterface((*WitnessInfo)(nil), nil)
    cdc.RegisterConcrete(LockInfo{}, "com.cosmos.peggy.LockInfo", nil)
    cdc.RegisterConcrete(WitnessData{}, "com.cosmos.peggy.WitnessData", nil)

    return WitnessMsgMapper {
        cdc: cdc,
        key: key,
    }
}

func (wmap WitnessMsgMapper) GetWitnessData(ctx sdk.Context, info WitnessInfo) WitnessData {
    store := ctx.KVStore(wmap.key)
    key, err := wmap.cdc.MarshalBinary(info)
    if err != nil {
        panic(err)
    }
    bz := store.Get(key)
    if bz == nil {
        return WitnessData {
            Info:      info,
            Witnesses: []crypto.Address{},
            Credited:  false,
        }
    }
    var data WitnessData
    if err := wmap.cdc.UnmarshalBinary(bz, &data); err != nil {
        panic(err)
    }
    return data
}

func (wmap WitnessMsgMapper) SetWitnessData(ctx sdk.Context, info WitnessInfo, data WitnessData) {
    store := ctx.KVStore(wmap.key)
    key, err := wmap.cdc.MarshalBinary(info)
    if err != nil {
        panic(err)
    }
    bz, err := wmap.cdc.MarshalBinary(data)
    if err != nil {
        panic(err)
    }
    store.Set(key, bz)
}
