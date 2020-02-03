package utils

// --------------------------------------------------------
//      Utils
//
//      Utils contains utility functionality for the ebrelayer.
// --------------------------------------------------------

import (
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

// IsZeroAddress checks an Ethereum address and returns a bool which indicates if it is the null address
func IsZeroAddress(address common.Address) bool {
	return address == common.HexToAddress(nullAddress)
}

// GetSymbolAmountFromCoin Parse (symbol, amount) from coin string
func GetSymbolAmountFromCoin(coin string) (string, *big.Int) {
	coinRune := []rune(coin)
	amount := new(big.Int)

	var symbol string

	// Set up regex
	isLetter := regexp.MustCompile(`[a-z]`)

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

	return symbol, amount
}
