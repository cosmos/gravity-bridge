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
	"strings"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	ethbridgeTypes "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// ProphecyClaimToSignedOracleClaim : packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToSignedOracleClaim(event events.NewProphecyClaimEvent) OracleClaim {
	// Parse relevant data into type byte[]
	// prophecyID := event.ProphecyID.Bytes()
	// sender := event.CosmosSender
	// recipient := []byte(event.EthereumReceiver.Hex())
	// token := []byte(event.TokenAddress.Hex())
	// amount := event.Amount.Bytes()
	// validator := []byte(event.ValidatorAddress.Hex())

	// Generate hash using ProphecyClaim data
	// claimHash := GenerateClaimHash(prophecyID, sender, recipient, token, amount, validator)

	// Sign the hash using the active validator's private key
	// hash, v, r, s := SignHash(claimHash)
	hash, r, s, v := Sign("hello")

	data := []byte("hello")
	SignFull(data)

	// Convert claimHash to [32]byte for packaging in OracleClaim
	// var byteHash [32]byte
	// copy(byteHash[:], hash)

	// Package the ProphecyID, Message, and Signature into an OracleClaim
	oracleClaim := OracleClaim{
		ProphecyID: event.ProphecyID,
		Message:    hash,
		V:          v,
		R:          r,
		S:          s,
		// Signature:  signature,
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
