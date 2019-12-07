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

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// BridgeBankABI : path to file containing BridgeBank smart contract ABI
const BridgeBankABI = "/src/github.com/cosmos/peggy/cmd/ebrelayer/contract/abi/BridgeBank.abi"

// CosmosBridgeABI : path to file containing CosmosBridge smart contract ABI
const CosmosBridgeABI = "/src/github.com/cosmos/peggy/cmd/ebrelayer/contract/abi/CosmosBridge.abi"

// LoadABI : loads a smart contract as an abi.ABI
func LoadABI(contractType txs.ContractRegistry) abi.ABI {
	// Open the file containing BridgeBank contract's ABI
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	var contractRaw []byte
	var err error

	switch contractType {
	case txs.CosmosBridge:
		contractRaw, err = ioutil.ReadFile(gopath + CosmosBridgeABI)
		if err != nil {
			panic(err)
		}
	case txs.BridgeBank:
		contractRaw, err = ioutil.ReadFile(gopath + BridgeBankABI)
		if err != nil {
			panic(err)
		}
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(string(contractRaw)))
	if err != nil {
		panic(err)
	}

	return contractABI
}
