package events

// -----------------------------------------------------
//    ethereumEvent : Creates LockEvents from new events on the
//			  Ethereum blockchain.
// -----------------------------------------------------

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// LockEvent : struct which represents a LogLock event
type LockEvent struct {
	Id     [32]byte
	From   common.Address
	To     []byte
	Token  common.Address
	Symbol string
	Value  *big.Int
	Nonce  *big.Int
}

// NewProphecyClaimEvent : struct which represents a LogNewProphecyClaim event
type NewProphecyClaimEvent struct {
	ProphecyID       *big.Int
	ClaimType        uint8
	CosmosSender     []byte
	EthereumReceiver common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	Symbol           string
	Amount           *big.Int
}

// UnpackLogLock : Handles new LogLock events
func UnpackLogLock(contractAbi abi.ABI, eventName string, eventData []byte) (lockEvent LockEvent) {
	event := LockEvent{}

	// Parse the event's attributes as Ethereum network variables
	err := contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		log.Fatalf("Error unpacking: %v", err)
	}

	PrintLockEvent(event)

	return event
}

// UnpackLogNewProphecyClaim : Handles new LogNewProphecyClaim events
func UnpackLogNewProphecyClaim(contractAbi abi.ABI, eventName string, eventData []byte) (newProphecyClaimEvent NewProphecyClaimEvent) {
	event := NewProphecyClaimEvent{}

	// Parse the event's attributes as Ethereum network variables
	err := contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		log.Fatalf("Error unpacking: %v", err)
	}

	PrintProphecyClaimEvent(event)

	return event
}

// PrintLockEvent : prints a LockEvent struct's information
func PrintLockEvent(event LockEvent) {
	// Convert the variables into a printable format
	value := event.Value
	symbol := event.Symbol
	token := event.Token.Hex()
	sender := event.From.Hex()
	recipient := string(event.To)
	nonce := event.Nonce

	// Print the event's information
	fmt.Printf("\nValue: %v\nToken Symbol: %v\nToken Address: %v\nSender: %v\nRecipient: %v\nNonce: %v\n\n",
		value, symbol, token, sender, recipient, nonce)
}

// PrintProphecyClaimEvent : prints a NewProphecyClaimEvent struct's information
func PrintProphecyClaimEvent(event NewProphecyClaimEvent) {
	// Convert the variables into a printable format
	id := event.ProphecyID
	claimType := event.ClaimType
	sender := string(event.CosmosSender)
	recipient := event.EthereumReceiver.Hex()
	symbol := event.Symbol
	token := event.TokenAddress.Hex()
	amount := event.Amount
	validator := event.ValidatorAddress.Hex()

	// Print the event's information
	fmt.Printf("\nProphecy ID: %v\nClaim Type: %v\nSender: %v\nRecipient %v\nSymbol %v\nToken %v\nAmount: %v\nValidator: %v\n\n",
		id, claimType, sender, recipient, symbol, token, amount, validator)
}
