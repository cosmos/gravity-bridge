package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultStatus for a prophecy to start in
var DefaultStatus = "pending"

// DefaultConsensusNeeded is the default fraction of validators needed to make claims on a prophecy in order for it to pass
var DefaultConsensusNeeded = 0.7

// BridgeClaim is a struct that contains the details of a single validators claims about a single bridge transaction from ethereum to cosmos
type BridgeClaim struct {
	Nonce          int            `json:"nonce"`
	EthereumSender string         `json:"ethereum_sender"`
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Amount         sdk.Coins      `json:"amount"`
}

// NewBridgeClaim returns a new BridgeClaim with the given data contained
func NewBridgeClaim(nonce int, ethereumSender string, cosmosReceiver sdk.AccAddress, amount sdk.Coins) BridgeClaim {
	return BridgeClaim{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
	}
}

// BridgeProphecy is a struct that contains all the metadata of an oracle ritual
type BridgeProphecy struct {
	Status        string        `json:"status"`
	Nonce         int           `json:"nonce"`
	MinimumClaims int           `json:"minimum_claims"` //The minimum number of claims needed before completion logic is checked
	BridgeClaims  []BridgeClaim `json:"bridge_claims"`
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(nonce int, ethereumSender string, cosmosReceiver sdk.AccAddress, amount sdk.Coins) BridgeClaim {
	return BridgeClaim{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
	}
}
