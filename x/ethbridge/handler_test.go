package ethbridge

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/peggy/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/peggy/x/ethbridge/types"
	keeperLib "github.com/cosmos/peggy/x/oracle/keeper"
	"github.com/stretchr/testify/require"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	cdc := keeperLib.MakeTestCodec()

	valAddress := validatorAddresses[0]

	handler := NewHandler(keeper, bankKeeper, types.DefaultCodespace, cdc)

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress)
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
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, valAddress)
	badCreateMsg.Nonce = -1
	err := badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	cdc := keeperLib.MakeTestCodec()

	valAddress := validatorAddresses[0]

	handler := NewHandler(keeper, bankKeeper, types.DefaultCodespace, cdc)
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress)
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
	ctx, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{2, 7, 1})
	cdc := keeperLib.MakeTestCodec()

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	handler := NewHandler(keeper, bankKeeper, types.DefaultCodespace, cdc)

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal2Pow7)
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
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1)
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
	ctx, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.71, []int64{3, 4, 3})
	cdc := keeperLib.MakeTestCodec()

	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]

	ethClaim1 := types.CreateTestEthClaim(t, valAddressVal1Pow3, types.NewEthereumAddress(types.TestEthereumAddress), types.TestCoins)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(t, valAddressVal2Pow4, types.NewEthereumAddress(types.AltTestEthereumAddress), types.TestCoins)
	ethMsg2 := NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(t, valAddressVal3Pow3, types.NewEthereumAddress(types.TestEthereumAddress), types.AltTestCoins)
	ethMsg3 := NewMsgCreateEthBridgeClaim(ethClaim3)

	handler := NewHandler(keeper, bankKeeper, types.DefaultCodespace, cdc)

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
