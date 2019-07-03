package txs

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
)

var TestValidator sdk.ValAddress
var TestEventData events.LockEvent

func init() {

	// Set up testing parameters for the parser
	testValidator, err := sdk.ValAddressFromBech32("cosmosvaloper18hf69vxn8a3tkladruxgxgv8tl8sl54gpuezxk")
	if err != nil {
		panic(err)
	}
	TestValidator = testValidator

	// Mock expected data from the parser
	TestEventData := events.LockEvent{}

	var arr [32]byte
	copy(arr[:], []byte("0xab85e2ceaa7d100af2f07cac01365f3777153a4e004342dca5db44e731b9d461"))
	TestEventData.Id = arr
	TestEventData.From = common.HexToAddress("0xC8Ee928625908D90d4B60859052aD200CBe2792A")
	TestEventData.To = []byte("0x6e656f")
	TestEventData.Token = common.HexToAddress("0xC8Ee928625908D90d4B60859052aD200CBe2792A")

	value := new(big.Int)
	value, _ = value.SetString("7", 10)
	TestEventData.Value = value

	nonce := new(big.Int)
	nonce, _ = nonce.SetString("39", 10)
	TestEventData.Nonce = nonce
}

// Set up data for parameters and to compare against
func TestParsePayload(t *testing.T) {
	result, err := ParsePayload(TestValidator, &TestEventData)
	require.NoError(t, err)

	fmt.Printf("%+v", result)

	// TODO: check each individual argument
	// require.Equal(t, "7", string(result.Nonce))
	// require.Equal(t, common.HexToAddress("0xC8Ee928625908D90d4B60859052aD200CBe2792A"), result.EthereumSender)
	// require.Equal(t, result.CosmosReceiver, "neo")
	// require.Equal(t, result.Validator, TestValidator)
	// require.Equal(t, result.Amount, 7)

}
