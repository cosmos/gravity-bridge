package witness

import (
	"github.com/cosmos/cosmos-sdk/examples/democoin/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	cdc *wire.Codec
	key sdk.StoreKey

	ork oracle.Keeper
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, valset sdk.ValidatorSet) Keeper {
	ork := oracle.NewKeeper(sdk.NewPrefixStoreGetter(key, []byte{0x00}), cdc, valset, sdk.NewRat(2, 3), 100)

	return Keeper{
		cdc: cdc,
		key: key,

		ork: ork,
	}
}
