package querier

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	keep "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

var (
	prophecyID0 = "0"
)

func TestNewQuerier(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper := keep.CreateTestKeepers(t, false, 1000)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(keeper, cdc)

	//Test wrong paths
	bz, err := querier(ctx, []string{"other"}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryDelegation(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper := keep.CreateTestKeepers(t, false, 10000)

	//Initial setup
	testProphecy := keep.CreateTestProphecy(t)
	err := keeper.CreateProphecy(ctx, testProphecy)
	require.NoError(t, err)

	bz, err2 := cdc.MarshalJSON(NewQueryProphecyParams(keep.TestID))
	require.Nil(t, err2)

	query := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz,
	}

	//Test query
	res, err3 := queryProphecy(ctx, cdc, query, keeper)
	require.Nil(t, err3)

	var prophecyResp types.BridgeProphecy
	err4 := cdc.UnmarshalJSON(res, &prophecyResp)
	require.Nil(t, err4)

	require.True(t, reflect.DeepEqual(prophecyResp, testProphecy))

	// Test error with bad request
	query.Data = bz[:len(bz)-1]

	_, err = queryProphecy(ctx, cdc, query, keeper)
	require.NotNil(t, err)

	// Test error with nonexistent request
	query.Data = bz[:len(bz)-1]
	bz2, err5 := cdc.MarshalJSON(NewQueryProphecyParams("wrongEthereumAddress0"))
	require.Nil(t, err5)

	query2 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz2,
	}

	_, err = queryProphecy(ctx, cdc, query2, keeper)
	require.NotNil(t, err)
}
