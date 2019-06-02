package contract

import (
  "testing"

  "github.com/stretchr/testify/require"
)

// Set up data for parameters and to compare against
func TestLoadABI(t *testing.T) {
  result, err := contract.LoadABI()

  require.NoError(t, err)

}