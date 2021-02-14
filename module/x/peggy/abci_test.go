package peggy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValsetSlashing(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper
	params := input.PeggyKeeper.GetParams(ctx)

	// This valset should be past the signed blocks window and trigger slashing
	vs := pk.GetCurrentValset(ctx)
	height := uint64(ctx.BlockHeight()) - (params.SignedValsetsWindow + 1)
	vs.Height = height

	// TODO: remove this once we are auto-incrementing the nonces
	vs.Nonce = height
	pk.StoreValsetUnsafe(ctx, vs)
	for i, val := range keeper.AccAddrs {
		if i == 0 {
			// don't sign with first validator
			continue
		}
		conf := types.NewMsgValsetConfirm(vs.Nonce, keeper.EthAddrs[i].String(), val, "dummysig")
		pk.SetValsetConfirm(ctx, *conf)
	}

	// Set the current valset to avoid setting a new valset in the switch
	pk.SetValsetRequest(ctx)
	EndBlocker(ctx, pk)

	// ensure that the  validator is jailed and slashed
	val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	require.True(t, val.IsJailed())

	// Ensure that the valset gets pruned properly
	valset := input.PeggyKeeper.GetValset(ctx, vs.Nonce)
	require.Nil(t, valset)

	// TODO: test balance of slashed tokens
}

func TestBatchSlashing(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper
	params := pk.GetParams(ctx)

	// First store a batch
	batch := &types.OutgoingTxBatch{
		BatchNonce:    1,
		Transactions:  []*types.OutgoingTransferTx{},
		TokenContract: keeper.TokenContractAddrs[0],
		Block:         uint64(ctx.BlockHeight() - int64(params.SignedBatchesWindow+1)),
	}
	pk.StoreBatchUnsafe(ctx, batch)

	for i, val := range keeper.AccAddrs {
		if i == 0 {
			// don't sign with first validator
			continue
		}
		pk.SetBatchConfirm(ctx, &types.MsgConfirmBatch{
			Nonce:         batch.BatchNonce,
			TokenContract: keeper.TokenContractAddrs[0],
			EthSigner:     keeper.EthAddrs[i].String(),
			Orchestrator:  val.String(),
		})
	}

	EndBlocker(ctx, pk)

	// ensure that the  validator is jailed and slashed
	val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	require.True(t, val.IsJailed())

	// Ensure that the batch gets pruned properly
	batch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, batch.TokenContract, batch.BatchNonce)
	require.Nil(t, batch)
}

func TestValsetEmission(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper

	// Store a validator set with a power change as the most recent validator set
	vs := pk.GetCurrentValset(ctx)
	vs.Nonce = vs.Nonce - 1
	delta := float64(types.BridgeValidators(vs.Members).TotalPower()) * 0.05
	vs.Members[0].Power = uint64(float64(vs.Members[0].Power) - delta/2)
	vs.Members[1].Power = uint64(float64(vs.Members[1].Power) + delta/2)
	pk.StoreValset(ctx, vs)

	// EndBlocker should set a new validator set
	EndBlocker(ctx, pk)
	require.NotNil(t, pk.GetValset(ctx, uint64(ctx.BlockHeight())))
	valsets := pk.GetValsets(ctx)
	require.True(t, len(valsets) == 2)
}

func TestValsetSetting(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper
	pk.SetValsetRequest(ctx)
	valsets := pk.GetValsets(ctx)
	require.True(t, len(valsets) == 1)
}

/// Test batch timeout
func TestBatchTimeout(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper
	params := pk.GetParams(ctx)
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		myTokenContractAddr = "0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5" // Pickle
		allVouchers         = sdk.NewCoins(
			types.NewERC20Token(99999, myTokenContractAddr).PeggyCoin(),
		)
	)

	require.Greater(t, params.AverageBlockTime, uint64(0))
	require.Greater(t, params.AverageEthereumBlockTime, uint64(0))

	// mint some vouchers first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, input.BankKeeper.SetBalances(ctx, mySender, allVouchers))

	// add some TX to the pool
	for i, v := range []uint64{2, 3, 2, 1, 5, 6} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).PeggyCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}

	// when
	ctx = ctx.WithBlockTime(now)
	ctx = ctx.WithBlockHeight(250)

	// check that we can make a batch without first setting an ethereum block height
	b1, err1 := pk.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err1)
	require.Equal(t, b1.BatchTimeout, uint64(0))

	pk.SetLastObservedEthereumBlockHeight(ctx, 500)

	b2, err2 := pk.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err2)
	// this is exactly block 500 plus twelve hours
	require.Equal(t, b2.BatchTimeout, uint64(504))

	// make sure the batches got stored in the first place
	gotFirstBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, b1.TokenContract, b1.BatchNonce)
	require.NotNil(t, gotFirstBatch)
	gotSecondBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, b2.TokenContract, b2.BatchNonce)
	require.NotNil(t, gotSecondBatch)

	// when, way into the future
	ctx = ctx.WithBlockTime(now)
	ctx = ctx.WithBlockHeight(9)

	b3, err2 := pk.BuildOutgoingTXBatch(ctx, myTokenContractAddr, 2)
	require.NoError(t, err2)

	EndBlocker(ctx, pk)

	// this had a timeout of zero should be deleted.
	gotFirstBatch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, b1.TokenContract, b1.BatchNonce)
	require.Nil(t, gotFirstBatch)
	// make sure the end blocker does not delete these, as the block height has not officially
	// been updated by a relay event
	gotSecondBatch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, b2.TokenContract, b2.BatchNonce)
	require.NotNil(t, gotSecondBatch)
	gotThirdBatch := input.PeggyKeeper.GetOutgoingTXBatch(ctx, b3.TokenContract, b3.BatchNonce)
	require.NotNil(t, gotThirdBatch)

	pk.SetLastObservedEthereumBlockHeight(ctx, 5000)
	EndBlocker(ctx, pk)

	// make sure the end blocker does delete these, as we've got a new Ethereum block height
	gotFirstBatch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, b1.TokenContract, b1.BatchNonce)
	require.Nil(t, gotFirstBatch)
	gotSecondBatch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, b2.TokenContract, b2.BatchNonce)
	require.Nil(t, gotSecondBatch)
	gotThirdBatch = input.PeggyKeeper.GetOutgoingTXBatch(ctx, b3.TokenContract, b3.BatchNonce)
	require.NotNil(t, gotThirdBatch)

}
