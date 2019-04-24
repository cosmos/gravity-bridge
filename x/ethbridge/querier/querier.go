package querier

import (
	"fmt"
	"strconv"

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

	id := strconv.Itoa(params.Nonce) + params.EthereumSender
	prophecy, err := keeper.GetProphecy(ctx, id)
	if err != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
	}

	bridgeClaims, err2 := MapOracleClaimsToEthBridgeClaims(params.Nonce, params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleClaim)
	if err2 != nil {
		return []byte{}, err2
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	bz, err3 := codec.MarshalJSONIndent(cdc, response)
	if err3 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func MapOracleClaimsToEthBridgeClaims(nonce int, ethereumSender string, oracleValidatorClaims map[string]string, f func(int, string, sdk.ValAddress, string) (types.EthBridgeClaim, sdk.Error)) ([]types.EthBridgeClaim, sdk.Error) {
	mappedClaims := make([]types.EthBridgeClaim, len(oracleValidatorClaims))
	i := 0
	for validatorBech32, validatorClaim := range oracleValidatorClaims {
		validatorAddress, parseErr := sdk.ValAddressFromBech32(validatorBech32)
		if parseErr != nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", parseErr))
		}
		mappedClaim, err := f(nonce, ethereumSender, validatorAddress, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}
	return mappedClaims, nil
}
