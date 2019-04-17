package keeper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

func TestCreateGetProphecy(t *testing.T) {
	ctx, _, keeper := CreateTestKeepers(t, false, 0)
	testProphecy := CreateTestProphecy(t)

	//Test normal Creation
	err := keeper.CreateProphecy(ctx, testProphecy)
	require.NoError(t, err)

	//Test bad Creation
	badProphecy := CreateTestProphecy(t)
	badProphecy.MinimumPower = -1
	err = keeper.CreateProphecy(ctx, badProphecy)

	badProphecy2 := CreateTestProphecy(t)
	badProphecy2.ID = ""
	err = keeper.CreateProphecy(ctx, badProphecy2)
	require.Error(t, err)

	badProphecy3 := CreateTestProphecy(t)
	badProphecy3.BridgeClaims = []types.BridgeClaim{}
	err = keeper.CreateProphecy(ctx, badProphecy3)
	require.Error(t, err)

	//Test retrieval
	prophecy, err := keeper.GetProphecy(ctx, testProphecy.ID)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(testProphecy, prophecy))
}
