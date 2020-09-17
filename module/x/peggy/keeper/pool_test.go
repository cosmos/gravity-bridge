package keeper

import (
	"bytes"
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToOutgoingPool(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)
	var (
		mySender             = bytes.Repeat([]byte{1}, sdk.AddrLen)
		myReceiver           = "eth receiver"
		myBridgeContractAddr = "my eth bridge contract address"
		myETHToken           = "myETHToken"
		voucherDenom         = types.NewVoucherDenom(myBridgeContractAddr, myETHToken)
	)
	// mint some voucher first
	allVouchers := sdk.Coins{sdk.NewInt64Coin(string(voucherDenom), 99999)}
	err := keepers.SupplyKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	keepers.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = keepers.BankKeeper.SetCoins(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// store counterpart
	k.SetCounterpartDenominator(ctx, myBridgeContractAddr, myETHToken)

	// when
	for i, v := range []int64{2, 3, 2, 1} {
		amount := sdk.NewInt64Coin(string(voucherDenom), int64(i+100))
		fee := sdk.NewInt64Coin(string(voucherDenom), v)
		r, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
		t.Logf("___ response: %#v", r)
	}
	// then
	var got []types.OutgoingTx
	err = k.IterateOutgoingPoolByFee(ctx, voucherDenom, func(_ uint64, tx types.OutgoingTx) bool {
		got = append(got, tx)
		return false
	})
	exp := []types.OutgoingTx{
		{
			BridgeFee:   sdk.NewInt64Coin(string(voucherDenom), 3),
			Sender:      mySender,
			DestAddress: myReceiver,
			Amount:      sdk.NewInt64Coin(string(voucherDenom), 101),
		},
		{
			BridgeFee:   sdk.NewInt64Coin(string(voucherDenom), 2),
			Sender:      mySender,
			DestAddress: myReceiver,
			Amount:      sdk.NewInt64Coin(string(voucherDenom), 100),
		},
		{
			BridgeFee:   sdk.NewInt64Coin(string(voucherDenom), 2),
			Sender:      mySender,
			DestAddress: myReceiver,
			Amount:      sdk.NewInt64Coin(string(voucherDenom), 102),
		},
		{
			BridgeFee:   sdk.NewInt64Coin(string(voucherDenom), 1),
			Sender:      mySender,
			DestAddress: myReceiver,
			Amount:      sdk.NewInt64Coin(string(voucherDenom), 103),
		},
	}
	assert.Equal(t, exp, got)
}
