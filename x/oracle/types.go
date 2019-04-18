package oracle

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var PendingStatus = "pending"
var CompleteStatus = "complete"

// DefaultConsensusNeeded is the default fraction of validators needed to make claims on a prophecy in order for it to pass
var DefaultConsensusNeeded = 0.7

// BridgeClaim is a struct that contains the details of a single validators claims about a single bridge transaction from ethereum to cosmos
type BridgeClaim struct {
	ID             string         `json:"id"`
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Validator      sdk.AccAddress `json:"validator"`
	Amount         sdk.Coins      `json:"amount"`
}

// NewBridgeClaim returns a new BridgeClaim with the given data contained
func NewBridgeClaim(id string, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) BridgeClaim {
	return BridgeClaim{
		ID:             id,
		CosmosReceiver: cosmosReceiver,
		Validator:      validator,
		Amount:         amount,
	}
}

// BridgeProphecy is a struct that contains all the metadata of an oracle ritual
type BridgeProphecy struct {
	ID            string        `json:"id"`
	Status        string        `json:"status"`
	MinimumClaims int           `json:"minimum_claims"` //The minimum number of claims needed before completion logic is checked
	BridgeClaims  []BridgeClaim `json:"bridge_claims"`
}

func (prophecy BridgeProphecy) String() string {
	prophecyJSON, err := json.Marshal(prophecy)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}

// NewBridgeProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewBridgeProphecy(id string, status string, minimumClaims int, bridgeClaims []BridgeClaim) BridgeProphecy {
	return BridgeProphecy{
		ID:            id,
		Status:        status,
		MinimumClaims: minimumClaims,
		BridgeClaims:  bridgeClaims,
	}
}

// NewEmptyBridgeProphecy returns a blank prophecy, used with errors
func NewEmptyBridgeProphecy() BridgeProphecy {
	return NewBridgeProphecy("", "", 0, nil)
}
