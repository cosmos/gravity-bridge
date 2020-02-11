package contract

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// TestLoadABI test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {

	// const AbiPath = "/cmd/ebrelayer/contract/abi/BridgeBank.abi"

	//Get the ABI ready
	abi := LoadABI(txs.BridgeBank)

	require.NotNil(t, abi.Events["LogLock"])
}
