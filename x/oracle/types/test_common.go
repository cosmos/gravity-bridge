package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestID = "ethereumAddress0"

	//Ethereum-bridge specific stuff
	TestAddress         = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator       = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce           = 0
	TestEthereumAddress = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
)

func CreateTestProphecy(t *testing.T) Prophecy {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	claim := NewClaim(TestID, testAddress, testValidator, amount)
	claims := []Claim{claim}
	newProphecy := NewProphecy(TestID, PendingStatus, 5, claims)
	return newProphecy
}

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T) MsgMakeEthBridgeClaim {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	ethMsg := NewMsgMakeEthBridgeClaim(TestNonce, TestEthereumAddress, testAddress, testValidator, amount)
	return ethMsg
}
