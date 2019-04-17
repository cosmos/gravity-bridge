package common

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestAddress         = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator       = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce           = 0
	TestEthereumAddress = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestID              = "ethereumAddress0"
)

func CreateTestProphecy(t *testing.T) types.BridgeProphecy {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	bridgeClaim := types.NewBridgeClaim(TestID, testAddress, testValidator, amount)
	bridgeClaims := []types.BridgeClaim{bridgeClaim}
	newProphecy := types.NewBridgeProphecy(TestID, types.PendingStatus, 5, bridgeClaims)
	return newProphecy
}

func CreateTestEthMsg(t *testing.T) types.MsgMakeBridgeEthClaim {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	ethMsg := types.NewMsgMakeEthBridgeClaim(TestNonce, TestEthereumAddress, testAddress, testValidator, amount)
	return ethMsg
}
