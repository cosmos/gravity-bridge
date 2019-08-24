package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Claim struct {
	ID               string         `json:"id"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Content          string         `json:"content"`
}

// NewClaim returns a new Claim
func NewClaim(id string, validatorAddress sdk.ValAddress, content string) Claim {
	return Claim{
		ID:               id,
		ValidatorAddress: validatorAddress,
		Content:          content,
	}
}
