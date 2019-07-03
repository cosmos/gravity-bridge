package querier

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/swishlabsco/peggy/x/ethbridge/types"
	keep "github.com/swishlabsco/peggy/x/oracle/keeper"
	oracletypes "github.com/swishlabsco/peggy/x/oracle/types"
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

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper keep.Keeper, codespace sdk.CodespaceType) (res []byte, errSdk sdk.Error) {
	var params types.QueryEthProphecyParams

	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	id := strconv.Itoa(params.Nonce) + params.EthereumSender
	prophecy, errSdk := keeper.GetProphecy(ctx, id)
	if errSdk != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
	}

	bridgeClaims, errSdk := MapOracleClaimsToEthBridgeClaims(params.Nonce, params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleString)
	if errSdk != nil {
		return []byte{}, errSdk
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	bz, err := codec.MarshalJSONIndent(cdc, response)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func MapOracleClaimsToEthBridgeClaims(nonce int, ethereumSender string, oracleValidatorClaims map[string]string, f func(int, gethCommon.Address, sdk.ValAddress, string) (types.EthBridgeClaim, sdk.Error)) ([]types.EthBridgeClaim, sdk.Error) {
	mappedClaims := make([]types.EthBridgeClaim, len(oracleValidatorClaims))
	i := 0
	ethereumAddress := gethCommon.HexToAddress(ethereumSender)
	for validatorBech32, validatorClaim := range oracleValidatorClaims {
		validatorAddress, parseErr := sdk.ValAddressFromBech32(validatorBech32)
		if parseErr != nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", parseErr))
		}
		mappedClaim, err := f(nonce, ethereumAddress, validatorAddress, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}
	return mappedClaims, nil
}
