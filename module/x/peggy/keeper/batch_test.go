package keeper

import (
	"bytes"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatches(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)
	var (
		mySender            = bytes.Repeat([]byte{1}, sdk.AddrLen)
		myReceiver          = types.NewEthereumAddress("eth receiver")
		myTokenContractAddr = types.NewEthereumAddress("my eth oken address")
		myETHToken          = "myETHToken"
		voucherDenom        = types.NewVoucherDenom(myTokenContractAddr, myETHToken)
		now                 = time.Now().UTC()
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
	k.StoreCounterpartDenominator(ctx, myTokenContractAddr, myETHToken)

	// add some TX to the pool
	for i, v := range []int64{2, 3, 2, 1} {
		amount := sdk.NewInt64Coin(string(voucherDenom), int64(i+100))
		fee := sdk.NewInt64Coin(string(voucherDenom), v)
		_, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}
	// when
	ctx = ctx.WithBlockTime(now)
	batchID, err := k.BuildOutgoingTXBatch(ctx, voucherDenom, 2)
	require.NoError(t, err)
	t.Logf("___ response: %#v", batchID)

	// then batch is persisted
	gotBatch := k.GetOutgoingTXBatch(ctx, batchID)
	require.NotNil(t, gotBatch)

	denominator := types.NewBridgedDenominator(myTokenContractAddr, myETHToken)
	expBatch := types.OutgoingTxBatch{
		Nonce: types.NewUInt64Nonce(1),
		Elements: []types.OutgoingTransferTx{
			{
				ID:          2,
				BridgeFee:   denominator.ToUint64ERC20Token(3),
				Sender:      mySender,
				DestAddress: myReceiver,
				Amount:      denominator.ToUint64ERC20Token(101),
			},
			{
				ID:          1,
				BridgeFee:   denominator.ToUint64ERC20Token(2),
				Sender:      mySender,
				DestAddress: myReceiver,
				Amount:      denominator.ToUint64ERC20Token(100),
			},
		},
		CreatedAt:          now,
		TotalFee:           denominator.ToUint64ERC20Token(5),
		BridgedDenominator: denominator,
		BatchStatus:        types.BatchStatusPending,
	}
	assert.Equal(t, expBatch, *gotBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []types.OutgoingTx
	k.IterateOutgoingPoolByFee(ctx, voucherDenom, func(_ uint64, tx types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []types.OutgoingTx{
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
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// ------
	// and when canceled

	err = k.CancelOutgoingTXBatch(ctx, batchID)
	require.NoError(t, err)

	// then
	gotBatch = k.GetOutgoingTXBatch(ctx, batchID)
	require.NotNil(t, gotBatch)
	assert.Equal(t, types.BatchStatusCancelled, gotBatch.BatchStatus)

	// and all TX added back to unbatched pool
	gotUnbatchedTx = nil
	k.IterateOutgoingPoolByFee(ctx, voucherDenom, func(_ uint64, tx types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []types.OutgoingTx{
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
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}
