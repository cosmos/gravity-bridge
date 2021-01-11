package peggy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
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
	delta := float64(types.BridgeValidators(vs.Members).TotalPower()) * 0.01
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
