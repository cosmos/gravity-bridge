package keeper

import (
	"math/big"
	"testing"

	"github.com/althea-net/peggy/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToOutgoingPool(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	var (
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5"
	)
	// mint some voucher first
	allVouchers := sdk.Coins{types.NewERC20Token(99999, myTokenContractAddr).GravityCoin()}
	err := input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = input.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// when
	for i, v := range []uint64{2, 3, 2, 1} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).GravityCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).GravityCoin()
		r, err := input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
		t.Logf("___ response: %#v", r)
	}
	// then
	var got []*types.OutgoingTx
	input.GravityKeeper.IterateOutgoingPoolByFee(ctx, myTokenContractAddr, func(_ uint64, tx *types.OutgoingTx) bool {
		got = append(got, tx)
		return false
	})
	exp := []*types.OutgoingTx{
		{
			BridgeFee: types.NewERC20Token(3, myTokenContractAddr).GravityCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(101, myTokenContractAddr).GravityCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(2, myTokenContractAddr).GravityCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(100, myTokenContractAddr).GravityCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(2, myTokenContractAddr).GravityCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(102, myTokenContractAddr).GravityCoin(),
		},
		{
			BridgeFee: types.NewERC20Token(1, myTokenContractAddr).GravityCoin(),
			Sender:    mySender.String(),
			DestAddr:  myReceiver,
			Amount:    types.NewERC20Token(103, myTokenContractAddr).GravityCoin(),
		},
	}
	assert.Equal(t, exp, got)
}

func TestTotalBatchFeeInPool(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	// token1
	var (
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5"
	)
	// mint some voucher first
	allVouchers := sdk.Coins{types.NewERC20Token(99999, myTokenContractAddr).GravityCoin()}
	err := input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = input.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// create outgoing pool
	for i, v := range []uint64{2, 3, 2, 1} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).GravityCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).GravityCoin()
		r, err := input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
		t.Logf("___ response: %#v", r)
	}

	// token 2 - Only top 100
	var (
		myToken2ContractAddr = "0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0"
	)
	// mint some voucher first
	allVouchers = sdk.Coins{types.NewERC20Token(18446744073709551615, myToken2ContractAddr).GravityCoin()}
	err = input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = input.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// Add

	// create outgoing pool
	for i := 0; i < 110; i++ {
		amount := types.NewERC20Token(uint64(i+100), myToken2ContractAddr).GravityCoin()
		fee := types.NewERC20Token(uint64(5), myToken2ContractAddr).GravityCoin()
		r, err := input.GravityKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
		t.Logf("___ response: %#v", r)
	}

	batchFees := input.GravityKeeper.CreateBatchFees(ctx)
	/*
		tokenFeeMap should be
		map[0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5:8 0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0:500]
		**/
	assert.Equal(t, batchFees[0].TopOneHundred.BigInt(), big.NewInt(int64(8)))
	assert.Equal(t, batchFees[1].TopOneHundred.BigInt(), big.NewInt(int64(500)))

}
