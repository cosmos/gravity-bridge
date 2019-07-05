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
const AbiPath = "cmd/ebrelayer/contract/PeggyABI.json"

// LoadABI : loads a smart contract as an abi.ABI from a .json file
func LoadABI() abi.ABI {
	// Open the file containing Peggy contract's ABI
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	rawContractAbi, err := ioutil.ReadFile(gopath + AbiPath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractAbi, err := abi.JSON(strings.NewReader(string(rawContractAbi)))
	if err != nil {
		panic(err)
	}

	return contractAbi
}
