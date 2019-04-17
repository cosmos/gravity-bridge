package querier

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keep "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//query endpoints supported by the oracle Querier
const (
	QueryProphecy = "prophecies"
)

// defines the params for the following queries:
// - 'custom/oracle/prophecies/'
type QueryProphecyParams struct {
	ID string
}

func NewQueryProphecyParams(id string) QueryProphecyParams {
	return QueryProphecyParams{
		ID: id,
	}
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper keep.Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProphecy:
			return queryProphecy(ctx, cdc, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryProphecyParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	prophecy, err := keeper.GetProphecy(ctx, params.ID)
	if err != nil {
		return []byte{}, types.ErrProphecyNotFound(keeper.Codespace())
	}

	bz, err2 := codec.MarshalJSONIndent(cdc, prophecy)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
