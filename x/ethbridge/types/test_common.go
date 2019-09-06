package types

import (
	"testing"

	"github.com/cosmos/peggy/x/oracle"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestAddress            = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator          = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce              = 0
	TestEthereumAddress    = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	AltTestEthereumAddress = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207344"
	TestCoins              = "10ethereum"
	AltTestCoins           = "12ethereum"
)

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T, validatorAddress sdk.ValAddress) MsgCreateEthBridgeClaim {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	ethClaim := CreateTestEthClaim(t, validatorAddress, testEthereumAddress, TestCoins)
	ethMsg := NewMsgCreateEthBridgeClaim(ethClaim)
	return ethMsg
}

func CreateTestEthClaim(t *testing.T, validatorAddress sdk.ValAddress, testEthereumAddress EthereumAddress, coins string) EthBridgeClaim {
	testCosmosAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	amount, err2 := sdk.ParseCoins(coins)
	require.NoError(t, err1)
	require.NoError(t, err2)
	ethClaim := NewEthBridgeClaim(TestNonce, testEthereumAddress, testCosmosAddress, validatorAddress, amount)
	return ethClaim
}

func CreateTestQueryEthProphecyResponse(cdc *codec.Codec, t *testing.T, validatorAddress sdk.ValAddress) QueryEthProphecyResponse {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	ethBridgeClaim := CreateTestEthClaim(t, validatorAddress, testEthereumAddress, TestCoins)
	oracleClaim, _ := CreateOracleClaimFromEthClaim(cdc, ethBridgeClaim)
	ethBridgeClaims := []EthBridgeClaim{ethBridgeClaim}

	return NewQueryEthProphecyResponse(
		oracleClaim.ID,
		oracle.NewStatus(oracle.PendingStatusText, ""),
		ethBridgeClaims,
	)
}
