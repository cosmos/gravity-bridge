package contract

// -------------------------------------------------------
//    Contract
//
//		Contains functionality related to the smart contract
// -------------------------------------------------------

import (
	"go/build"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const ABI_PATH = "/src/github.com/cosmos/peggy/cmd/ebrelayer/contract/PeggyABI.json"

func LoadABI() abi.ABI {
	// Open the file containing Peggy contract's ABI
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	rawContractAbi, err := ioutil.ReadFile(gopath + ABI_PATH)
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
