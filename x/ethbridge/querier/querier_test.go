package querier

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	oracleLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
	keeperLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
)

var (
	prophecyID0 = "0"
)

func TestNewQuerier(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper, _, _ := keeperLib.CreateTestKeepers(t, false, 0.7, nil, []int64{3, 3})

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(keeper, cdc, types.DefaultCodespace)

	//Test wrong paths
	bz, err := querier(ctx, []string{"other"}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryEthProphecy(t *testing.T) {
	cdc := codec.New()
	initialEthBridgeClaim := types.CreateTestEthClaim(t)
	initialClaim := types.CreateOracleClaimFromEthClaim(cdc, initialEthBridgeClaim)
	ctx, _, keeper, _, _ := keeperLib.CreateTestKeepers(t, false, 0.7, []oracleLib.Claim{initialClaim}, []int64{3, 7})

	testResponse := types.CreateTestQueryEthProphecyResponse(cdc, t)

	bz, err2 := cdc.MarshalJSON(types.NewQueryEthProphecyParams(initialClaim.ID))
	require.Nil(t, err2)

	query := abci.RequestQuery{
		Path: "/custom/ethbridge/prophecies",
		Data: bz,
	}

	//Test query
	res, err3 := queryEthProphecy(ctx, cdc, query, keeper, types.DefaultCodespace)
	require.Nil(t, err3)

	var ethProphecyResp types.QueryEthProphecyResponse
	err4 := cdc.UnmarshalJSON(res, &ethProphecyResp)
	require.Nil(t, err4)
	require.True(t, reflect.DeepEqual(ethProphecyResp, testResponse))

	// Test error with bad request
	query.Data = bz[:len(bz)-1]

	_, err5 := queryEthProphecy(ctx, cdc, query, keeper, types.DefaultCodespace)
	require.NotNil(t, err5)

	// Test error with nonexistent request
	query.Data = bz[:len(bz)-1]
	bz2, err6 := cdc.MarshalJSON(types.NewQueryEthProphecyParams("badProphecyID"))
	require.Nil(t, err6)

	query2 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz2,
	}

	_, err7 := queryEthProphecy(ctx, cdc, query2, keeper, types.DefaultCodespace)
	require.NotNil(t, err7)
}
