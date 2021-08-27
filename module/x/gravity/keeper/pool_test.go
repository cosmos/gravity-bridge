package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestAddToOutgoingPool(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")
	)
	// mint some voucher first
	allVouchers := sdk.Coins{types.NewERC20Token(99999, myTokenContractAddr.Hex()).GravityCoin()}
	err := input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, fundAccount(ctx, input.BankKeeper, mySender, allVouchers))

	// when
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1)

	// then
	var got []*types.SendToEthereum
	input.GravityKeeper.IterateUnbatchedSendToEthereums(ctx, func(tx *types.SendToEthereum) bool {
		got = append(got, tx)
		return false
	})

	exp := []*types.SendToEthereum{
		types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
		types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
		types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 100, 2),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 103, 1),
	}

	require.Equal(t, exp, got)
	require.EqualValues(t, exp[0], got[0])
	require.EqualValues(t, exp[1], got[1])
	require.EqualValues(t, exp[2], got[2])
	require.EqualValues(t, exp[3], got[3])
	require.Len(t, got, 4)
}
