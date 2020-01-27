package querier

import (
	"fmt"
	"strconv"

	"github.com/cosmos/peggy/x/ethbridge/types"
	oracletypes "github.com/cosmos/peggy/x/oracle/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper types.OracleKeeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper types.OracleKeeper) (res []byte, errSdk error) {
	var params types.QueryEthProphecyParams

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse params: %s", err.Error()))
	}
	id := strconv.Itoa(params.EthereumChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, errSdk := keeper.GetProphecy(ctx, id)
	if errSdk != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound
	}

	bridgeClaims, errSdk := types.MapOracleClaimsToEthBridgeClaims(params.EthereumChainID, params.BridgeContractAddress, params.Nonce, params.Symbol, params.TokenContractAddress, params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleString)
	if errSdk != nil {
		return []byte{}, errSdk
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	bz, err := cdc.MarshalJSONIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}

	return bz, nil
}
