package events

import "log"

// -----------------------------------------------------
// 	Events: Events maintains a mapping of events to an array
//		of claims made by validators.
// -----------------------------------------------------

// EventRecords : map of transaction hashes to LockEvent structs
var EventRecords = make(map[string]LockEvent)

// NewEventWrite : add a validator's address to the official claims list
func NewEventWrite(txHash string, event LockEvent) {
	EventRecords[txHash] = event
}

// IsEventRecorded : checks the sessions stored events for this transaction hash
func IsEventRecorded(txHash string) bool {
	return EventRecords[txHash].Nonce != nil
}

// PrintLockEventByTx : prints any witnessed events associated with a given transaction hash
func PrintLockEventByTx(txHash string) {
	if IsEventRecorded(txHash) {
		PrintLockEvent(EventRecords[txHash])
	} else {
		log.Printf("\nNo records from this session for tx: %v\n", txHash)
	}
}

// PrintLockEvents : prints all the claims made on this event
func PrintLockEvents() {
	// For each claim, print the validator which submitted the claim
	for txHash, event := range EventRecords {
		log.Printf("\nTransaction: %v\n", txHash)
		PrintLockEvent(event)
	}
}
