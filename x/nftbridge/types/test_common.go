package types

import (
	"testing"

	// "github.com/cosmos/peggy/x/nftbridge/types"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/oracle"
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
	TestSymbol                = "eth"
	TestTokenContractAddress  = "0x0000000000000000000000000000000000000000"
	TestEthereumAddress       = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	AltTestEthereumAddress    = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207344"
	TestDenom                 = "denom"
	TestID                    = "id1"
	AltTestID                 = "id2"
	AltTestDenom              = "denom2"
)

//Ethereum-bridge specific stuff
func CreateTestNFTMsg(t *testing.T, validatorAddress sdk.ValAddress, claimType ethbridge.ClaimType) MsgCreateNFTBridgeClaim {
	testEthereumAddress := ethbridge.NewEthereumAddress(TestEthereumAddress)
	testContractAddress := ethbridge.NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := ethbridge.NewEthereumAddress(TestTokenContractAddress)
	nftClaim := CreateTestNFTClaim(t, testContractAddress, testTokenAddress, validatorAddress, testEthereumAddress, TestDenom, TestID, claimType)
	nftMsg := NewMsgCreateNFTBridgeClaim(nftClaim)
	return nftMsg
}

func CreateTestNFTClaim(t *testing.T, testContractAddress ethbridge.EthereumAddress, testTokenAddress ethbridge.EthereumAddress, validatorAddress sdk.ValAddress, testEthereumAddress ethbridge.EthereumAddress, denom, id string, claimType ethbridge.ClaimType) BridgeClaim {
	testCosmosAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err1)
	nftClaim := NewNFTBridgeClaim(TestEthereumChainID, testContractAddress, TestNonce, TestSymbol, testTokenAddress, testEthereumAddress, testCosmosAddress, validatorAddress, denom, id, claimType)
	return nftClaim
}

func CreateTestBurnMsg(t *testing.T, testCosmosSender string, ethereumReceiver ethbridge.EthereumAddress, denom, id string) MsgBurnNFT {
	testTokenContractAddress := ethbridge.NewEthereumAddress(TestTokenContractAddress)
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	burnEth := NewMsgBurnNFT(TestEthereumChainID, testTokenContractAddress, testCosmosAddress, ethereumReceiver, denom, id)
	return burnEth
}

func CreateTestQueryNFTProphecyResponse(cdc *codec.Codec, t *testing.T, validatorAddress sdk.ValAddress, claimType ethbridge.ClaimType) QueryNFTProphecyResponse {
	testEthereumAddress := ethbridge.NewEthereumAddress(TestEthereumAddress)
	testContractAddress := ethbridge.NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := ethbridge.NewEthereumAddress(TestTokenContractAddress)
	nftBridgeClaim := CreateTestNFTClaim(t, testContractAddress, testTokenAddress, validatorAddress, testEthereumAddress, TestDenom, TestID, claimType)
	oracleClaim, _ := CreateOracleClaimFromNFTClaim(cdc, nftBridgeClaim)
	nftBridgeClaims := []BridgeClaim{nftBridgeClaim}

	return NewQueryNFTProphecyResponse(
		oracleClaim.ID,
		oracle.NewStatus(oracle.PendingStatusText, ""),
		nftBridgeClaims,
	)
}
