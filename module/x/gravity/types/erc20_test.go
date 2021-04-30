package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateGravityDenom(t *testing.T) {
	testCases := []struct {
		name string
		denom string
		expError bool
	}{
		{"invalid denom prefix", "fake/denom", true },
		{"invalid format, correct prefix, no address", "gravity/test-denom", true},
		{"valid existing denom", "testdenom", false},
	}
	for _, tc := range testCases {
		err := ValidateGravityDenom(tc.denom)
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestValidateGravityDenomWithAddress(t *testing.T) {
	address, err := GenerateTestEthAddress()
	require.NoError(t, err, "unable to generate eth address")

	denom := fmt.Sprintf("gravity/%s", address.Hex())
	err = ValidateGravityDenom(denom)
	require.NoError(t, err, "valid external token with eth address not accepted")
}

