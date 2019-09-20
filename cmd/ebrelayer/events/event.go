package events

// -----------------------------------------------------
//    Event : Creates LockEvents from new events on the ethereum
//			  Ethereum blockchain.
// -----------------------------------------------------

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// LockEvent : struct which represents a single smart contract event
type LockEvent struct {
	Chain    string
	Contract common.Address
	Id       [32]byte
	From     common.Address
	To       []byte
	Token    common.Address
	Symbol   string
	Value    *big.Int
	Nonce    *big.Int
}

// NewLockEvent : parses LogLock events using go-ethereum's accounts/abi library
func NewLockEvent(contractAbi abi.ABI, contractAddress string, eventName string, eventData []byte) LockEvent {
	// Check event name
	if eventName != "LogLock" {
		log.Fatal("Only LogLock events are currently supported.")
	}

	// Parse the event's attributes as Ethereum network variables
	event := LockEvent{}

	if !common.IsHexAddress(contractAddress) {
		log.Fatalf("Only Ethereum contracts are currently supported. Invalid address: %v", contractAddress)
	}

	event.Chain = "Ethereum"
	event.Contract = common.HexToAddress(contractAddress)

	err := contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		log.Fatalf("Unpacking: %v", err)
	}

	PrintEvent(event)

	return event
}

// PrintEvent : prints a LockEvent struct's information
func PrintEvent(event LockEvent) {
	// Convert the variables into a printable format
	chain := event.Chain
	contract := event.Contract
	id := hex.EncodeToString(event.Id[:])
	sender := event.From.Hex()
	recipient := string(event.To)
	token := event.Token.Hex()
	symbol := event.Symbol
	value := event.Value
	nonce := event.Nonce

	// Print the event's information
	fmt.Printf("\nChain: %v\nContract: %v\nEvent ID: %v\nToken Symbol: %v\nToken Address: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\n\n",
		chain, contract, id, symbol, token, sender, recipient, value, nonce)
}
