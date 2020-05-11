package txs

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/stretchr/testify/require"
)

func TestGenerateClaimMessage(t *testing.T) {
	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Confirm that the generated message matches the expected generated message
	require.Equal(t, TestExpectedMessage, hex.EncodeToString(message))
}

func TestPrepareMessageForSigning(t *testing.T) {
	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Simulate message hashing, prefixing
	prefixedMessage := solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), solsha3.Bytes32(message))

	// Prepare the message for signing
	preparedMessage := PrefixMsg(message)

	// Confirm that the prefixed message matches the prepared message
	require.Equal(t, preparedMessage, prefixedMessage)
}

func TestSignClaim(t *testing.T) {
	// Set and get env variables to replicate relayer
	os.Setenv(EthereumPrivateKey, TestPrivHex)
	rawKey := os.Getenv(EthereumPrivateKey)

	// Load signer's private key and address
	key, _ := crypto.HexToECDSA(rawKey)
	signerAddr := common.HexToAddress(TestAddrHex)

	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)

	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Prepare the message (required for signature verification on contract)
	prefixedHashedMsg := PrefixMsg(message)

	// Sign the message using the validator's private key
	signature, err := SignClaim(prefixedHashedMsg, key)
	require.NoError(t, err)

	// Recover signer's public key and address
	recoveredPub, err := crypto.Ecrecover(prefixedHashedMsg, signature)
	require.NoError(t, err)
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Confirm that the recovered address is correct
	require.Equal(t, recoveredAddr, signerAddr)
}
