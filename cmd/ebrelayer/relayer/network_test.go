package relayer

// ------------------------------------------------------------
//    Network_test
//
//    Tests network.go functionality.
//
// ------------------------------------------------------------

import (
  "testing"
  "fmt"

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
  client, err := SetupWebsocketEthClient(Client)

	require.NoError(t, err)
	fmt.Printf("%+v", client)
}
