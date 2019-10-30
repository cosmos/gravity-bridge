package events

import (
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

// CosmosMsg : contains data from MsgBurn and MsgLock events
type CosmosMsg struct {
	ClaimType            Event
	CosmosSender         []byte
	EthereumReceiver     common.Address
	Symbol               string
	Amount               *big.Int
	TokenContractAddress common.Address
}

// NewCosmosMsg : parses CosmosMsg data
func NewCosmosMsg(claimType Event, eventData [3]string) CosmosMsg {
	// Declare a new CosmosEvent
	cosmosMsg := CosmosMsg{}

	// Parse Cosmos sender
	cosmosSender := []byte(eventData[0])

	// Parse Ethereum receiver
	if !common.IsHexAddress(eventData[1]) {
		log.Fatal("Invalid recipient address:", eventData[1])
	}

	ethereumReceiver := common.HexToAddress(eventData[1])

	// Parse symbol, amount from coin
	coinRune := []rune(eventData[2])
	amount := new(big.Int)

	var symbol string

	// Set up regex
	isLetter, err := regexp.Compile(`[a-z]`)
	if err != nil {
		log.Fatal("Regex compilation error:", err)
	}

	// Iterate over each rune in the coin string
	for i, char := range coinRune {
		// Regex will match first letter [a-z] (lowercase)
		matched := isLetter.MatchString(string(char))

		// On first match, split the coin into (amount, symbol)
		if matched {
			amount, _ = amount.SetString(string(coinRune[0:i]), 10)
			symbol = string(coinRune[i:])

			break
		}
	}

	// Parse token contract address
	// TODO: Add tokenContractAddress to MsgBurn event
	tokenContractAddressString := "0xbeddb076fa4df04859098a9873591dce3e9c404d"
	if !common.IsHexAddress(tokenContractAddressString) {
		log.Fatal("Invalid token address:", tokenContractAddressString)
	}

	tokenContractAddress := common.HexToAddress(tokenContractAddressString)

	// Package the information in a MsgEvent struct
	cosmosMsg.ClaimType = claimType
	cosmosMsg.CosmosSender = cosmosSender
	cosmosMsg.EthereumReceiver = ethereumReceiver
	cosmosMsg.Symbol = symbol
	cosmosMsg.Amount = amount
	cosmosMsg.TokenContractAddress = tokenContractAddress

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
