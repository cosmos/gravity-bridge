package contract

// -------------------------------------------------------
//    Contract
//
//		Contains functionality related to the smart contract
// -------------------------------------------------------

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func LoadABI() abi.ABI {
	// Open the file containing Peggy contract's ABI
	rawContractAbi, errorMsg := ioutil.ReadFile("cmd/ebrelayer/contract/PeggyABI.json")
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
