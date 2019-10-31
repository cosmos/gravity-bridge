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
	Signature  []byte
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
func SignHash(hash common.Hash) []byte {
	// Load the validator's private key
	privateKey, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Signing hash:", hash.Hex())

	// Sign the hash using the validator's private key
	rawSignature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Verifying raw signature...")
	verifySig(hash, rawSignature)

	fmt.Println("\nRecovering signature components...")
	r, s, v := SigRSV(rawSignature)
	fmt.Println("r:", hexutil.Encode(r[:])[2:])
	fmt.Println("s:", hexutil.Encode(s[:])[2:])
	fmt.Println("v:", v)

	return rawSignature
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
	hash.Bytes()

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

// type Sig struct {
// 	Raw  []byte
// 	Hash [32]byte
// 	R    [32]byte
// 	S    [32]byte
// 	V    uint8
// }

// func Sign(message string) Sig {

// 	privateKey, err := LoadPrivateKey()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	hashRaw := crypto.Keccak256([]byte(message))
// 	signature, err := crypto.Sign(hashRaw, privateKey)

// 	// Convert hash to [32]byte
// 	var hashRaw32 [32]byte
// 	copy(hashRaw32[:], hashRaw)

// 	// s := string(byteArray[:n])

// 	// return Sig{
// 	// 	signature,
// 	// 	hashRaw32,
// 	// 	[32]byte(signature[:32]),
// 	// 	p.bytes32(signature[32:64]),
// 	// 	uint8(int(signature[65])) + 27, // Yes add 27, weird Ethereum quirk
// 	// }
// }
