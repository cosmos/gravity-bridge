package contract

// -------------------------------------------------------
//    Contract Contains functionality for loading the
//				 smart contract
// -------------------------------------------------------

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/trinhtan/peggy/cmd/ebrelayer/txs"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// File paths to Peggy smart contract ABIs
const (
	BridgeBankABI   = "/generated/abi/BridgeBank/BridgeBank.abi"
	CosmosBridgeABI = "/generated/abi/CosmosBridge/CosmosBridge.abi"
)

// LoadABI loads a smart contract as an abi.ABI
func LoadABI(contractType txs.ContractRegistry) abi.ABI {
	var (
		_, b, _, _ = runtime.Caller(0)
		dir        = filepath.Dir(b)
	)

	var filePath string
	switch contractType {
	case txs.CosmosBridge:
		filePath = CosmosBridgeABI
	case txs.BridgeBank:
		filePath = BridgeBankABI
	}

	// Read the file containing the contract's ABI
	contractRaw, err := ioutil.ReadFile(dir + filePath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(string(contractRaw)))
	if err != nil {
		panic(err)
	}
	return contractABI
}
