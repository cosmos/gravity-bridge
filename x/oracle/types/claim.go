package types

// // Claim is a struct that contains the details of a single validator's claim
// // Note: We use strings for both fields as each field is used as a lookup key
// // in an index stored in the prophecy. ValidatorBech32 is a field that uniquely identifies this claim
// // within a specific prophecy.
// type Claim struct {
// 	ValidatorBech32 string `json:"validator_bech32"`
// 	ClaimJSON       string `json:"claim_json"`
// }

// // NewClaim returns a new Claim with the given data contained
// func NewClaim(validatorBech32 string, claimJSON string) Claim {
// 	return Claim{
// 		ValidatorBech32: validatorBech32,
// 		ClaimJSON:       claimJSON,
// 	}
// }
