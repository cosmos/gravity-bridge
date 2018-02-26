package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
)

type PeggyAccount struct {
	auth.BaseAccount
}

var _ sdk.Account = (*AppAccount)(nil)

// AccountMapper creates an account mapper given a storekey
func AccountMapper(capKey sdk.StoreKey) sdk.AccountMapper {
	var accountMapper = auth.NewAccountMapper(
		capKey,          // target store
		&PeggyAccount{}, // prototype
	)

	// Register all interfaces and concrete types that
	// implement those interfaces, here.
	cdc := accountMapper.WireCodec()
	// XXX: What does this do?
	crypto.RegisterWire(cdc)

	// Make WireCodec inaccessible before sealing
	// XXX: What does this do?
	res := accountMapper.Seal()
	return res
}
