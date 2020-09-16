package txs

import (
	"crypto/ecdsa"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	solsha3 "github.com/miguelmota/go-solidity-sha3"

	"github.com/trinhtan/peggy/cmd/ebrelayer/types"
)

// LoadPrivateKey loads the validator's private key from environment variables
func LoadPrivateKey() (key *ecdsa.PrivateKey, err error) {
	// Load config file containing environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
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

// LoadSender uses the validator's private key to load the validator's address
func LoadSender() (address common.Address, err error) {
	key, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, nil
}

// GenerateClaimMessage Generates a hashed message containing a ProphecyClaim event's data
func GenerateClaimMessage(event types.ProphecyClaimEvent) []byte {
	prophecyID := solsha3.Int256(event.ProphecyID)
	sender := solsha3.String(event.CosmosSender)
	recipient := solsha3.Int256(event.EthereumReceiver.Hex())
	token := solsha3.String(event.TokenAddress.Hex())
	amount := solsha3.Int256(event.Amount)

	// Generate claim message using ProphecyClaim data
	return solsha3.SoliditySHA3(prophecyID, sender, recipient, token, amount)
}

// PrefixMsg prefixes a message for verification, mimics behavior of web3.eth.sign
func PrefixMsg(msg []byte) []byte {
	return solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), msg)
}

// SignClaim Signs the prepared message with validator's private key
func SignClaim(msg []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	// Sign the message
	sig, err := crypto.Sign(msg, key)
	if err != nil {
		panic(err)
	}
	return sig, nil
}
