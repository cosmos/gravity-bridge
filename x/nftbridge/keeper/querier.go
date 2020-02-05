package keeper

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
	oracletypes "github.com/cosmos/peggy/x/oracle/types"
)

// TODO: move to x/oracle

// NewQuerier is the module level router for state queries
func NewQuerier(keeper ethbridge.OracleKeeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryNFTProphecy:
			return queryNFTProphecy(ctx, cdc, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nftbridge query endpoint")
		}
	}
}

func queryNFTProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper ethbridge.OracleKeeper) ([]byte, error) {
	var params types.QueryNFTProphecyParams

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse params: %s", err.Error()))
	}

	id := strconv.Itoa(params.EthereumChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, found := keeper.GetProphecy(ctx, id)
	if !found {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, id)
	}

	bridgeClaims, err := types.MapOracleClaimsToNFTBridgeClaims(params.EthereumChainID, params.BridgeContractAddress, params.Nonce, params.Symbol, params.TokenContractAddress, params.EthereumSender, prophecy.ValidatorClaims, types.CreateNFTClaimFromOracleString)
	if err != nil {
		return nil, err
	}

	response := types.NewQueryNFTProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	return cdc.MarshalJSONIndent(response, "", "  ")
}
