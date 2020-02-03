package relayer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ValidTestnetWebSocket = "wss://ropsten.infura.io/ws"
	ValidLocalWebSocket   = "ws://127.0.0.1:7545/"
	InvalidWebSocket      = "http://localhost:7545"
)

// TestIsWebsocketURL test identification of Ethereum websocket URLs
func TestIsWebsocketURL(t *testing.T) {
	validTestnetRes := IsWebsocketURL(ValidTestnetWebSocket)
	require.True(t, validTestnetRes)

	validLocalRes := IsWebsocketURL(ValidLocalWebSocket)
	require.True(t, validLocalRes)

	invalidRes := IsWebsocketURL(InvalidWebSocket)
	require.False(t, invalidRes)
}

// TestSetupWebsocketEthClient test initialization of Ethereum websocket
func TestSetupWebsocketEthClient(t *testing.T) {
	_, err := SetupWebsocketEthClient(InvalidWebSocket)
	require.Error(t, err, "invalid websocket eth client URL: "+InvalidWebSocket)
}
