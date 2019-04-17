package oracle

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/common"
	keeperLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
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
	normalCreateMsg := common.CreateTestEthMsg(t)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Bad Creation
	badCreateMsg := common.CreateTestEthMsg(t)
	badCreateMsg.Nonce = -1
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum nonce provided"))

	badCreateMsg = common.CreateTestEthMsg(t)
	badCreateMsg.EthereumSender = "badAddress"
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum address provided"))
}
