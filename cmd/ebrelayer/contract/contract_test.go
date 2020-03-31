package contract

import (
	"testing"

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
	"github.com/stretchr/testify/require"
)

// TestLoadABI test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {
	abi := LoadABI(txs.BridgeBank)
	require.NotNil(t, abi.Events["LogLock"])
}
