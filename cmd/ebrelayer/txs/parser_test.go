package txs

import (
	"math/big"
	"os"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/trinhtan/peggy/cmd/ebrelayer/types"
	ethbridge "github.com/trinhtan/peggy/x/ethbridge/types"
)

func TestLogLockToEthBridgeClaim(t *testing.T) {
	// Set up testing variables
	testBridgeContractAddress := ethbridge.NewEthereumAddress(TestBridgeContractAddress)
	testTokenContractAddress := ethbridge.NewEthereumAddress(TestEthTokenAddress)
	testEthereumAddress := ethbridge.NewEthereumAddress(TestEthereumAddress1)
	// Cosmos account address
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress1)
	require.NoError(t, err)
	// Cosmos validator address
	testRawCosmosValidatorAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress2)
	require.NoError(t, err)
	testCosmosValidatorBech32Address := sdk.ValAddress(testRawCosmosValidatorAddress)

	// Set up expected EthBridgeClaim
	expectedEthBridgeClaim := ethbridge.NewEthBridgeClaim(
		TestEthereumChainID, testBridgeContractAddress, TestNonce, strings.ToLower(TestSymbol), testTokenContractAddress,
		testEthereumAddress, testCosmosAddress, testCosmosValidatorBech32Address, TestAmount, TestLockClaimType)

	// Create test ethereum event
	ethereumEvent := CreateTestLogEthereumEvent(t)

	ethBridgeClaim, err := EthereumEventToEthBridgeClaim(testCosmosValidatorBech32Address, &ethereumEvent)
	require.NoError(t, err)

	require.Equal(t, expectedEthBridgeClaim, ethBridgeClaim)
}

func TestProphecyClaimToSignedOracleClaim(t *testing.T) {
	// Set ETHEREUM_PRIVATE_KEY env variable
	os.Setenv(EthereumPrivateKey, TestPrivHex)
	// Get and load private key from env variables
	rawKey := os.Getenv(EthereumPrivateKey)
	privateKey, _ := crypto.HexToECDSA(rawKey)

	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Prepare the message (required for signature verification on contract)
	prefixedHashedMsg := PrefixMsg(message)

	// Sign the message using the validator's private key
	signature, err := SignClaim(prefixedHashedMsg, privateKey)
	require.NoError(t, err)

	var message32 [32]byte
	copy(message32[:], message)

	// Set up expected OracleClaim
	expectedOracleClaim := OracleClaim{
		ProphecyID: big.NewInt(int64(TestProphecyID)),
		Message:    message32,
		Signature:  signature,
	}

	// Map the test ProphecyClaim to a signed OracleClaim
	oracleClaim, err := ProphecyClaimToSignedOracleClaim(prophecyClaimEvent, privateKey)
	require.NoError(t, err)

	require.Equal(t, expectedOracleClaim, oracleClaim)
}

func TestBurnEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgBurn
	expectedMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)

	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgBurn)
	msgBurn := BurnLockEventToCosmosMsg(types.MsgBurn, cosmosMsgAttributes)

	require.Equal(t, expectedMsgBurn, msgBurn)
}

func TestLockEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgLock
	expectedMsgLock := CreateTestCosmosMsg(t, types.MsgLock)

	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgLock)
	msgLock := BurnLockEventToCosmosMsg(types.MsgLock, cosmosMsgAttributes)

	require.Equal(t, expectedMsgLock, msgLock)
}

func TestMsgBurnToProphecyClaim(t *testing.T) {
	// Parse expected symbol
	res := strings.SplitAfter(TestSymbol, "PEGGY")
	symbol := strings.Join(res[1:], "")

	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:        types.MsgBurn,
		CosmosSender:     []byte(TestCosmosAddress1),
		EthereumReceiver: common.HexToAddress(TestEthereumAddress1),
		Symbol:           symbol,
		Amount:           big.NewInt(int64(TestAmount)),
	}

	// Create a MsgBurn as input parameter
	testCosmosMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgBurn)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}

func TestMsgLockToProphecyClaim(t *testing.T) {
	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:        types.MsgLock,
		CosmosSender:     []byte(TestCosmosAddress1),
		EthereumReceiver: common.HexToAddress(TestEthereumAddress1),
		Symbol:           TestSymbol,
		Amount:           big.NewInt(int64(TestAmount)),
	}

	// Create a MsgLock as input parameter
	testCosmosMsgLock := CreateTestCosmosMsg(t, types.MsgLock)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgLock)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}

func TestIsZeroAddress(t *testing.T) {
	falseRes := isZeroAddress(common.HexToAddress(TestOtherAddress))
	require.False(t, falseRes)

	trueRes := isZeroAddress(common.HexToAddress(TestNullAddress))
	require.True(t, trueRes)
}
