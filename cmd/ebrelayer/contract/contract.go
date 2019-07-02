package contract

// -------------------------------------------------------
//    Contract
//
//		Contains functionality related to the smart contract
// -------------------------------------------------------

import (
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const ABI_PATH = "cmd/ebrelayer/contract/PeggyABI.json"

func LoadABI() abi.ABI {
	// Open the file containing Peggy contract's ABI
	rawContractAbi, err := ioutil.ReadFile(ABI_PATH)
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
