package events

// -----------------------------------------------------
//    Events
//
// 		Events maintains a mapping of events to an array
//		of claims made by validators.
// -----------------------------------------------------

import (
	"fmt"

)

var EventRecords = make(map[string]LockEvent)

// Add a validator's address to the official claims list
func NewEventWrite(txHash string, event LockEvent) bool {
	EventRecords[txHash] = event

	return true
}

// Checks the sessions stored events for this transaction hash
func IsEventRecorded(txHash string) bool {
	if EventRecords[txHash].Nonce == nil  {
		return false
	}
	return true
}

func PrintEventByTx(txHash string) {
	if IsEventRecorded(txHash) {
		PrintEvent(EventRecords[txHash])
	} else {
		fmt.Printf("\nNo records from this sesson for tx: %v\n", txHash)
	}
}

// Prints all the claims made on this event
func PrintEvents() error {

 	// For each claim, print the validator which submitted the claim
  for tx, event := range EventRecords {
    fmt.Printf("\nTransaction: %v\n", tx)
    PrintEvent(event)
  }

  return nil
}
