package contract

// -------------------------------------------------------
//    Contract Contains functionality for loading the
//				 smart contract
// -------------------------------------------------------

import (
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// BridgeBankABI path to file containing BridgeBank smart contract ABI
const BridgeBankABI = "/abi/BridgeBank.abi"

// CosmosBridgeABI path to file containing CosmosBridge smart contract ABI
const CosmosBridgeABI = "/abi/CosmosBridge.abi"

// LoadABI loads a smart contract as an abi.ABI
func LoadABI(cosmosSupport bool) abi.ABI {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("could not find filename")
	}
	dir := path.Dir(filename)

	var contractRaw []byte
	var err error
	switch cosmosSupport {
	case true:
		contractRaw, err = ioutil.ReadFile(dir + CosmosBridgeABI)
		if err != nil {
			panic(err)
		}
	case false:
		contractRaw, err = ioutil.ReadFile(dir + BridgeBankABI)
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
