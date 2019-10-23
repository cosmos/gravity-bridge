package events

import (
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

// MsgEvent : contains data from MsgBurn and MsgLock events
type MsgEvent struct {
	EventName            string // TODO: enum
	CosmosSender         []byte
	EthereumReceiver     common.Address
	Symbol               string
	Amount               *big.Int
	TokenContractAddress common.Address
}

// NewMsgEvent : parses MsgEvent data
func NewMsgEvent(eventName string, eventData [3]string) MsgEvent {

	// Check event name
	if eventName != "burn" && eventName != "lock" {
		log.Fatal("Only burn/lock events are supported.")
	}

	// Declare a new MsgEvent
	msgEvent := MsgEvent{}

	// Parse Cosmos sender
	cosmosSender := []byte(eventData[0])

	// Parse Ethereum receiver
	if !common.IsHexAddress(eventData[1]) {
		log.Fatal("Invalid recipient address: %v", eventData[1])
	}
	ethereumReceiver := common.HexToAddress(eventData[1])

	// Parse symbol, amount from coin
	coinRune := []rune(eventData[2])
	amount := new(big.Int)
	var symbol string
	// Iterate over each rune in the coin string
	for i, char := range coinRune {
		// Regex will match first letter [a-z] (lowercase)
		matched, err := regexp.MatchString(`[a-z]`, string(char))
		if err != nil {
			log.Fatal("Coin symbol/amount parsing error: %v", err)
		}
		// On first match, split the coin into (amount, symbol)
		if matched {
			amount, _ = amount.SetString(string(coinRune[0:i]), 10)
			symbol = string(coinRune[i:len(coinRune)])
			break
		}
	}

	// Parse token contract address
	// TODO: Add tokenContractAddress to MsgBurn event
	tokenContractAddressString := "0xbeddb076fa4df04859098a9873591dce3e9c404d"
	if !common.IsHexAddress(tokenContractAddressString) {
		log.Fatal("Invalid token address: %v", tokenContractAddressString)
	}
	tokenContractAddress := common.HexToAddress(tokenContractAddressString)

	// Package the information in a MsgEvent struct
	msgEvent.EventName = eventName
	msgEvent.CosmosSender = cosmosSender
	msgEvent.EthereumReceiver = ethereumReceiver
	msgEvent.Symbol = symbol
	msgEvent.Amount = amount
	msgEvent.TokenContractAddress = tokenContractAddress

	PrintMsgEvent(msgEvent)

	return msgEvent
}

// PrintMsgEvent : prints a MsgEvent struct's information
func PrintMsgEvent(event MsgEvent) {
	eventName := event.EventName
	cosmosSender := string(event.CosmosSender)
	ethereumReceiver := event.EthereumReceiver.Hex()
	tokenContractAddress := event.TokenContractAddress.Hex()
	symbol := event.Symbol
	amount := event.Amount

	fmt.Printf("\nMsg Type: %v\nCosmos Sender: %v\nEthereum Recipient: %v\nToken Address: %v\nSymbol: %v\nAmount: %v\n",
		eventName, cosmosSender, ethereumReceiver, tokenContractAddress, symbol, amount)
}
