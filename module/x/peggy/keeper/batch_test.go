package keeper

import (
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatches(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5" // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(99999, myTokenContractAddr).PeggyCoin(),
		)
	)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, input.BankKeeper.SetBalances(ctx, mySender, allVouchers))

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	for i, v := range []uint64{2, 3, 2, 1} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).PeggyCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch, err := input.PeggyKeeper.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// then batch is persisted
	gotFirstBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, firstBatch.TokenContract, firstBatch.BatchNonce)
	require.NotNil(t, gotFirstBatch)

	expFirstBatch := &types.OutgoingTxBatch{
		BatchNonce: 1,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          2,
				Erc20Fee:    types.NewERC20Token(3, myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewERC20Token(101, myTokenContractAddr),
			},
			{
				Id:          1,
				Erc20Fee:    types.NewERC20Token(2, myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewERC20Token(100, myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr,
		Block:         1234567,
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.OutgoingTx
	input.PeggyKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.OutgoingTx{
		{
			BridgeFee: types.NewERC20Token(2, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(102, myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(1, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(103, myTokenContractAddr).PeggyCoin(),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	for i, v := range []uint64{4, 5} {

		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).PeggyCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch, err := input.PeggyKeeper.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.OutgoingTxBatch{
		BatchNonce: 2,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          6,
				Erc20Fee:    types.NewERC20Token(5, myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewERC20Token(101, myTokenContractAddr),
			},
			{
				Id:          5,
				Erc20Fee:    types.NewERC20Token(4, myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewERC20Token(100, myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr,
		Block:         1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	input.PeggyKeeper.OutgoingTxBatchExecuted(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)

	// check batch has been deleted
	gotSecondBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.PeggyKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.OutgoingTx{
		{
			BridgeFee: types.NewERC20Token(3, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(101, myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(2, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(100, myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(2, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(102, myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(1, myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(103, myTokenContractAddr).PeggyCoin(),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}

// tests that batches work with large token amounts, mostly a duplicate of the above
// tests but using much bigger numbers
func TestBatchesFullCoins(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5"   // Pickle
		totalCoins, _       = sdk.NewIntFromString("1500000000000000000000") // 1,500 ETH worth
		oneEth, _           = sdk.NewIntFromString("1000000000000000000")
		allVouchers         = sdk.NewCoins(
			types.NewSDKIntERC20Token(totalCoins, myTokenContractAddr).PeggyCoin(),
		)
	)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, input.BankKeeper.SetBalances(ctx, mySender, allVouchers))

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	for _, v := range []uint64{20, 300, 25, 10} {
		vAsSDKInt := sdk.NewIntFromUint64(v)
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).PeggyCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch, err := input.PeggyKeeper.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// then batch is persisted
	gotFirstBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, firstBatch.TokenContract, firstBatch.BatchNonce)
	require.NotNil(t, gotFirstBatch)

	expFirstBatch := &types.OutgoingTxBatch{
		BatchNonce: 1,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          2,
				Erc20Fee:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr),
			},
			{
				Id:          3,
				Erc20Fee:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr,
		Block:         1234567,
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.OutgoingTx
	input.PeggyKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.OutgoingTx{
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr).PeggyCoin(),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	for _, v := range []uint64{4, 5} {
		vAsSDKInt := sdk.NewIntFromUint64(v)
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).PeggyCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch, err := input.PeggyKeeper.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.OutgoingTxBatch{
		BatchNonce: 2,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          1,
				Erc20Fee:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(20)), myTokenContractAddr),
			},
			{
				Id:          4,
				Erc20Fee:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
				Sender:      mySender.String(),
				DestAddress: myReceiver,
				Erc20Token:  types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(10)), myTokenContractAddr),
			},
		},
		TokenContract: myTokenContractAddr,
		Block:         1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	input.PeggyKeeper.OutgoingTxBatchExecuted(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)

	// check batch has been deleted
	gotSecondBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.PeggyKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.OutgoingTx{
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(300)), myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(25)), myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(5)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(5)), myTokenContractAddr).PeggyCoin(),
		},
		{
			BridgeFee: types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(4)), myTokenContractAddr).PeggyCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewSDKIntERC20Token(oneEth.Mul(sdk.NewIntFromUint64(4)), myTokenContractAddr).PeggyCoin(),
		},
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}
