package contract

// -----------------------------------------------------
//    Contract
//
// -----------------------------------------------------

import(
	"io/ioutil"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

  // TODO: Use package instead of direct path
  // "ABI/PeggyABI"
)

func LoadAbi() abi.ABI {
  // Open the file containing Peggy contract's ABI
  rawContractAbi, errorMsg := ioutil.ReadFile("/Users/denali/go/src/github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/contract/PeggyABI.json")
  if errorMsg != nil {
    log.Fatal(errorMsg)
  }

  // Convert the raw abi into a usable format
  contractAbi, err := abi.JSON(strings.NewReader(string(rawContractAbi)))
  if err != nil {
  	log.Fatal(err)
  }

  return contractAbi
}