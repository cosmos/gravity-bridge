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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmCommon "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/utils"
	ethbridgeTypes "github.com/cosmos/peggy/x/ethbridge/types"
)

// LogLockToEthBridgeClaim : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func LogLockToEthBridgeClaim(valAddr sdk.ValAddress, event *events.LockEvent) (ethbridgeTypes.EthBridgeClaim, error) {
	witnessClaim := ethbridgeTypes.EthBridgeClaim{}

	// chainID type casting (*big.Int -> int)
	chainID := int(event.EthereumChainID.Int64())

	bridgeContractAddress := ethbridgeTypes.NewEthereumAddress(event.BridgeContractAddress.Hex())

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

	// Sender type casting (address.common -> string)
	tokenContractAddress := ethbridgeTypes.NewEthereumAddress(event.Token.Hex())

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && !utils.IsZeroAddress(event.Token) {
		return witnessClaim, errors.New("symbol \"eth\" must have null address set as token address")
	}

	// Amount type casting (*big.Int -> sdk.Coins)
	coins := sdk.Coins{sdk.NewInt64Coin(symbol, event.Value.Int64())}

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Int64())

	// Package the information in a unique EthBridgeClaim
	witnessClaim.EthereumChainID = chainID
	witnessClaim.BridgeContractAddress = bridgeContractAddress
	witnessClaim.Nonce = nonce
	witnessClaim.TokenContractAddress = tokenContractAddress
	witnessClaim.Symbol = symbol
	witnessClaim.EthereumSender = sender
	witnessClaim.ValidatorAddress = valAddr
	witnessClaim.CosmosReceiver = recipient
	witnessClaim.Amount = coins

	return witnessClaim, nil
}

// ProphecyClaimToSignedOracleClaim : packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToSignedOracleClaim(
	event events.NewProphecyClaimEvent,
	key *ecdsa.PrivateKey,
) (OracleClaim, error) {
	// Set up new OracleClaim struct
	oracleClaim := OracleClaim{}

	// Generate a hashed claim message which contains ProphecyClaim's data
	fmt.Println("Generating unique message for ProphecyClaim", event.ProphecyID)
	message := GenerateClaimMessage(event)

	// Prepare the message (required for signature verification on contract)
	prefixedHashedMsg := PrepareMsgForSigning(message.Hex())

	// Sign the message using the validator's private key
	fmt.Println("Signing message...")
	signature, err := SignClaim(prefixedHashedMsg, key)
	if err != nil {
		return oracleClaim, err
	}
	fmt.Println("Signature generated:", hexutil.Encode(signature))

	// Package the ProphecyID, Message, and Signature into an OracleClaim
	oracleClaim.ProphecyID = event.ProphecyID
	oracleClaim.Message = message.Hex()
	oracleClaim.Signature = signature

	return oracleClaim, nil
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
			symbol, amount = utils.GetSymbolAmountFromCoin(val)
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
