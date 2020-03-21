package txs

import (
	"math/big"
	"strconv"
	"testing"

	tmKv "github.com/tendermint/tendermint/libs/kv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
)

const (
	// EthereumPrivateKey config field which holds the user's private key
	EthereumPrivateKey        = "ETHEREUM_PRIVATE_KEY"
	TestEthereumChainID       = 3
	TestBridgeContractAddress = "0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"
	TestLockClaimType         = 0
	TestBurnClaimType         = 1
	TestProphecyID            = 20
	TestNonce                 = 19
	TestEthTokenAddress       = "0x0000000000000000000000000000000000000000"
	TestSymbol                = "eth"
	TestAmount                = 5
	TestEthereumAddress1      = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestEthereumAddress2      = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestCosmosAddress1        = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestCosmosAddress2        = "cosmos1l5h2x255pvdy9l4z0hf9tr8zw7k657s97wyz7y"
	TestExpectedMessage       = "0xfc3c746e966d5f48af553b166b0870b0fa6b6921b353fba67de4e2230392f48b"
	TestExpectedSignature     = "0xac349f2452d50d14e11f72de8fc7acde0b47f280a47792470198dcff59358e42425315c0db810dc5d2a7ba5eda7d9cf35cea4f13d550bfa03484df739249c4d401" //nolint:lll
	TestAddrHex               = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex               = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	TestNullAddress           = "0x0000000000000000000000000000000000000000"
	TestOtherAddress          = "0x1000000000000000000000000000000000000000"
)

// CreateTestLogLockEvent creates a sample LockEvent event for testing purposes
func CreateTestLogLockEvent(t *testing.T) events.LockEvent {
	testEthereumChainID := big.NewInt(int64(TestEthereumChainID))
	testBridgeContractAddress := common.HexToAddress(TestBridgeContractAddress)
	testEthereumSender := common.HexToAddress(TestEthereumAddress1)
	testCosmosRecipient := []byte(TestCosmosAddress1)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := big.NewInt(int64(TestAmount))
	testNonce := big.NewInt(int64(TestNonce))

	// Create new LockEvent
	lockEvent := events.LockEvent{
		EthereumChainID:       testEthereumChainID,
		BridgeContractAddress: testBridgeContractAddress,
		From:                  testEthereumSender,
		To:                    testCosmosRecipient,
		Token:                 testTokenAddress,
		Symbol:                TestSymbol,
		Value:                 testAmount,
		Nonce:                 testNonce,
	}

	return lockEvent
}

// CreateTestProphecyClaimEvent creates a sample ProphecyClaimEvent for testing purposes
func CreateTestProphecyClaimEvent(t *testing.T) events.NewProphecyClaimEvent {
	testProphecyID := big.NewInt(int64(TestProphecyID))
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testValidatorAddress := common.HexToAddress(TestEthereumAddress2)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := big.NewInt(int64(TestAmount))

	// Create new ProphecyClaimEvent
	prophecyClaimEvent := events.NewProphecyClaimEvent{
		ProphecyID:       testProphecyID,
		ClaimType:        TestBurnClaimType,
		EthereumReceiver: testEthereumReceiver,
		ValidatorAddress: testValidatorAddress,
		TokenAddress:     testTokenAddress,
		Symbol:           TestSymbol,
		Amount:           testAmount,
	}

	return prophecyClaimEvent
}

// CreateTestCosmosMsg creates a sample Cosmos Msg for testing purposes
func CreateTestCosmosMsg(t *testing.T, claimType events.Event) events.CosmosMsg {
	testCosmosSender := []byte(TestCosmosAddress1)
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testAmount := big.NewInt(int64(TestAmount))
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)

	// Create new Cosmos Msg
	cosmosMsg := events.NewCosmosMsg(
		claimType, testCosmosSender, testEthereumReceiver, TestSymbol, testAmount, testTokenAddress)

	return cosmosMsg
}

// CreateCosmosMsgAttributes creates expected attributes for a MsgBurn/MsgLock for testing purposes
func CreateCosmosMsgAttributes(t *testing.T) []tmKv.Pair {
	attributes := [4]tmKv.Pair{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := tmKv.Pair{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosAddress1),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := tmKv.Pair{
		Key:   []byte("ethereum_receiver"),
		Value: []byte(common.HexToAddress(TestEthereumAddress1).Hex()), // .Bytes() doesn't seem to work here
	}

	// (key, value) pairing for "amount" key
	pairAmount := tmKv.Pair{
		Key:   []byte("amount"),
		Value: []byte(strconv.Itoa(TestAmount) + TestSymbol),
	}

	// (key, value) pairing for "token_contract_address" key
	pairTokenContract := tmKv.Pair{
		Key:   []byte("token_contract_address"),
		Value: []byte(common.HexToAddress(TestEthTokenAddress).Hex()),
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairEthereumReceiver
	attributes[2] = pairTokenContract
	attributes[3] = pairAmount

	return attributes[:]
}
