package txs

import (
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/x/ethbridge/types"
)

func TestLogLockToEthBridgeClaim(t *testing.T) {
	// Set up testing variables
	testBridgeContractAddress := types.NewEthereumAddress(TestBridgeContractAddress)
	testTokenContractAddress := types.NewEthereumAddress(TestTokenAddress)
	testEthereumAddress := types.NewEthereumAddress(TestEthereumAddress1)
	// Cosmos account address
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress1)
	require.NoError(t, err)
	// Cosmos validator address
	testRawCosmosValidatorAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress2)
	require.NoError(t, err)
	testCosmosValidatorBech32Address := sdk.ValAddress(testRawCosmosValidatorAddress)
	// Construct coins from TestAmount and TestSymbol
	coins := strconv.Itoa(TestAmount) + TestSymbol
	testCoins, err := sdk.ParseCoins(coins)
	require.NoError(t, err)

	// Set up expected EthBridgeClaim
	expectedEthBridgeClaim := types.NewEthBridgeClaim(TestEthereumChainID, testBridgeContractAddress, TestNonce, TestSymbol, testTokenContractAddress, testEthereumAddress, testCosmosAddress, testCosmosValidatorBech32Address, testCoins, TestLockClaimType)

	// Create test LogLockEvent
	logLockEvent := CreateTestLogLockEvent(t)

	ethBridgeClaim, err := LogLockToEthBridgeClaim(testCosmosValidatorBech32Address, &logLockEvent)
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
	prefixedHashedMsg := PrepareMsgForSigning(message.Hex())

	// Sign the message using the validator's private key
	signature, err := SignClaim(prefixedHashedMsg, privateKey)
	require.NoError(t, err)

	// Set up expected OracleClaim
	expectedOracleClaim := OracleClaim{
		ProphecyID: big.NewInt(int64(TestProphecyID)),
		Message:    message.Hex(),
		Signature:  signature,
	}

	// Map the test ProphecyClaim to a signed OracleClaim
	oracleClaim, err := ProphecyClaimToSignedOracleClaim(prophecyClaimEvent, privateKey)
	require.NoError(t, err)

	require.Equal(t, expectedOracleClaim, oracleClaim)
}

func TestBurnEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgBurn
	expectedMsgBurn := CreateTestCosmosMsg(t, events.MsgBurn)

	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t)
	msgBurn := BurnLockEventToCosmosMsg(events.MsgBurn, cosmosMsgAttributes)

	require.Equal(t, msgBurn, expectedMsgBurn)
}

func TestLockEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgLock
	expectedMsgLock := CreateTestCosmosMsg(t, events.MsgLock)

	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t)
	msgLock := BurnLockEventToCosmosMsg(events.MsgLock, cosmosMsgAttributes)

	require.Equal(t, expectedMsgLock, msgLock)
}

func TestMsgBurnToProphecyClaim(t *testing.T) {
	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:            events.MsgBurn,
		CosmosSender:         []byte(TestCosmosAddress1),
		EthereumReceiver:     common.HexToAddress(TestEthereumAddress1),
		TokenContractAddress: common.HexToAddress(TestTokenAddress),
		Symbol:               TestSymbol,
		Amount:               big.NewInt(int64(TestAmount)),
	}

	// Create a MsgBurn as input parameter
	testCosmosMsgBurn := CreateTestCosmosMsg(t, events.MsgBurn)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgBurn)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}

func TestMsgLockToProphecyClaim(t *testing.T) {
	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:            events.MsgLock,
		CosmosSender:         []byte(TestCosmosAddress1),
		EthereumReceiver:     common.HexToAddress(TestEthereumAddress1),
		TokenContractAddress: common.HexToAddress(TestTokenAddress),
		Symbol:               TestSymbol,
		Amount:               big.NewInt(int64(TestAmount)),
	}

	// Create a MsgLock as input parameter
	testCosmosMsgLock := CreateTestCosmosMsg(t, events.MsgLock)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgLock)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}
