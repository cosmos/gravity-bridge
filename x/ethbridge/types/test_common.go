package types

import (
	"testing"

	"github.com/trinhtan/peggy/x/oracle"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestEthereumChainID       = 3
	TestBridgeContractAddress = "0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB"
	TestAddress               = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator             = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce                 = 0
	TestTokenContractAddress  = "0x0000000000000000000000000000000000000000"
	TestEthereumAddress       = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	AltTestEthereumAddress    = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207344"
	TestCoinsAmount           = 10
	TestCoinsSymbol           = "eth"
	TestCoinsLockedSymbol     = "peggyeth"
	AltTestCoinsAmount        = 12
	AltTestCoinsSymbol        = "eth"
)

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType) MsgCreateEthBridgeClaim {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethClaim := CreateTestEthClaim(
		t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	ethMsg := NewMsgCreateEthBridgeClaim(ethClaim)
	return ethMsg
}

func CreateTestEthClaim(
	t *testing.T, testContractAddress EthereumAddress, testTokenAddress EthereumAddress,
	validatorAddress sdk.ValAddress, testEthereumAddress EthereumAddress, amount int64, symbol string, claimType ClaimType,
) EthBridgeClaim {
	testCosmosAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err1)
	ethClaim := NewEthBridgeClaim(
		TestEthereumChainID, testContractAddress, TestNonce, symbol,
		testTokenAddress, testEthereumAddress, testCosmosAddress, validatorAddress, amount, claimType)
	return ethClaim
}

func CreateTestBurnMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount int64, coinsSymbol string) MsgBurn {
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	burnEth := NewMsgBurn(TestEthereumChainID, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol)
	return burnEth
}

func CreateTestQueryEthProphecyResponse(
	cdc *codec.Codec, t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType,
) QueryEthProphecyResponse {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethBridgeClaim := CreateTestEthClaim(t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	oracleClaim, _ := CreateOracleClaimFromEthClaim(cdc, ethBridgeClaim)
	ethBridgeClaims := []EthBridgeClaim{ethBridgeClaim}

	return NewQueryEthProphecyResponse(
		oracleClaim.ID,
		oracle.NewStatus(oracle.PendingStatusText, ""),
		ethBridgeClaims,
	)
}
