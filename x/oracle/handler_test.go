package oracle

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	keeperLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, _, keeper := keeperLib.CreateTestKeepers(t, false, 0)
	handler := NewHandler(keeper)

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized oracle message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t)
	badCreateMsg.Nonce = -1
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum nonce provided"))

	badCreateMsg = types.CreateTestEthMsg(t)
	badCreateMsg.EthereumSender = "badAddress"
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum address provided"))
}

func TestDuplicateMsgs(t *testing.T) {
	//Setup
	ctx, _, keeper := keeperLib.CreateTestKeepers(t, false, 0)
	handler := NewHandler(keeper)
	normalCreateMsg := types.CreateTestEthMsg(t)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Duplicate message from same validator
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Not yet implemented"))

}
