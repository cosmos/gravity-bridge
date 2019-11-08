package events

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// CosmosMsg : contains data from MsgBurn and MsgLock events
type CosmosMsg struct {
	ClaimType            Event
	CosmosSender         []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}

// NewCosmosMsg : creates a new CosmosMsg
func NewCosmosMsg(
	claimType Event,
	cosmosSender []byte,
	ethereumReceiver common.Address,
	symbol string,
	amount *big.Int,
	tokenContractAddress common.Address,
) CosmosMsg {
	// Package data into a CosmosMsg
	cosmosMsg := CosmosMsg{
		ClaimType:            claimType,
		CosmosSender:         cosmosSender,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
		TokenContractAddress: tokenContractAddress,
	}

	PrintCosmosMsg(cosmosMsg)

	return cosmosMsg
}

// PrintCosmosMsg : prints a CosmosMsg struct's information
func PrintCosmosMsg(event CosmosMsg) {
	claimType := event.ClaimType.String()
	cosmosSender := string(event.CosmosSender)
	ethereumReceiver := event.EthereumReceiver.Hex()
	tokenContractAddress := event.TokenContractAddress.Hex()
	symbol := event.Symbol
	amount := event.Amount

	fmt.Printf("\nClaim Type: %v\nCosmos Sender: %v\nEthereum Recipient: %v\nToken Address: %v\nSymbol: %v\nAmount: %v\n",
		claimType, cosmosSender, ethereumReceiver, tokenContractAddress, symbol, amount)
}
