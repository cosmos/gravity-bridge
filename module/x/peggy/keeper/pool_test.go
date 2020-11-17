package keeper

import (
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToOutgoingPool(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)
	var (
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = types.NewEthereumAddress("eth receiver")
		myETHToken          = "myETHToken"
		myTokenContractAddr = types.NewEthereumAddress("my eth oken address")
		voucherDenom        = types.NewVoucherDenom(myTokenContractAddr, myETHToken)
	)
	// mint some voucher first
	allVouchers := sdk.Coins{sdk.NewInt64Coin(string(voucherDenom), 99999)}
	err := keepers.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	keepers.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = keepers.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// store counterpart
	k.StoreCounterpartDenominator(ctx, myTokenContractAddr, myETHToken)

	// when
	for i, v := range []int64{2, 3, 2, 1} {
		amount := sdk.NewInt64Coin(string(voucherDenom), int64(i+100))
		fee := sdk.NewInt64Coin(string(voucherDenom), v)
		r, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
		t.Logf("___ response: %#v", r)
	}
	// then
	var got []*types.OutgoingTx
	k.IterateOutgoingPoolByFee(ctx, voucherDenom, func(_ uint64, tx *types.OutgoingTx) bool {
		got = append(got, tx)
		return false
	})
	exp := []*types.OutgoingTx{
		{
			BridgeFee: sdk.NewInt64Coin(string(voucherDenom), 3),
			Sender:    mySender.String(),
			DestAddr:  myReceiver.Bytes(),
			Amount:    sdk.NewInt64Coin(string(voucherDenom), 101),
		},
		{
			BridgeFee: sdk.NewInt64Coin(string(voucherDenom), 2),
			Sender:    mySender.String(),
			DestAddr:  myReceiver.Bytes(),
			Amount:    sdk.NewInt64Coin(string(voucherDenom), 100),
		},
		{
			BridgeFee: sdk.NewInt64Coin(string(voucherDenom), 2),
			Sender:    mySender.String(),
			DestAddr:  myReceiver.Bytes(),
			Amount:    sdk.NewInt64Coin(string(voucherDenom), 102),
		},
		{
			BridgeFee: sdk.NewInt64Coin(string(voucherDenom), 1),
			Sender:    mySender.String(),
			DestAddr:  myReceiver.Bytes(),
			Amount:    sdk.NewInt64Coin(string(voucherDenom), 103),
		},
	}
	assert.Equal(t, exp, got)
}
