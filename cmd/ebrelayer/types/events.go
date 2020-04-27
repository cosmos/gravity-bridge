package types

import "log"

// TODO: This should be moved to new 'events' directory and expanded so that it can
// serve as a local store of witnessed events and allow for re-trying failed relays.

// EventRecords map of transaction hashes to EthereumEvent structs
var EventRecords = make(map[string]EthereumEvent)

// NewEventWrite add a validator's address to the official claims list
func NewEventWrite(txHash string, event EthereumEvent) {
	EventRecords[txHash] = event
}

// IsEventRecorded checks the sessions stored events for this transaction hash
func IsEventRecorded(txHash string) bool {
	return EventRecords[txHash].Nonce != nil
}

// PrintEventByTx prints any witnessed events associated with a given transaction hash
func PrintEventByTx(txHash string) {
	if IsEventRecorded(txHash) {
		log.Println(EventRecords[txHash].String())
	} else {
		log.Printf("\nNo records from this session for tx: %v\n", txHash)
	}
}

// PrintEvents prints all the claims made on this event
func PrintEvents() {
	// For each claim, print the validator which submitted the claim
	for txHash, event := range EventRecords {
		log.Printf("\nTransaction: %v\n", txHash)
		log.Println(event.String())
	}
}
