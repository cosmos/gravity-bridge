package txs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
)

const (
	TestValidator = sdk.AccAddress
	TestEventData = events.LockEvent
)

func init() {

	// Set up testing parameters for the parser
	TestValidator = sdk.AccAddress("cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq")

	// Mock expected data from the parser
	//
	// TestEvent := events.LockEvent{
	// 	Id    := "0xab85e2ceaa7d100af2f07cac01365f3777153a4e004342dca5db44e731b9d461",
	// 	From  := "0xC8Ee928625908D90d4B60859052aD200CBe2792A",
	// 	To    := "0x6e656f",
	// 	Token := "0x0000000000000000000000000000000000000000",
	// 	Value := "7",
	// 	Nonce := "39"
	// }
}

// Set up data for parameters and to compare against
func TestParsePayload(t *testing.T) {
	result, err := ParsePayloadAndRelay(TestValidator, featureVector)

	require.NoError(t, err)
	require.Equal(t, result.Nonce, "7")
	require.Equal(t, result.EthereumSender, "0xC8Ee928625908D90d4B60859052aD200CBe2792A")
	require.Equal(t, result.CosmosReceiver, "neo")
	require.Equal(t, result.Validator, TestValidator)
	require.Equal(t, result.Amount, "7")

}