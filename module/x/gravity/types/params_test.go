package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_ValidateBasic(t *testing.T) {
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

func TestParams_validateBoundedUInt64(t *testing.T) {
	testCases := []struct {
		name string
		i interface{}
		lower uint64
		upper uint64
		expError bool
	}{
		{"valid input", toInterface(uint64(50)), 0, 100, false},
		{"incorrect type", toInterface("notUInt64"), 0, 100, true},
		{"too small value", toInterface(uint64(10)), 20, 100, true},
		{"too large value", toInterface(uint64(110)), 20, 100, true},
	}
	for _, tc := range testCases {
		err := validateBoundedUInt64(tc.i, tc.lower, tc.upper)
		if err != nil {
			if tc.expError {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err, tc.name)
			}
		}
	}
}

func TestParams_validateBoundedDec(t *testing.T) {
	testCases := []struct {
		name string
		i interface{}
		lower sdk.Dec
		upper sdk.Dec
		expError bool
	}{
		{"valid input", toInterface(sdk.MustNewDecFromStr("0.6")), sdk.ZeroDec(), sdk.OneDec(), false},
		{"incorrect type", toInterface("notDec"), sdk.ZeroDec(), sdk.NewDec(100), true},
		{"value too small", toInterface(sdk.NewDec(10)), sdk.NewDec(20), sdk.NewDec(100), true},
		{"value too large", toInterface(sdk.NewDec(110)), sdk.NewDec(20), sdk.NewDec(100), true},
	}
	for _, tc := range testCases {
		err := validateBoundedDec(tc.i, tc.lower, tc.upper)
		if err != nil {
			if tc.expError {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err, tc.name)
			}
		}
	}
}