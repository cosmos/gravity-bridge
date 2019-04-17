package querier

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	keep "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	oracletypes "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper keep.Keeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper, codespace)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper keep.Keeper, codespace sdk.CodespaceType) (res []byte, err sdk.Error) {
	var params types.QueryEthProphecyParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	prophecy, err := keeper.GetProphecy(ctx, params.ID)
	if err != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
	}

	bridgeClaims, err2 := MapOracleClaimsToEthBridgeClaims(cdc, prophecy.Claims, ConvertOracleClaimToEthBridgeClaim)
	if err2 != nil {
		return []byte{}, err2
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, prophecy.MinimumPower, bridgeClaims)

	bz, err3 := codec.MarshalJSONIndent(cdc, response)
	if err3 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func MapOracleClaimsToEthBridgeClaims(cdc *codec.Codec, oracleClaims []oracletypes.Claim, f func(*codec.Codec, oracletypes.Claim) (types.EthBridgeClaim, sdk.Error)) ([]types.EthBridgeClaim, sdk.Error) {
	mappedClaims := make([]types.EthBridgeClaim, len(oracleClaims))
	for i, oracleClaim := range oracleClaims {
		mappedClaim, err := f(cdc, oracleClaim)
		if err != nil {
			return []types.EthBridgeClaim{}, err
		}
		mappedClaims[i] = mappedClaim
	}
	return mappedClaims, nil
}

func ConvertOracleClaimToEthBridgeClaim(cdc *codec.Codec, oracleClaim oracletypes.Claim) (types.EthBridgeClaim, sdk.Error) {
	var ethBridgeClaim types.EthBridgeClaim

	errRes := cdc.UnmarshalJSON(oracleClaim.ClaimBytes, &ethBridgeClaim)
	if errRes != nil {
		return types.EthBridgeClaim{}, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", errRes))
	}
	return ethBridgeClaim, nil
}
