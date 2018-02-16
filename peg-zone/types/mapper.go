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
type WitnessTxMapper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
}

func NewWitnessTxMapper(key sdk.StoreKey) WitnessTxMapper {
	return WitnessTxMapper{
		key: key,
	}
}

// Implements sdk.AccountMapper.
func (wtx WitnessTxMapper) GetWitnessData(ctx sdk.Context, addr crypto.Address) sdk.Account {
	store := ctx.KVStore(am.key)
	bz := store.Get(addr)
	if bz == nil {
		return nil
	}
	acc := am.decodeAccount(bz)
	return acc
}

// Implements sdk.AccountMapper.
func (wtx WitnessTxMapper) SetAccount(ctx sdk.Context, acc sdk.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(am.key)
	bz := am.encodeAccount(acc)
	store.Set(addr, bz)
}
