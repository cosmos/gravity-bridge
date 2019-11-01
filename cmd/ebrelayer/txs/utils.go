package txs

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	solsha3 "github.com/miguelmota/go-solidity-sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

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

// GenerateClaimHash : Generates an OracleClaim hash from a ProphecyClaim's event data
func GenerateClaimHash(prophecyID []byte, sender []byte, recipient []byte, token []byte, amount []byte, validator []byte) string {
	// Generate a hash containing the information
	rawHash := crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)

	// Cast hash to hex encoded string
	return rawHash.Hex()
}

// SignClaim : Signs hashed message with validator's private key
func SignClaim(hash string) []byte {
	key, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	signer := hex.EncodeToString(crypto.PubkeyToAddress(key.PublicKey).Bytes())
	fmt.Printf("\nAttempting to sign message \"%v\" with account \"%v\"...\n", hash, signer)

	rawSignature, _ := prefixMessage(hash, key)

	signature := hexutil.Bytes(rawSignature)
	fmt.Println("Success! Signature:", signature)

	return signature
}

func prefixMessage(message string, key *ecdsa.PrivateKey) ([]byte, []byte) {
	// Turn the message into a 32-byte hash
	hash := solsha3.SoliditySHA3(solsha3.String(message))
	// Prefix and then hash to mimic behavior of eth_sign
	prefixed := solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), solsha3.Bytes32(hash))
	sig, err := secp256k1.Sign(prefixed, math.PaddedBigBytes(key.D, 32))

	if err != nil {
		panic(err)
	}

	return sig, prefixed
}

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

// verifySignature: utility function for signature verification
func verifySignature(hash common.Hash, signature []byte) {
	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	// public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// public key which signed this message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)

	// compare
	matches := bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true
}

// SigRSV : utility function which breaks a signature down into [R, S, V] components
func SigRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := sigstr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}
