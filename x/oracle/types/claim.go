package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim is a struct that contains the details of a single validator's claim
type Claim struct {
	ID             string         `json:"id"`
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Validator      sdk.AccAddress `json:"validator"`
	Amount         sdk.Coins      `json:"amount"`
}

// NewClaim returns a new Claim with the given data contained
func NewClaim(id string, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) Claim {
	return Claim{
		ID:             id,
		CosmosReceiver: cosmosReceiver,
		Validator:      validator,
		Amount:         amount,
	}
}
