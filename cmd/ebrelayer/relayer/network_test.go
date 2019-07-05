package relayer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Client = "wss://ropsten.infura.io/ws"
)

// TestIsWebsocketURL : test identification of Ethereum websocket URLs
func TestIsWebsocketURL(t *testing.T) {
	result := IsWebsocketURL(Client)
	require.True(t, result)
}

// TestSetupWebsocketEthClient : test initialization of Ethereum websocket
func TestSetupWebsocketEthClient(t *testing.T) {
	_, err := SetupWebsocketEthClient(Client)

	require.NoError(t, err)
}
