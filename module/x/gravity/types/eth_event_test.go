package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDepositEvent_Validate(t *testing.T) {
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")
	ethAddr2, err := newEthAddress()
	require.NoError(t, err, "unable to generate second ethereum address")

	testCases := []struct {
		name     string
		event    DepositEvent
		expError bool
	}{
		{
			"default input",
			DepositEvent{
				1,
				ethAddr.String(),
				sdk.NewInt(10),
				ethAddr2.String(),
				orchAddr.String(),
				30,
			},
			false,
		},
		{
			"zero nonce",
			DepositEvent{
				0,
				ethAddr.String(),
				sdk.NewInt(10),
				ethAddr2.String(),
				orchAddr.String(),
				30,
			},
			true,
		},
		{
			"invalid contract address",
			DepositEvent{
				1,
				"not an addr",
				sdk.NewInt(10),
				ethAddr2.String(),
				orchAddr.String(),
				30,
			},
			true,
		},
		{
			"negative amount",
			DepositEvent{
				1,
				ethAddr.String(),
				sdk.NewInt(-10),
				ethAddr2.String(),
				orchAddr.String(),
				30,
			},
			true,
		},
		{
			"invalid sender addr",
			DepositEvent{
				1,
				ethAddr.String(),
				sdk.NewInt(10),
				"not an addr",
				orchAddr.String(),
				30,
			},
			true,
		},
		{
			"invalid receiver addr",
			DepositEvent{
				1,
				ethAddr.String(),
				sdk.NewInt(10),
				ethAddr2.String(),
				"not an addr",
				30,
			},
			true,
		},
		{
			"zero ethereum height",
			DepositEvent{
				1,
				ethAddr.String(),
				sdk.NewInt(10),
				ethAddr2.String(),
				orchAddr.String(),
				0,
			},
			true,
		},
	}
	for _, tc := range testCases {
		err := tc.event.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestWithdrawEvent_Validate(t *testing.T) {
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")

	testCases := []struct {
		name     string
		event    WithdrawEvent
		expError bool
	}{
		{"valid input", WithdrawEvent{[]byte("txid"), 10, ethAddr.String(), 10}, false},
		{"zero height", WithdrawEvent{[]byte("txid"), 10, ethAddr.String(), 0}, true},
		{"invalid contract address", WithdrawEvent{[]byte("txid"), 10, "not an addr", 10}, true},
		{"zero nonce", WithdrawEvent{[]byte("txid"), 0, ethAddr.String(), 10}, true},
		{"empty transaction ID", WithdrawEvent{[]byte(""), 10, ethAddr.String(), 10}, true},
	}
	for _, tc := range testCases {
		err := tc.event.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestCosmosERC20DeployedEvent_Validate(t *testing.T) {
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")

	testCases := []struct {
		name     string
		event    CosmosERC20DeployedEvent
		expError bool
	}{
		{"valid input", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			ethAddr.String(),
			"testname",
			"TT",
			3,
			10,
		}, false},
		{"zero nonce", CosmosERC20DeployedEvent{
			0,
			"testtoken",
			ethAddr.String(),
			"testname",
			"TT",
			3,
			10,
		}, true},
		{"broken denom", CosmosERC20DeployedEvent{
			10,
			"gravity/&$@",
			ethAddr.String(),
			"testname",
			"TT",
			3,
			10,
		}, true},
		{"invalid contract address", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			"not an addr",
			"testname",
			"TT",
			3,
			10,
		}, true},
		{"blank token name", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			ethAddr.String(),
			"",
			"TT",
			3,
			10,
		}, true},
		{"blank token symbol", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			ethAddr.String(),
			"testname",
			"",
			3,
			10,
		}, true},
		{"zero degree precision", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			ethAddr.String(),
			"testname",
			"TT",
			0,
			10,
		}, true},
		{"zero height", CosmosERC20DeployedEvent{
			10,
			"testtoken",
			ethAddr.String(),
			"testname",
			"TT",
			3,
			0,
		}, true},
	}
	for _, tc := range testCases {
		err := tc.event.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestLogicCallExecutedEvent_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		event    LogicCallExecutedEvent
		expError bool
	}{
		{"valid input", LogicCallExecutedEvent{10, []byte("invalidationID"), 20, 30}, false},
		{"zero nonce", LogicCallExecutedEvent{0, []byte("invalidationID"), 20, 30}, true},
		{"empty invalidation id", LogicCallExecutedEvent{10, []byte(""), 20, 30}, true},
		{"zero invalidation nonce", LogicCallExecutedEvent{10, []byte("invalidationID"), 0, 30}, true},
		{"zero height", LogicCallExecutedEvent{10, []byte("invalidationID"), 20, 0}, true},
	}
	for _, tc := range testCases {
		err := tc.event.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}
