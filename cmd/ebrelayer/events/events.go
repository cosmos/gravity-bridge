package events

// -----------------------------------------------------
//    Events
//
// 		Events maintains a mapping of events to an array
//		of claims made by validators.
// -----------------------------------------------------

// Declare the official map:(EVENT_HASH => []ETH_CLAIM_ID)
const officialClaims := make(map[string][]string);

// Adds a new event to the official mapping, allowing claims to be made upon it by validators
// TODO: Replace the eventHash with the event's unique _id
func AddEvent(eventHash) {
	officialClaims[eventHash] = []string
}

// Add a validator's address to the official claims list
func ValidatorMakeClaim(eventHash string, validator sdk.AccAddress) {
	officialClaims[eventHash] = append(officialClaims[eventHash, validator])
	fmt.Printf("Validator %s has witnessed event %s", validator, eventHash)
}

// Prints all the claims made on this event
func PrintClaims(eventHash string) {
 	ethClaimsSubmitted := officialClaims(eventHash)

 	// For each claim, print the validator which submitted the claim
 	fmt.Printf("Event Hash: %s", eventHash)
  for i, claim := range ethClaimsSubmitted {
    fmt.Printf("Witness %d: %s", i, claim.Validator);
  }
}

// Returns all validators that have made claims on this event as []string
func Claims(eventHash string) []string {
 	return officialClaims(eventHash)
}
