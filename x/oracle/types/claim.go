package types

// Claim is a struct that contains the details of a single validator's claim
type Claim struct {
	ID         string `json:"id"`
	ClaimBytes []byte `json:"claim_bytes"`
}

// NewClaim returns a new Claim with the given data contained
func NewClaim(id string, claimBytes []byte) Claim {
	return Claim{
		ID:         id,
		ClaimBytes: claimBytes,
	}
}
