package ethbridge

import (
	"strings"
	"testing"

	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	keeperLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses, _ := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	accAddress := sdk.AccAddress(validatorAddresses[0])

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, accAddress)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, accAddress)
	badCreateMsg.Nonce = -1
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum nonce provided"))

	badCreateMsg = types.CreateTestEthMsg(t, accAddress)
	badCreateMsg.EthereumSender = "badAddress"
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum address provided"))
}

func TestDuplicateMsgs(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses, _ := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	accAddress := sdk.AccAddress(validatorAddresses[0])

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)
	normalCreateMsg := types.CreateTestEthMsg(t, accAddress)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	require.Equal(t, res.Log, oracle.PendingStatus)

	//Duplicate message from same validator
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Already processed message from validator for this id"))

}

func TestMintSuccess(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses, _ := keeperLib.CreateTestKeepers(t, 0.7, []int64{2, 7, 1})
	accAddressVal1Pow2 := sdk.AccAddress(validatorAddresses[0])
	accAddressVal2Pow7 := sdk.AccAddress(validatorAddresses[1])
	accAddressVal3Pow1 := sdk.AccAddress(validatorAddresses[2])

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, accAddressVal1Pow2)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, accAddressVal2Pow7)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins, err := sdk.ParseCoins(types.TestCoins)
	require.NoError(t, err)
	require.True(t, receiverCoins.IsEqual(expectedCoins))
	require.Equal(t, res.Log, oracle.SuccessStatus)

	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, accAddressVal3Pow1)
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Prophecy already finalized"))
	receiverCoins = bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins, err = sdk.ParseCoins(types.TestCoins)
	require.NoError(t, err)
	require.True(t, receiverCoins.IsEqual(expectedCoins))

}

func TestNoMintFail(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses, _ := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 4, 3})
	accAddressVal1Pow3 := sdk.AccAddress(validatorAddresses[0])
	accAddressVal2Pow4 := sdk.AccAddress(validatorAddresses[1])
	accAddressVal3Pow3 := sdk.AccAddress(validatorAddresses[2])

	ethClaim1 := types.CreateTestEthClaim(t, accAddressVal1Pow3, types.TestEthereumAddress, types.TestCoins)
	ethMsg1 := NewMsgMakeEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(t, accAddressVal2Pow4, types.AltTestEthereumAddress, types.TestCoins)
	ethMsg2 := NewMsgMakeEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(t, accAddressVal3Pow3, types.TestEthereumAddress, types.AltTestCoins)
	ethMsg3 := NewMsgMakeEthBridgeClaim(ethClaim3)

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)

	//Initial message
	res := handler(ctx, ethMsg1)
	require.True(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, oracle.PendingStatus))
	require.Equal(t, res.Log, oracle.PendingStatus)

	//Different message from second validator succeeds
	res = handler(ctx, ethMsg2)
	require.True(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, oracle.PendingStatus))
	require.Equal(t, res.Log, oracle.PendingStatus)

	//Different message from third validator succeeds but results in failed prophecy with no minting
	res = handler(ctx, ethMsg3)
	require.True(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, oracle.FailedStatus))
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1Coins := bankKeeper.GetCoins(ctx, receiverAddress)
	require.True(t, receiver1Coins.IsZero())
}
