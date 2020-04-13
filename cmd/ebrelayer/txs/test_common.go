package txs

import (
	"encoding/binary"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	tmKv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/peggy/cmd/ebrelayer/types"
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
	TestAmount                = int64(5)
	TestEthereumAddress1      = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestEthereumAddress2      = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestCosmosAddress1        = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestCosmosAddress2        = "cosmos1l5h2x255pvdy9l4z0hf9tr8zw7k657s97wyz7y"
	TestExpectedMessage       = "0x34669ae046add4b9b45863caf90f623dcda0a869ad6163087fe7f9fc41f93355"
	TestExpectedSignature     = "0x2952d5757ea8d7d810e4c8ae5a4897a1e76a0d4741ee1468d0ce1d6e0ccf4d414512391a09952c4b1d26995f532e38f289256d35bd882e00a62c3ac9e8e9ef9601" //nolint:lll
	TestAddrHex               = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex               = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	TestNullAddress           = "0x0000000000000000000000000000000000000000"
	TestOtherAddress          = "0x1000000000000000000000000000000000000000"
)

// CreateTestLogLockEvent creates a sample LockEvent event for testing purposes
func CreateTestLogLockEvent(t *testing.T) types.LockEvent {
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

	return types.NewLockEvent(testEthereumChainID, testBridgeContractAddress,
		testProphecyID32, testEthereumSender, testCosmosRecipient, testTokenAddress,
		TestSymbol, testAmount, testNonce)
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
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)

	// Create new Cosmos Msg
	cosmosMsg := types.NewCosmosMsg(claimType, testCosmosSender,
		testEthereumReceiver, TestSymbol, testAmount)

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
		Value: []byte(strconv.FormatInt(TestAmount, 10) + TestSymbol),
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
