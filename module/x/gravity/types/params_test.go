package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamsValidation(t *testing.T) {
	params := DefaultParams()
	err := params.ValidateBasic()
	require.NoError(t, err, "default parameter validation failed")

	params.BridgeChainId = -1
	err = params.ValidateBasic()
	require.Error(t, err, "negative chain ID accepted")

	params = DefaultParams()
	params.TargetBatchTimeout = 120000
	require.Error(t, err, "too long target batch timeout accepted")
}