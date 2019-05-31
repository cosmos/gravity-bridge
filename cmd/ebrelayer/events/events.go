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

var OfficialClaims map[string][]sdk.AccAddress

// Declare the official map:(EVENT_HASH => []ETH_CLAIM_ID)
func main() {
	officialClaims := map[string][]sdk.AccAddress{
	    "first": {},
	    // "second": []string{"one", "two", "three", "four", "five"},
	    // "third": []string{"quarter", "half"},
	}

	OfficialClaims = officialClaims
}

// Adds a new event to the official mapping, allowing claims to be made upon it by validators
// TODO: Replace the eventHash with the event's unique _id
func AddEvent(eventHash string) {
	OfficialClaims[eventHash] = []sdk.AccAddress{}
}

// Add a validator's address to the official claims list
func ValidatorMakeClaim(eventHash string, validator sdk.AccAddress) {
	OfficialClaims[eventHash] = append(OfficialClaims[eventHash], validator)
	fmt.Printf("Validator %s has witnessed event %s", validator, eventHash)
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

// // Returns all validators that have made claims on this event as []string
// func Claims(eventHash string) []string {
//  	return OfficialClaims(eventHash)
// }

// ClaimingValidators returns a list of validators which have submitted
// claims on the event.
func Claims(event string) []sdk.AccAddress {
	return OfficialClaims[event]
}
