package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestBatches(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(99999, myTokenContractAddr.Hex()).GravityCoin(),
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
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1)

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch, err := input.GravityKeeper.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// then batch is persisted
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, firstBatch.GetStoreIndex())
	require.NotNil(t, gotFirstBatch)

	expFirstBatch := &types.BatchTx{
		Nonce: 1,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
			types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 100, 2),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.SendToEthereum
	input.GravityKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.SendToEthereum{
		types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 103, 1),
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 4, 5)

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch, err := input.GravityKeeper.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.BatchTx{
		Nonce: 2,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(6, myTokenContractAddr, mySender, myReceiver, 101, 5),
			types.NewSendToEthereumTx(5, myTokenContractAddr, mySender, myReceiver, 100, 4),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	err = input.GravityKeeper.BatchTxExecuted(ctx, common.HexToAddress(secondBatch.TokenContract), secondBatch.Nonce)
	require.NoError(t, err)

	// check batch has been deleted
	gotSecondBatch := input.GravityKeeper.GetOutgoingTx(ctx, secondBatch.GetStoreIndex())
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.GravityKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.SendToEthereum{
		types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
		types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 100, 2),
		types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 103, 1),
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
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		totalCoins, _       = sdk.NewIntFromString("1500000000000000000000")                    // 1,500 ETH worth
		oneEth, _           = sdk.NewIntFromString("1000000000000000000")
		allVouchers         = sdk.NewCoins(
			types.NewSDKIntERC20Token(totalCoins, myTokenContractAddr).GravityCoin(),
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
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		_, err := input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver.Hex(), amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch, err := input.GravityKeeper.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// then batch is persisted
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, firstBatch.GetStoreIndex())
	require.NotNil(t, gotFirstBatch)

	expFirstBatch := &types.BatchTx{
		Nonce: 1,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 300, 300),
			types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 25, 25),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.SendToEthereum
	input.GravityKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx := []*types.SendToEthereum{
		types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 20, 20),
		types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 10, 10),
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)

	// CREATE SECOND, MORE PROFITABLE BATCH
	// ====================================

	// add some more TX to the pool to create a more profitable batch
	for _, v := range []uint64{4, 5} {
		vAsSDKInt := sdk.NewIntFromUint64(v)
		amount := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		fee := types.NewSDKIntERC20Token(oneEth.Mul(vAsSDKInt), myTokenContractAddr).GravityCoin()
		_, err = input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver.Hex(), amount, fee)
		require.NoError(t, err)
	}

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch, err := input.GravityKeeper.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// check that the more profitable batch has the right txs in it
	expSecondBatch := &types.BatchTx{
		Nonce: 2,
		Transactions: []*types.SendToEthereum{
			types.NewSendToEthereumTx(1, myTokenContractAddr, mySender, myReceiver, 20, 20),
			types.NewSendToEthereumTx(4, myTokenContractAddr, mySender, myReceiver, 10, 10),
		},
		TokenContract: myTokenContractAddr.Hex(),
		Height:        1234567,
	}

	assert.Equal(t, expSecondBatch, secondBatch)

	// EXECUTE THE MORE PROFITABLE BATCH
	// =================================

	// Execute the batch
	err = input.GravityKeeper.BatchTxExecuted(ctx, common.HexToAddress(secondBatch.TokenContract), secondBatch.Nonce)
	require.NoError(t, err)

	// check batch has been deleted
	gotSecondBatch := input.GravityKeeper.GetOutgoingTx(ctx, secondBatch.GetStoreIndex())
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	input.GravityKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.SendToEthereum) bool {
		gotUnbatchedTx = append(gotUnbatchedTx, tx)
		return false
	})
	expUnbatchedTx = []*types.SendToEthereum{
		types.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 300, 300),
		types.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 25, 25),
		types.NewSendToEthereumTx(6, myTokenContractAddr, mySender, myReceiver, 5, 5),
		types.NewSendToEthereumTx(5, myTokenContractAddr, mySender, myReceiver, 4, 4),
	}
	assert.Equal(t, expUnbatchedTx, gotUnbatchedTx)
}

func TestPoolTxRefund(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5" // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(414, myTokenContractAddr).GravityCoin(),
		)
		myDenom = types.NewERC20Token(1, myTokenContractAddr).GravityCoin().Denom
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
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).GravityCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).GravityCoin()
		_, err := input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	_, err := input.GravityKeeper.BuildBatchTx(ctx, common.HexToAddress(myTokenContractAddr), 2)
	require.NoError(t, err)

	// try to refund a tx that's in a batch
	err1 := input.GravityKeeper.RemoveFromOutgoingPoolAndRefund(ctx, 1, mySender)
	require.Error(t, err1)

	// try to refund a tx that's in the pool
	err2 := input.GravityKeeper.RemoveFromOutgoingPoolAndRefund(ctx, 4, mySender)
	require.NoError(t, err2)

	// make sure refund was issued
	balances := input.BankKeeper.GetAllBalances(ctx, mySender)
	require.Equal(t, sdk.NewInt(104), balances.AmountOf(myDenom))
}
