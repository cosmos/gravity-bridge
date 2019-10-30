package events

// -----------------------------------------------------
//    ethereumEvent : Creates LockEvents from new events on the
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
	ProphecyID *big.Int
	ClaimType  Event
	From       []byte
	To         common.Address
	Validator  common.Address
	Token      common.Address
	Symbol     string
	Amount     *big.Int
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
	id := hex.EncodeToString(event.Id[:])
	sender := event.From.Hex()
	recipient := string(event.To)
	token := event.Token.Hex()
	symbol := event.Symbol
	value := event.Value
	nonce := event.Nonce

	// Print the event's information
	fmt.Printf("\nEvent ID: %v\nToken Symbol: %v\nToken Address: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\n\n",
		id, symbol, token, sender, recipient, value, nonce)
}

// PrintProphecyClaimEvent : prints a NewProphecyClaimEvent struct's information
func PrintProphecyClaimEvent(event NewProphecyClaimEvent) {
	// Convert the variables into a printable format
	id := event.ProphecyID
	claimType := event.ClaimType
	sender := string(event.From)
	recipient := event.To.Hex()
	symbol := event.Symbol
	token := event.Token.Hex()
	amount := event.Amount
	validator := event.Validator.Hex()

	// Print the event's information
	fmt.Printf("\nProphecy ID: %v\nClaim Type: %v\nSender: %v\nRecipient %v\nSymbol %v\nToken %v\nAmount: %v\nValidator: %v\n\n",
		id, claimType, sender, recipient, symbol, token, amount, validator)
}
