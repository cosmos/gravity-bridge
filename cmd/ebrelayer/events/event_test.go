package events

import (
  "math/big"
  "testing"
  "reflect"

  "github.com/ethereum/go-ethereum/common"
  "github.com/stretchr/testify/require"

  "github.com/swishlabsco/peggy/cmd/ebrelayer/contract"
  "github.com/swishlabsco/peggy/cmd/ebrelayer/events"
)

var TestEvent events.LockEvent

func init() {
  var arr [32]byte
  copy(arr[:], []byte("0xab85e2ceaa7d100af2f07cac01365f3777153a4e004342dca5db44e731b9d461"))
  TestEvent.Id = arr
  TestEvent.From = common.BytesToAddress([]byte("0xC8Ee928625908D90d4B60859052aD200CBe2792A"))
  TestEvent.To = []byte("0x6e656f")
  TestEvent.Token = common.BytesToAddress([]byte("0x0000000000000000000000000000000000000000"))

  value := new(big.Int)
  value, _ = value.SetString("7", 10)
  TestEvent.Value = value

  nonce := new(big.Int)
  nonce, _ = nonce.SetString("39", 10)
  TestEvent.Nonce = nonce
}

// Set up data for parameters and to compare against
func TestNewLockEvent(t *testing.T) {
  myEventData := []byte{59, 181, 100, 123, 215, 35, 69, 95, 37, 244, 14, 173, 33, 170, 131, 44, 236, 189, 193, 251, 2, 43, 207, 229, 242, 183, 141, 180, 236, 119, 201, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 200, 238, 146, 134, 37, 144, 141, 144, 212, 182, 8, 89, 5, 42, 210, 0, 203, 226, 121, 42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 149, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 45, 99, 111, 115, 109, 111, 115, 49, 112, 106, 116, 103, 117, 48, 118, 97, 117, 50, 109, 53, 50, 110, 114, 121, 107, 100, 112, 122, 116, 114, 116, 56, 56, 55, 97, 121, 107, 117, 101, 48, 104, 113, 55, 100, 102, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  contractABI := contract.LoadABI()

  lockEvent := events.NewLockEvent(contractABI, "LogLock", myEventData)
  require.True(t, reflect.DeepEqual(lockEvent, TestEvent))
}