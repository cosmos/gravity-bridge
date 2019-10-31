package txs

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
)

// OracleClaim :
type OracleClaim struct {
	ProphecyID *big.Int
	Message    [32]byte
	V          uint8
	R          [32]byte
	S          [32]byte
}

// ProphecyClaim :
type ProphecyClaim struct {
	ClaimType            events.Event
	CosmosSender         []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
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
func GenerateClaimHash(prophecyID []byte, sender []byte, recipient []byte, token []byte, amount []byte, validator []byte) common.Hash {
	hash := crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)
	return hash
}

// SignHash : signs a specified hash using the validator's private key
func SignHash(hash []byte) ([32]byte, uint8, [32]byte, [32]byte) {
	// Load the validator's private key
	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Signing hash:", hash)

	// Sign the hash using the validator's private key
	rawSignature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nRecovering signature components...")
	r, s, v := SigRSV(rawSignature)

	fmt.Println("v:", v)
	fmt.Println("r:", hexutil.Encode(r[:])[2:])
	fmt.Println("s:", hexutil.Encode(s[:])[2:])

	var byteHash [32]byte
	copy(byteHash[:], hash)
	// byteHash := [32]byte{hash}

	fmt.Println("Verifying raw signature...")
	verifySig(byteHash, rawSignature)

	// return rawSignature
	return byteHash, v, r, s
}

func verifySig(hash common.Hash, signature []byte) {
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

// SigRSV signatures R S V returned as arrays
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

func SignFull(data []byte) {

	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	// SignFull :
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hexutil.Encode(signature)) // 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches) // true

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true
}

// Sign :
func Sign(message string) ([32]byte, [32]byte, [32]byte, uint8) {
	// Load the validator's private key
	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	hashRaw := crypto.Keccak256([]byte(message))
	signature, err := crypto.Sign(hashRaw, privateKey)

	var hash [32]byte
	copy(hash[:], hashRaw)

	var r [32]byte
	copy(r[:], signature[:32])

	var s [32]byte
	copy(s[:], signature[32:64])

	// v := uint8(int(signature[65])) + 27
	var v uint8
	rawV := int(signature[64])
	if rawV < 27 {
		v = uint8(rawV + 27)
	} else {
		v = uint8(rawV)
	}

	return hash, r, s, v
}
