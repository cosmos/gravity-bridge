package txs

import (
	"testing"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	"github.com/ethereum/go-ethereum/common"
)

var TestValidator sdk.AccAddress
var TestEventData events.LockEvent

func init() {

	// Set up testing parameters for the parser
	testValidator, err := sdk.AccAddressFromBech32("cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq")
  if err != nil {
    fmt.Errorf("%s", err)
  }
  TestValidator = testValidator

	// Mock expected data from the parser
	TestEventData := events.LockEvent{}

	var arr [32]byte
	copy(arr[:], []byte("0xab85e2ceaa7d100af2f07cac01365f3777153a4e004342dca5db44e731b9d461"))
	TestEventData.Id = arr
	TestEventData.From = common.BytesToAddress([]byte("0xC8Ee928625908D90d4B60859052aD200CBe2792A"))
	TestEventData.To = []byte("0x6e656f")
	TestEventData.Token = common.BytesToAddress([]byte("0x0000000000000000000000000000000000000000"))

	value := new(big.Int)
	value, okValue := value.SetString("7", 10)
	if !okValue {
	  fmt.Println("SetString: error")
  }
	TestEventData.Value = value

	nonce := new(big.Int)
	nonce, okNonce := nonce.SetString("39", 10)
	if !okNonce {
	  fmt.Println("SetString: error")
  }
	TestEventData.Nonce = nonce

	fmt.Printf("%+v", TestEventData)
	
}

// Set up data for parameters and to compare against
func TestParsePayload(t *testing.T) {
	result, err := ParsePayload(TestValidator, &TestEventData)

	require.NoError(t, err)
	fmt.Printf("%+v", result)

	// TODO: check each individual argument
	// require.Equal(t, "7", string(result.Nonce))
	// require.Equal(t, common.BytesToAddress([]byte("0xC8Ee928625908D90d4B60859052aD200CBe2792A")), result.EthereumSender)
	// require.Equal(t, result.CosmosReceiver, "neo")
	// require.Equal(t, result.Validator, TestValidator)
	// require.Equal(t, result.Amount, 7)

}