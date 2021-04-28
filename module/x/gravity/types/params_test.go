package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamsValidation(t *testing.T) {
	params := DefaultParams()
	err := params.ValidateBasic()
	require.NoError(t, err, "default parameter validation failed")

	params = DefaultParams()
	params.TargetBatchTimeout = 50000
	err = params.ValidateBasic()
	require.Error(t, err, "invalid params accepted")
}

func toInterface(i interface{}) interface{} {
	return i
}

func TestValidateBoundedUInt64(t *testing.T) {
	err := validateBoundedUInt64(toInterface(uint64(50)), 0, 100)
	require.NoError(t, err, "failed valid uint64 validation")

	err = validateBoundedUInt64(toInterface("notUInt64"), 0, 100)
	require.Error(t, err, "accepted incorrect type")

	err = validateBoundedUInt64(toInterface(uint64(10)), 20, 100)
	require.Error(t, err, "accepted too small value")

	err = validateBoundedUInt64(toInterface(uint64(110)), 20, 100)
	require.Error(t, err, "accepted too large value")
}

func TestValidateBoundedDev(t *testing.T) {
	err := validateBoundedDec(toInterface(sdk.MustNewDecFromStr("0.6")), sdk.ZeroDec(), sdk.OneDec())
	require.NoError(t, err, "failed valid dec validation")

	err = validateBoundedDec(toInterface("notDec"), sdk.ZeroDec(), sdk.NewDec(100))
	require.Error(t, err, "accepted incorrect type")

	err = validateBoundedDec(toInterface(sdk.NewDec(10)), sdk.NewDec(20), sdk.NewDec(100))
	require.Error(t, err, "accepted too small value")

	err = validateBoundedDec(toInterface(sdk.NewDec(110)), sdk.NewDec(20), sdk.NewDec(100))
	require.Error(t, err, "accepted too large value")
}