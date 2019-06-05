package relayer

// ------------------------------------------------------------
//    Relayer_Test
//
//    Tests Relayer functionality.
//
//		`go test network.go relayer.go relayer_test.go`
// ------------------------------------------------------------

import (
	"testing"
	"fmt"
	"strings"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	app "github.com/swishlabsco/cosmos-ethereum-bridge"
)

const (
	ChainID          = "testing"
	Socket           = "wss://ropsten.infura.io/ws"
	ContractAddress  =  "3de4ef81Ba6243A60B0a32d3BCeD4173b6EA02bb"
	EventSig         = "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	Validator        = "validator"
)

func TestInitRelayer(t *testing.T) {
	cdc := app.MakeCodec()

	// Parse the address of the deployed contract
	bytesContractAddress, err := hex.DecodeString(ContractAddress)
	if err != nil {
		fmt.Printf("Invalid contract-address: %v", bytesContractAddress)
	}
	contractAddress := common.BytesToAddress(bytesContractAddress)

	err = InitRelayer(cdc, ChainID, Socket, contractAddress, EventSig, Validator)

	//TODO: add validator key processing for relayer init
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Key validator not found"))
}
