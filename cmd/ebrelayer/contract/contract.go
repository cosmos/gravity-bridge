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

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// File paths to Peggy smart contract ABIs
const (
	BridgeBankABI   = "/abi/BridgeBank.abi"
	CosmosBridgeABI = "/abi/CosmosBridge.abi"
)

// LoadABI loads a smart contract as an abi.ABI
func LoadABI(contractType txs.ContractRegistry) abi.ABI {
	// Open the file containing BridgeBank contract's ABI
	var (
		_, b, _, _ = runtime.Caller(0)
		dir        = filepath.Dir(b)
	)

	var contractRaw []byte
	var err error
	switch contractType {
	case txs.CosmosBridge:
		contractRaw, err = ioutil.ReadFile(dir + CosmosBridgeABI)
		if err != nil {
			panic(err)
		}
	case txs.BridgeBank:
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
