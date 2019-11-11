package txs

import (
	"crypto/ecdsa"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
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

// GenerateClaimMessage : Generates a hased message containing a ProphecyClaim event's data
func GenerateClaimMessage(event events.NewProphecyClaimEvent) common.Hash {
	// Cast event field values to byte[]
	prophecyID := event.ProphecyID.Bytes()
	sender := event.CosmosSender
	recipient := []byte(event.EthereumReceiver.Hex())
	token := []byte(event.TokenAddress.Hex())
	amount := event.Amount.Bytes()
	validator := []byte(event.ValidatorAddress.Hex())

	// Generate claim message using ProphecyClaim data
	return crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)
}

// PrepareMsgForSigning : prefixes a message for verification by a Smart Contract
func PrepareMsgForSigning(msg string) []byte {
	// Turn the message into a 32-byte hash
	hashedMsg := solsha3.SoliditySHA3(solsha3.String(msg))

	// Prefix and then hash to mimic behavior of eth_sign
	return solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), solsha3.Bytes32(hashedMsg))
}

// SignClaim : Signs the prepared message with validator's private key
func SignClaim(msg []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	// Sign the message
	sig, err := secp256k1.Sign(msg, math.PaddedBigBytes(key.D, 32))
	if err != nil {
		panic(err)
	}

	return sig, nil
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
