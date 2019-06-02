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
func IsStoredEvent(eventHash string) bool {
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
func PrintClaims(event string) error {
	fmt.Println("\nEvent id: ", event)

 	ethClaims := OfficialClaims[event]

 	// Check to see if there are any claims to report
 	if ethClaims == nil {
 		fmt.Println("\nThis event has 0 claims")
 		return fmt.Errorf("\nThis event has 0 claims")
 	} else {
 		fmt.Println("\nClaim count: ", ClaimCount)

 	}

 	// For each claim, print the validator which submitted the claim
  for i, claim := range ethClaims {
    fmt.Printf("Witness %d: %s", i, claim)
  }

  return nil
}
