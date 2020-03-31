package types

import "log"

// TODO: This should be moved to new 'events' directory and expanded so that it can
// serve as a local store of witnessed events and allow for re-trying failed relays.

// EventRecords map of transaction hashes to LockEvent structs
var EventRecords = make(map[string]LockEvent)

// NewEventWrite add a validator's address to the official claims list
func NewEventWrite(txHash string, event LockEvent) {
	EventRecords[txHash] = event
}

// IsEventRecorded checks the sessions stored events for this transaction hash
func IsEventRecorded(txHash string) bool {
	return EventRecords[txHash].Nonce != nil
}

// PrintLockEventByTx prints any witnessed events associated with a given transaction hash
func PrintLockEventByTx(txHash string) {
	if IsEventRecorded(txHash) {
		log.Println(EventRecords[txHash].String())
	} else {
		log.Printf("\nNo records from this session for tx: %v\n", txHash)
	}
}

// PrintLockEvents prints all the claims made on this event
func PrintLockEvents() {
	// For each claim, print the validator which submitted the claim
	for txHash, event := range EventRecords {
		log.Printf("\nTransaction: %v\n", txHash)
		log.Println(event.String())
	}
}
