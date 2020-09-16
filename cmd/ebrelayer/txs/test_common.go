package txs

import (
	"encoding/binary"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	tmKv "github.com/tendermint/tendermint/libs/kv"

	"github.com/trinhtan/peggy/cmd/ebrelayer/types"
	ethbridge "github.com/trinhtan/peggy/x/ethbridge/types"
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
	TestSymbol                = "PEGGYETH"
	TestAmount                = 5
	TestEthereumAddress1      = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestEthereumAddress2      = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestCosmosAddress1        = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestCosmosAddress2        = "cosmos1l5h2x255pvdy9l4z0hf9tr8zw7k657s97wyz7y"
	TestExpectedMessage       = "d39d3a837b322ea6355a4de856bb88e0a1a33a1a9655767df2fa947f42ebcda6"
	TestExpectedSignature     = "f3b43b87b8b3729d6b380a640954d14e425acd603bc98f5da8437cba9e492e7805c437f977900cf9cfbeb9e0e2e6fc5b189723b0979efff1fc2796a2daf4fd3a01" //nolint:lll
	TestAddrHex               = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex               = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	TestNullAddress           = "0x0000000000000000000000000000000000000000"
	TestOtherAddress          = "0x1000000000000000000000000000000000000000"
)

// CreateTestLogEthereumEvent creates a sample EthereumEvent event for testing purposes
func CreateTestLogEthereumEvent(t *testing.T) types.EthereumEvent {
	testEthereumChainID := big.NewInt(int64(TestEthereumChainID))
	testBridgeContractAddress := common.HexToAddress(TestBridgeContractAddress)
	// Convert int to [32]byte
	var testProphecyID []byte
	var testProphecyID32 [32]byte
	testProphecyID = make([]byte, 32)
	binary.LittleEndian.PutUint64(testProphecyID, uint64(TestProphecyID))
	copy(testProphecyID32[:], testProphecyID)
	testEthereumSender := common.HexToAddress(TestEthereumAddress1)
	testCosmosRecipient := []byte(TestCosmosAddress1)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := big.NewInt(int64(TestAmount))
	testNonce := big.NewInt(int64(TestNonce))

	return types.EthereumEvent{testEthereumChainID, testBridgeContractAddress,
		testProphecyID32, testEthereumSender, testCosmosRecipient, testTokenAddress,
		TestSymbol, testAmount, testNonce, ethbridge.LockText}
}

// CreateTestProphecyClaimEvent creates a sample ProphecyClaimEvent for testing purposes
func CreateTestProphecyClaimEvent(t *testing.T) types.ProphecyClaimEvent {
	testProphecyID := big.NewInt(int64(TestProphecyID))
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testValidatorAddress := common.HexToAddress(TestEthereumAddress2)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := big.NewInt(int64(TestAmount))

	return types.NewProphecyClaimEvent([]byte(TestCosmosAddress1), TestSymbol,
		testProphecyID, testAmount, testEthereumReceiver, testValidatorAddress,
		testTokenAddress, TestBurnClaimType)
}

// CreateTestCosmosMsg creates a sample Cosmos Msg for testing purposes
func CreateTestCosmosMsg(t *testing.T, claimType types.Event) types.CosmosMsg {
	testCosmosSender := []byte(TestCosmosAddress1)
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testAmount := big.NewInt(int64(TestAmount))

	var symbol string
	if claimType == types.MsgBurn {
		res := strings.SplitAfter(TestSymbol, "PEGGY")
		symbol = strings.Join(res[1:], "")
	} else {
		symbol = TestSymbol
	}

	// Create new Cosmos Msg
	cosmosMsg := types.NewCosmosMsg(claimType, testCosmosSender,
		testEthereumReceiver, symbol, testAmount)

	return cosmosMsg
}

// CreateCosmosMsgAttributes creates expected attributes for a MsgBurn/MsgLock for testing purposes
func CreateCosmosMsgAttributes(t *testing.T, claimType types.Event) []tmKv.Pair {
	attributes := [5]tmKv.Pair{}

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

	// (key, value) pairing for "symbol" key
	var symbol string
	if claimType == types.MsgBurn {
		symbol = strings.ToLower(TestSymbol)
	} else {
		symbol = TestSymbol
	}
	pairSymbol := tmKv.Pair{
		Key:   []byte("symbol"),
		Value: []byte(symbol),
	}

	// (key, value) pairing for "amount" key
	pairAmount := tmKv.Pair{
		Key:   []byte("amount"),
		Value: []byte(strconv.Itoa(TestAmount)),
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
	attributes[3] = pairSymbol
	attributes[4] = pairAmount

	return attributes[:]
}
