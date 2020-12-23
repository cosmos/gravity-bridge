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

	// SETUP
	// =====

	k, ctx, keepers := CreateTestEnv(t)
	var (
		// Sender on peggy
		mySender, _ = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		// Random ETH address
		myReceiver = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		// Pickle token contract
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5"
		now                 = time.Now().UTC()
	)
	// mint some voucher first

	allVouchers := sdk.Coins{types.NewERC20Token(99999, myTokenContractAddr).PeggyCoin()}
	err := keepers.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	keepers.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = keepers.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	for i, v := range []uint64{2, 3, 2, 1} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).PeggyCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).PeggyCoin()
		_, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}
	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	firstBatch, err := k.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err)

	// then batch is persisted
	gotFirstBatch := k.GetOutgoingTXBatch(ctx, firstBatch.TokenContract, firstBatch.BatchNonce)
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
		// Valset:        &types.Valset{Nonce: 0x12d687, Members: types.BridgeValidators(nil)},
	}
	assert.Equal(t, expFirstBatch, gotFirstBatch)

	// and verify remaining available Tx in the pool
	var gotUnbatchedTx []*types.OutgoingTx
	k.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
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
		_, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// create the more profitable batch
	ctx = ctx.WithBlockTime(now)
	// tx batch size is 2, so that some of them stay behind
	secondBatch, err := k.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
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
	k.OutgoingTxBatchExecuted(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)

	// check batch has been deleted
	gotSecondBatch := k.GetOutgoingTXBatch(ctx, secondBatch.TokenContract, secondBatch.BatchNonce)
	require.Nil(t, gotSecondBatch)

	// check that txs from first batch have been freed
	gotUnbatchedTx = nil
	k.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
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
