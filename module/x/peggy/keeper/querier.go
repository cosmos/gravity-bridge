package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	// QueryResolve = "resolve"
	// QueryWhois   = "whois"
	// QueryNames   = "names"
	QueryCurrentValset = "currentValset"
	QueryValsetRequest = "valsetRequest"
	QueryValsetConfirm = "valsetConfirm"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		// case QueryResolve:
		// 	return queryResolve(ctx, path[1:], req, keeper)
		// case QueryWhois:
		// 	return queryWhois(ctx, path[1:], req, keeper)
		// case QueryNames:
		// 	return queryNames(ctx, req, keeper)
		case QueryCurrentValset:
			return queryCurrentValset(ctx, keeper)
		case QueryValsetRequest:
			return queryValsetRequest(ctx, path[1:], keeper)
		case QueryValsetConfirm:
			return queryValsetConfirm(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}

func queryCurrentValset(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentValset(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValsetRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, err
	}

	valset := keeper.GetValsetRequest(ctx, nonce)
	res, err := codec.MarshalJSONIndent(keeper.cdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValsetConfirm(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	accAddress, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	valset := keeper.GetValsetConfirm(ctx, nonce, accAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// // nolint: unparam
// func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
// 	value := keeper.ResolveName(ctx, path[0])

// 	if value == "" {
// 		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not resolve name")
// 	}

// 	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResResolve{Value: value})
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
// 	}

// 	return res, nil
// }

// // nolint: unparam
// func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
// 	whois := keeper.GetWhois(ctx, path[0])

// 	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
// 	}

// 	return res, nil
// }

// func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
// 	var namesList types.QueryResNames

// 	iterator := keeper.GetNamesIterator(ctx)

// 	for ; iterator.Valid(); iterator.Next() {
// 		namesList = append(namesList, string(iterator.Key()))
// 	}

// 	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
// 	}

// 	return res, nil
// }
