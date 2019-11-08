package txs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      Cosmos Bridge.
// --------------------------------------------------------

import (
	"errors"
	"log"
	"math/big"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	tmCommon "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	ethbridgeTypes "github.com/cosmos/peggy/x/ethbridge/types"
)

// LogLockToEthBridgeClaim : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func LogLockToEthBridgeClaim(valAddr sdk.ValAddress, event *events.LockEvent) (ethbridgeTypes.EthBridgeClaim, error) {
	witnessClaim := ethbridgeTypes.EthBridgeClaim{}

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Int64())

	// Sender type casting (address.common -> string)
	sender := ethbridgeTypes.NewEthereumAddress(event.From.Hex())

	// Recipient type casting ([]bytes -> sdk.AccAddress)
	recipient, err := sdk.AccAddressFromBech32(string(event.To))
	if err != nil {
		return witnessClaim, err
	}
	if recipient.Empty() {
		return witnessClaim, errors.New("empty recipient address")
	}

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && event.Token != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return witnessClaim, errors.New("symbol \"eth\" must have null address set as token address")
	}

	// Amount type casting (*big.Int -> sdk.Coins)
	coins := sdk.Coins{sdk.NewInt64Coin(symbol, event.Value.Int64())}

	// Package the information in a unique EthBridgeClaim
	witnessClaim.Nonce = nonce
	witnessClaim.EthereumSender = sender
	witnessClaim.ValidatorAddress = valAddr
	witnessClaim.CosmosReceiver = recipient
	witnessClaim.Amount = coins

	return witnessClaim, nil
}

// BurnLockEventToCosmosMsg : parses data from a Burn/Lock event witnessed on Cosmos into a CosmosMsg struct
func BurnLockEventToCosmosMsg(claimType events.Event, attributes []tmCommon.KVPair) events.CosmosMsg {
	// Set up variables
	var cosmosSender []byte
	var ethereumReceiver, tokenContractAddress common.Address
	var symbol string
	var amount *big.Int

	// Iterate over attributes
	for _, attribute := range attributes {
		// Get (key, value) for each attribute
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		// Set variable based on value of CosmosMsgAttributeKey
		switch key {
		case events.CosmosSender.String():
			// Parse sender's Cosmos address
			cosmosSender = []byte(val)
		case events.EthereumReceiver.String():
			// Confirm recipient is valid Ethereum address
			if !common.IsHexAddress(val) {
				log.Fatal("Invalid recipient address:", val)
			}
			// Parse recipient's Ethereum address
			ethereumReceiver = common.HexToAddress(val)
		case events.Coin.String():
			// Parse symbol and amount from coin string
			symbol, amount = getSymbolAmountFromCoin(val)
		case events.TokenContractAddress.String():
			// Confirm token contract address is valid Ethereum address
			if !common.IsHexAddress(val) {
				log.Fatal("Invalid token address:", val)
			}
			// Parse token contract address
			tokenContractAddress = common.HexToAddress(val)
		}
	}

	// Package the event data into a CosmosMsg
	return events.NewCosmosMsg(claimType, cosmosSender, ethereumReceiver, symbol, amount, tokenContractAddress)
}

// ProphecyClaimToSignedOracleClaim : packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToSignedOracleClaim(event events.NewProphecyClaimEvent) OracleClaim {
	// Parse relevant data into type byte[]
	prophecyID := event.ProphecyID.Bytes()
	sender := event.CosmosSender
	recipient := []byte(event.EthereumReceiver.Hex())
	token := []byte(event.TokenAddress.Hex())
	amount := event.Amount.Bytes()
	validator := []byte(event.ValidatorAddress.Hex())

	// Generate rawHash using ProphecyClaim data
	hash := GenerateClaimHash(prophecyID, sender, recipient, token, amount, validator)

	// Sign the hash using the active validator's private key
	signature := SignClaim(hash)

	// Package the ProphecyID, Message, and Signature into an OracleClaim
	oracleClaim := OracleClaim{
		ProphecyID: event.ProphecyID,
		Message:    hash,
		Signature:  signature,
	}

	return oracleClaim
}

// CosmosMsgToProphecyClaim : parses event data from a CosmosMsg, packaging it as a ProphecyClaim
func CosmosMsgToProphecyClaim(event events.CosmosMsg) ProphecyClaim {
	claimType := event.ClaimType
	cosmosSender := event.CosmosSender
	ethereumReceiver := event.EthereumReceiver
	tokenContractAddress := event.TokenContractAddress
	symbol := strings.ToLower(event.Symbol)
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		CosmosSender:         cosmosSender,
		EthereumReceiver:     ethereumReceiver,
		TokenContractAddress: tokenContractAddress,
		Symbol:               symbol,
		Amount:               amount,
	}

	return prophecyClaim
}

// getSymbolAmountFromCoin : Parse (symbol, amount) from coin string
func getSymbolAmountFromCoin(coin string) (string, *big.Int) {
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
