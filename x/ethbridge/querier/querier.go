package querier

import (
	"strconv"

	"github.com/cosmos/peggy/x/ethbridge/types"
	oracletypes "github.com/cosmos/peggy/x/oracle/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper types.OracleKeeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper, codespace)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper types.OracleKeeper, codespace sdk.CodespaceType) (res []byte, errSdk sdk.Error) {
	var params types.QueryEthProphecyParams

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return []byte{}, sdk.ErrInternal(sdk.AppendMsgToErr("failed to parse params: %s", err.Error()))
	}
	id := strconv.Itoa(params.EthereumChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, errSdk := keeper.GetProphecy(ctx, id)
	if errSdk != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
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
