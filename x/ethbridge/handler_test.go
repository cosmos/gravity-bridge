package ethbridge

import (
	"strings"
	"testing"

	"github.com/swishlabsco/peggy/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/peggy/x/ethbridge/types"
	keeperLib "github.com/swishlabsco/peggy/x/oracle/keeper"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	valAddress := validatorAddresses[0]

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, valAddress)
	badCreateMsg.Nonce = -1
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum nonce provided"))
}

func TestDuplicateMsgs(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 7})
	valAddress := validatorAddresses[0]

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress)
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
	ctx, _, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{2, 7, 1})
	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	handler := NewHandler(keeper, bankKeeper, cdc, types.DefaultCodespace)

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
	require.Equal(t, res.Log, oracle.SuccessStatus)

	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1)
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
	ctx, _, keeper, bankKeeper, validatorAddresses := keeperLib.CreateTestKeepers(t, 0.7, []int64{3, 4, 3})
	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]

	ethClaim1 := types.CreateTestEthClaim(t, valAddressVal1Pow3, gethCommon.HexToAddress(types.TestEthereumAddress), types.TestCoins)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(t, valAddressVal2Pow4, gethCommon.HexToAddress(types.AltTestEthereumAddress), types.TestCoins)
	ethMsg2 := NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(t, valAddressVal3Pow3, gethCommon.HexToAddress(types.TestEthereumAddress), types.AltTestCoins)
	ethMsg3 := NewMsgCreateEthBridgeClaim(ethClaim3)

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
	require.Equal(t, res.Log, oracle.FailedStatus)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1Coins := bankKeeper.GetCoins(ctx, receiverAddress)
	require.True(t, receiver1Coins.IsZero())
}
