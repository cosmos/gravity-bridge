package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim is a struct that contains the details of a single validator's claim
type Claim struct {
	ID         string         `json:"id"`
	Validator  sdk.AccAddress `json:"validator"`
	ClaimBytes []byte         `json:"claim_bytes"`
}

// NewClaim returns a new Claim with the given data contained
func NewClaim(id string, validator sdk.AccAddress, claimBytes []byte) Claim {
	return Claim{
		ID:         id,
		Validator:  validator,
		ClaimBytes: claimBytes,
	}
}
