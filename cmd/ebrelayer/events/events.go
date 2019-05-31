package events

// -----------------------------------------------------
//    Events
//
// 		Events maintains a mapping of events to an array
//		of claims made by validators.
// -----------------------------------------------------

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var OfficialClaims = make(map[string][]sdk.AccAddress)

// Add a validator's address to the official claims list
func ValidatorMakeClaim(eventHash string, validator sdk.AccAddress) int {
	OfficialClaims[eventHash] = append(OfficialClaims[eventHash], validator)

	fmt.Printf("\nValidator \"%s\" has witnessed event \"%s\".\n", validator, eventHash)

	return ClaimCount(eventHash)
}


// Returns all validators that have made claims on this event as []string
func isStoredEvent(eventHash string) bool {
	if OfficialClaims[eventHash] != nil {
		return true
	}
	return false
}

// Submitted Claims returns the list of validators which have submitted
// 	claims on an event.
func SubmittedClaims(eventHash string) []sdk.AccAddress {
		return OfficialClaims[eventHash]
}

func ClaimCount(eventHash string) int {
	return len(OfficialClaims[eventHash])
}

// Prints all the claims made on this event
func PrintClaims(event string) {
 	ethClaimsSubmitted := OfficialClaims[event]

 	// For each claim, print the validator which submitted the claim
 	fmt.Printf("Event Hash: %s", event)
  for i, claim := range ethClaimsSubmitted {
    fmt.Printf("Witness %d: %s", i, claim);
  }
}
