package types

import (
	"testing"

	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestAddress         = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator       = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce           = 0
	TestEthereumAddress = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
)

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T) MsgMakeEthBridgeClaim {
	ethClaim := CreateTestEthClaim(t)
	ethMsg := NewMsgMakeEthBridgeClaim(ethClaim)
	return ethMsg
}

func CreateTestEthClaim(t *testing.T) EthBridgeClaim {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	ethClaim := NewEthBridgeClaim(TestNonce, TestEthereumAddress, testAddress, testValidator, amount)
	return ethClaim
}

func CreateTestQueryEthProphecyResponse(cdc *codec.Codec, t *testing.T) QueryEthProphecyResponse {
	ethBridgeClaim := CreateTestEthClaim(t)
	oracleClaim := CreateOracleClaimFromEthClaim(cdc, ethBridgeClaim)
	ethBridgeClaims := []EthBridgeClaim{ethBridgeClaim}
	resp := NewQueryEthProphecyResponse(oracleClaim.ID, oracle.PendingStatus, ethBridgeClaims)
	return resp
}
