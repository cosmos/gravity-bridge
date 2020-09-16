package ethbridge

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/trinhtan/peggy/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trinhtan/peggy/x/ethbridge/types"
	"github.com/stretchr/testify/require"
)

const (
	moduleString = "module"
	statusString = "status"
	senderString = "sender"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	//Unrecognized type
	res, err := handler(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "module":
				require.Equal(t, value, types.ModuleName)
			case senderString:
				require.Equal(t, value, valAddress.String())
			case "ethereum_sender":
				require.Equal(t, value, types.TestEthereumAddress)
			case "cosmos_receiver":
				require.Equal(t, value, types.TestAddress)
			case "amount":
				require.Equal(t, value, strconv.FormatInt(types.TestCoinsAmount, 10))
			case "symbol":
				require.Equal(t, value, types.TestCoinsSymbol)
			case "token_contract_address":
				require.Equal(t, value, types.TestTokenContractAddress)
			case statusString:
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
	err = badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res, err := handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Duplicate message from same validator
	res, err = handler(ctx, normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2, types.LockText)
	res, err := handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal2Pow7, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins := sdk.Coins{sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinsAmount)}
	require.True(t, receiverCoins.IsEqual(expectedCoins))
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.SuccessStatusText])
			}
		}
	}

	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
	receiverCoins = bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins = sdk.Coins{sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinsAmount)}
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

	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow3, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal2Pow4, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.LockText)
	ethMsg2 := NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal3Pow3, testEthereumAddress, types.AltTestCoinsAmount, types.AltTestCoinsSymbol, types.LockText)
	ethMsg3 := NewMsgCreateEthBridgeClaim(ethClaim3)

	//Initial message
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from second validator succeeds
	res, err = handler(ctx, ethMsg2)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from third validator succeeds but results in failed prophecy with no minting
	res, err = handler(ctx, ethMsg3)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
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
	coinsToMintAmount := int64(7)
	coinsToMintSymbol := "ether"
	coinsToMintSymbolLocked := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToMintSymbol)

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	mintedCoins := sdk.Coins{sdk.NewInt64Coin(coinsToMintSymbolLocked, coinsToMintAmount)}
	require.True(t, receiverCoins.IsEqual(mintedCoins))

	coinsToBurnAmount := int64(3)
	coinsToBurnSymbol := "ether"
	coinsToBurnSymbolPrefixed := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToBurnSymbol)

	ethereumReceiver := types.NewEthereumAddress(types.AltTestEthereumAddress)

	// Second message succeeds, burns eth and fires correct event
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount,
		coinsToBurnSymbolPrefixed)
	res, err = handler(ctx, burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	senderAddress := receiverAddress
	burnedCoins := sdk.Coins{sdk.NewInt64Coin(coinsToBurnSymbolPrefixed, coinsToBurnAmount)}
	remainingCoins := mintedCoins.Sub(burnedCoins)
	senderCoins := bankKeeper.GetCoins(ctx, senderAddress)
	require.True(t, senderCoins.IsEqual(remainingCoins))
	eventEthereumChainID := ""
	eventCosmosSender := ""
	eventEthereumReceiver := ""
	eventAmount := ""
	eventSymbol := ""
	eventCoins := ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case senderString:
				require.Equal(t, value, senderAddress.String())
			case "recipient":
				require.Equal(t, value, moduleAccountAddress.String())
			case moduleString:
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value
			case "cosmos_sender":
				eventCosmosSender = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "amount":
				eventAmount = value
			case "symbol":
				eventSymbol = value
			case "coins":
				eventCoins = value
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventAmount, strconv.FormatInt(coinsToBurnAmount, 10))
	require.Equal(t, eventSymbol, coinsToBurnSymbolPrefixed)
	require.Equal(t, eventCoins, sdk.Coins{sdk.NewInt64Coin(coinsToBurnSymbolPrefixed, coinsToBurnAmount)}.String())

	// Third message succeeds, burns more eth and fires correct event
	res, err = handler(ctx, burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	remainingCoins = remainingCoins.Sub(burnedCoins)
	senderCoins = bankKeeper.GetCoins(ctx, senderAddress)
	require.True(t, senderCoins.IsEqual(remainingCoins))
	eventEthereumChainID = ""
	eventCosmosSender = ""
	eventEthereumReceiver = ""
	eventAmount = ""
	eventSymbol = ""
	eventCoins = ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case senderString:
				require.Equal(t, value, senderAddress.String())
			case "recipient":
				require.Equal(t, value, moduleAccountAddress.String())
			case moduleString:
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value

			case "cosmos_sender":
				eventCosmosSender = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "amount":
				eventAmount = value
			case "symbol":
				eventSymbol = value
			case "coins":
				eventCoins = value
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventAmount, strconv.FormatInt(coinsToBurnAmount, 10))
	require.Equal(t, eventSymbol, coinsToBurnSymbolPrefixed)
	require.Equal(t, eventCoins, sdk.Coins{sdk.NewInt64Coin(coinsToBurnSymbolPrefixed, coinsToBurnAmount)}.String())

	// Fourth message fails, not enough eth
	res, err = handler(ctx, burnMsg)
	require.Error(t, err)
	require.Nil(t, res)
	senderCoins = bankKeeper.GetCoins(ctx, senderAddress)
	require.True(t, senderCoins.IsEqual(remainingCoins))
}
