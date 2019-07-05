package relayer

// ------------------------------------------------------------
//    Network_test
//
//    Tests network.go functionality.
//
// ------------------------------------------------------------

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Client = "wss://ropsten.infura.io/ws"
)

func TestIsWebsocketURL(t *testing.T) {
	result := IsWebsocketURL(Client)
	require.True(t, result)
}

func TestSetupWebsocketEthClient(t *testing.T) {
	_, err := SetupWebsocketEthClient(Client)

	require.NoError(t, err)
}
