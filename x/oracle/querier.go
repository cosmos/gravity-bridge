package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//query endpoints supported by the oracle Querier
const (
	QueryProphecy = "prophecy"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProphecy:
			return queryProphecy(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryProphecy(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	id := path[0]
	prophecy, err := keeper.GetProphecy(ctx, id)
	if err != nil {
		return []byte{}, ErrNotFound(DefaultCodespace)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, prophecy)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
