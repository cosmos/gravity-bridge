package txs

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// OracleClaim :
type OracleClaim struct {
	ProphecyID *big.Int
	Message    [32]byte
	Signature  []byte
}

// ProphecyClaim :
type ProphecyClaim struct {
	ClaimType            *big.Int
	CosmosSender         []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}

// uint8(msgData.ClaimType), msgData.CosmosSender, msgData.EthereumReceiver, msgData.TokenContractAddress, msgData.Symbol, msgData.Amount

// LoadSender : uses the validator's private key to load the validator's address
func LoadSender() (address common.Address, err error) {
	key, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Parse public key
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return fromAddress, nil
}

// LoadPrivateKey : loads the validator's private key from environment variables
func LoadPrivateKey() (key *ecdsa.PrivateKey, err error) {
	// Load config file containing environment variables
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Private key for validator's Ethereum address must be set as an environment variable
	rawPrivateKey := os.Getenv("ETHEREUM_PRIVATE_KEY")
	if strings.TrimSpace(rawPrivateKey) == "" {
		log.Fatal("Error loading ETHEREUM_PRIVATE_KEY from .env file")
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey, nil
}

// SignHash : signs a specified hash using the validator's private key
func SignHash(hash common.Hash) []byte {
	// Load the validator's private key
	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Signing hash:", hash.Hex())

	// Sign the hash using the validator's private key
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return signature
}

// ProphecyClaimToOracleClaim : packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToOracleClaim(event events.NewProphecyClaimEvent) OracleClaim {
	// Parse relevant data into type byte[]
	prophecyID := event.ProphecyID.Bytes()
	sender := event.From
	recipient := []byte(event.To.Hex())
	token := []byte(event.Token.Hex())
	amount := event.Amount.Bytes()
	validator := []byte(event.Validator.Hex())

	// Generate hash using ProphecyClaim data
	hash := crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)

	// Sign the hash using the active validator's private key
	signature := SignHash(hash)

	// Convert hash to [32]byte for packaging in OracleClaim
	var byteHash [32]byte
	copy(byteHash[:], hash.Hex())

	// Package the ProphecyID, Message, and Signature into an OracleClaim
	oracleClaim := OracleClaim{
		ProphecyID: event.ProphecyID,
		Message:    byteHash,
		Signature:  signature,
	}

	return oracleClaim
}

// CosmosMsgToProphecyClaim : parses event data from a CosmosMsg, packaging it as a ProphecyClaim
func CosmosMsgToProphecyClaim(event events.CosmosMsg) events.NewProphecyClaimEvent {

	// type CosmosMsg struct {
	// 	ClaimType            Event
	// 	CosmosSender         []byte
	// 	EthereumReceiver     common.Address
	// 	Symbol               string
	// 	Amount               *big.Int
	// 	TokenContractAddress common.Address
	// }

	// // Parse relevant data into type byte[]
	// prophecyID := event.ProphecyID.Bytes()
	// sender := event.From
	// recipient := []byte(event.To.Hex())
	// token := []byte(event.Token.Hex())
	// amount := event.Amount.Bytes()
	// validator := []byte(event.Validator.Hex())
}
