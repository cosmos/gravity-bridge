package contract

// -------------------------------------------------------
//    Contract : Contains functionality for loading the
//				 smart contract
// -------------------------------------------------------

import (
	"go/build"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// AbiPath : path to the file containing the smart contract's ABI
const AbiPath = "/src/github.com/cosmos/peggy/cmd/ebrelayer/contract/abi/BridgeBank.abi"

// LoadABI : loads a smart contract as an abi.ABI
func LoadABI() abi.ABI {
	// Open the file containing BridgeBank contract's ABI
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	peggyABI, err := ioutil.ReadFile(gopath + AbiPath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractAbi, err := abi.JSON(strings.NewReader(string(peggyABI)))
	if err != nil {
		panic(err)
	}

	return contractAbi
}
