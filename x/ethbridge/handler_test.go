package ethbridge

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/peggy/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/stretchr/testify/require"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "module":
				require.Equal(t, value, types.ModuleName)
			case "sender":
				require.Equal(t, value, valAddress.String())
			case "ethereum_sender":
				require.Equal(t, value, types.TestEthereumAddress)
			case "cosmos_receiver":
				require.Equal(t, value, types.TestAddress)
			case "amount":
				require.Equal(t, value, types.TestCoins)
			case "status":
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			case "claim_type":
				require.Equal(t, value, types.ClaimTypeToString[types.LockText])
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	badCreateMsg.Nonce = -1
	err := badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == "status" {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Duplicate message from same validator
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "already processed message from validator for this id"))
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2, types.LockText)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal2Pow7, types.LockText)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins, err := sdk.ParseCoins(types.TestCoins)
	require.NoError(t, err)
	require.True(t, receiverCoins.IsEqual(expectedCoins))
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == "status" {
				require.Equal(t, value, oracle.StatusTextToString[oracle.SuccessStatusText])
			}
		}
	}

	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1, types.LockText)
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "prophecy already finalized"))
	receiverCoins = bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins, err = sdk.ParseCoins(types.TestCoins)
	require.NoError(t, err)
	require.True(t, receiverCoins.IsEqual(expectedCoins))
}

func TestNoMintFail(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.71, []int64{3, 4, 3})

	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow3, testEthereumAddress, types.TestCoins, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal2Pow4, testEthereumAddress, types.TestCoins, types.LockText)
	ethMsg2 := NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal3Pow3, testEthereumAddress, types.AltTestCoins, types.LockText)
	ethMsg3 := NewMsgCreateEthBridgeClaim(ethClaim3)

	//Initial message
	res := handler(ctx, ethMsg1)
	require.True(t, res.IsOK())
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == "status" {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from second validator succeeds
	res = handler(ctx, ethMsg2)
	require.True(t, res.IsOK())
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == "status" {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from third validator succeeds but results in failed prophecy with no minting
	res = handler(ctx, ethMsg3)
	require.True(t, res.IsOK())
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == "status" {
				require.Equal(t, value, oracle.StatusTextToString[oracle.FailedStatusText])
			}
		}
	}
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1Coins := bankKeeper.GetCoins(ctx, receiverAddress)
	require.True(t, receiver1Coins.IsZero())
}

func TestBurnEthFail(t *testing.T) {

}

func TestBurnEthSuccess(t *testing.T) {
	ctx, _, bankKeeper, supplyKeeper, _, validatorAddresses, handler := CreateTestHandler(t, 0.5, []int64{5})
	valAddressVal1Pow5 := validatorAddresses[0]

	moduleAccount := supplyKeeper.GetModuleAccount(ctx, ModuleName)
	moduleAccountAddress := moduleAccount.GetAddress()

	// Initial message to mint some eth
	coinsToMint := "7ethereum"

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow5, testEthereumAddress, coinsToMint, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res := handler(ctx, ethMsg1)
	require.True(t, res.IsOK())
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	mintedCoins, err := sdk.ParseCoins(coinsToMint)
	require.NoError(t, err)
	require.True(t, receiverCoins.IsEqual(mintedCoins))

	coinsToBurn := "3ethereum"
	ethereumReceiver := types.NewEthereumAddress(types.AltTestEthereumAddress)

	// Second message succeeds, burns eth and fires correct event
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, coinsToBurn)
	res = handler(ctx, burnMsg)
	require.True(t, res.IsOK())
	senderAddress := receiverAddress
	require.NoError(t, err)
	burnedCoins, err := sdk.ParseCoins(coinsToBurn)
	require.NoError(t, err)
	remainingCoins := mintedCoins.Sub(burnedCoins)
	senderCoins := bankKeeper.GetCoins(ctx, senderAddress)
	require.NoError(t, err)
	require.True(t, senderCoins.IsEqual(remainingCoins))
	eventEthereumChainID := ""
	eventTokenContract := ""
	eventCosmosSender := ""
	eventEthereumReceiver := ""
	eventAmount := ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "sender":
				require.Equal(t, value, senderAddress.String())
			case "recipient":
				require.Equal(t, value, moduleAccountAddress.String())
			case "module":
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value
			case "token_contract_address":
				eventTokenContract = value
			case "cosmos_sender":
				eventCosmosSender = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "amount":
				eventAmount = value
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventTokenContract, types.TestTokenContractAddress)
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventAmount, coinsToBurn)

	// Third message succeeds, burns more eth and fires correct event
	res = handler(ctx, burnMsg)
	require.True(t, res.IsOK())
	require.NoError(t, err)
	remainingCoins = remainingCoins.Sub(burnedCoins)
	senderCoins = bankKeeper.GetCoins(ctx, senderAddress)
	require.NoError(t, err)
	require.True(t, senderCoins.IsEqual(remainingCoins))
	eventEthereumChainID = ""
	eventTokenContract = ""
	eventCosmosSender = ""
	eventEthereumReceiver = ""
	eventAmount = ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "sender":
				require.Equal(t, value, senderAddress.String())
			case "recipient":
				require.Equal(t, value, moduleAccountAddress.String())
			case "module":
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value
			case "token_contract_address":
				eventTokenContract = value
			case "cosmos_sender":
				eventCosmosSender = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "amount":
				eventAmount = value
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventTokenContract, types.TestTokenContractAddress)
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventAmount, coinsToBurn)

	// Fourth message fails, not enough eth
	res = handler(ctx, burnMsg)
	require.False(t, res.IsOK())
	require.NoError(t, err)
	senderCoins = bankKeeper.GetCoins(ctx, senderAddress)
	require.NoError(t, err)
	require.True(t, senderCoins.IsEqual(remainingCoins))
}
