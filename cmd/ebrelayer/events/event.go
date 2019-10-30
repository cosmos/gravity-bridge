package events

// -----------------------------------------------------
//    Event : Creates LockEvents from new events on the ethereum
//			  Ethereum blockchain.
// -----------------------------------------------------

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// LockEvent : struct which represents a single smart contract event
type LockEvent struct {
	EthereumChainID       *big.Int
	BridgeContractAddress common.Address
	From                  common.Address
	To                    []byte
	Token                 common.Address
	Symbol                string
	Value                 *big.Int
}

// NewLockEvent : parses LogLock events using go-ethereum's accounts/abi library
func NewLockEvent(contractAbi abi.ABI, clientChainID *big.Int, contractAddress string, eventName string, eventData []byte) LockEvent {
	// Check event name
	if eventName != "LogLock" {
		log.Fatal("Only LogLock events are currently supported.")
	}

	// Parse the event's attributes as Ethereum network variables
	event := LockEvent{}

	if !common.IsHexAddress(contractAddress) {
		log.Fatalf("Only Ethereum contracts are currently supported. Invalid address: %v", contractAddress)
	}

	event.EthereumChainID = clientChainID
	event.BridgeContractAddress = common.HexToAddress(contractAddress)

	err := contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		log.Fatalf("Error unpacking: %v", err)
	}

	PrintEvent(event)

	return event
}

// PrintEvent : prints a LockEvent struct's information
func PrintEvent(event LockEvent) {
	// Convert the variables into a printable format
	chainID := event.EthereumChainID
	bridgeContractAddress := event.BridgeContractAddress.Hex()
	sender := event.From.Hex()
	token := event.Token.Hex()
	recipient := string(event.To)
	symbol := event.Symbol
	value := event.Value

	// Print the event's information
	fmt.Printf("\nChain ID: %v\nBridge contract address: %v\nToken symbol: %v\nToken contract address: %v\nSender: %v\nRecipient: %v\nValue: %v\n\n",
		chainID, bridgeContractAddress, symbol, token, sender, recipient, value)
}
