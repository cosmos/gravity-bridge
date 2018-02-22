package types

import (
	"fmt"
	"reflect"

	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"

	sdk "github.com/cosmos/cosmos-sdk/types"   
)

// Implements sdk.AccountMapper.
// This AccountMapper encodes/decodes accounts using the
// go-wire (binary) encoding/decoding library.
// WitnessTxMapper : WitnessTx -> WitnessData
type WitnessTxMapper struct {
	// The (unexposed) key used to access the store from the Context.
    cdc *wire.Codec
	key sdk.StoreKey
}

func NewWitnessTxMapper(key sdk.StoreKey) WitnessTxMapper {
    cdc := wire.NewCodec()
    cdc.RegisterConcrete(WitnessData{}, "com.tendermint/WitnessData", nil)
    cdc.RegisterConcrete(WitnessTx{}, "com.tendermint/WitnessTx", nil)

    return WitnessTxMapper{
        cdc: cdc,
		key: key,
	}
}

// Implements sdk.AccountMapper.
func (wtx WitnessTxMapper) GetWitnessData(ctx sdk.Context, tx WitnessTx) WitnessData {
	store := ctx.KVStore(wtx.key)
    key, err := wtx.cdc.MarshalBinary(tx)
    if err != nil {
        panic(err)
    }
    bz := store.Get(key)
	if data == nil {
		return nil
	}
    var data WitnessData
    if err := wtx.cdc.UnmarshalBinary(bz, &data); err != nil {
        panic(err)
    }
	return data
}

func (wtx WitnessTxMapper) SetWitnessData(ctx sdk.Context, tx WitnessTx, data WitnessData) {
    store := ctx.KVStore(wtx.key)
    key, err := wtx.cdc.MarshalBinary(tx)
    if err != nil {
        panic(err)
    }
    bz := wtx.cdc.MarshalBinary(data)
    store.Set(key, bz)
}

// Implements sdk.AccountMapper.
func (wtx WitnessTxMapper) SetAccount(ctx sdk.Context, acc sdk.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(wtx.key)
	bz := am.encodeAccount(acc)
	store.Set(addr, bz)
}

// ValidatorMapper : crypto.Address -> bool

type ValidatorMapper struct {
    key sdk.StoreKey
    am sdk.AccountMapper
}

func NewValidatorMapper(key sdk.StoreKey) ValidatorMapper {
    return ValidatorMapper{
		key: key,
	}
}

func (val ValidatorMapper) GetValidators(ctx sdk.Context) []crypto.Address {
    res := []crypto.Address{}
    store := ctx.KVStore(val.key)
    for iter := store.Iterator([]byte{}, []byte(nil)); iter.Valid(); iter.Next() {
        res = append(res, iter.Key())
    }
    iter.Close()
    return res
}

func (val ValidatorMapper) IsValidator(ctx sdk.Context, addr crypto.Address) bool {
    store := ctx.KVStore(val.key)
    return store.Get(addr)
}
