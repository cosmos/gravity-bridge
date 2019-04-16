package keeper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGetProphecy(t *testing.T) {
	ctx, _, keeper := CreateTestKeepers(t, false, 0)
	testProphecy := CreateTestProphecy(t)
	err := keeper.CreateProphecy(ctx, testProphecy)
	require.NoError(t, err)

	prophecy, err := keeper.GetProphecy(ctx, testProphecy.ID)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(testProphecy, prophecy))
}

//TODO: Test sad paths have errors
