package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoadABI test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {

	const AbiPath = "/cmd/ebrelayer/contract/abi/BridgeBank.abi"

	//Get the ABI ready
	abi := LoadABI(true)

	require.NotNil(t, abi.Events["LogLock"])
}
