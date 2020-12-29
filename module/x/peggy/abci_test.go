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
	nonce := uint64(ctx.BlockHeight()) - (params.SignedBlocksWindow + 1)
	vs.Nonce = nonce
	pk.StoreValset(ctx, vs)
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
	valset := input.PeggyKeeper.GetValset(ctx, nonce)
	require.Nil(t, valset)

	// TODO: test balance of slashed tokens
}

func TestValsetEmission(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.PeggyKeeper

	// Store a validator set with a power change as the most recent validator set
	vs := pk.GetCurrentValset(ctx)
	vs.Nonce = vs.Nonce - 1
	delta := float64(types.BridgeValidators(vs.Members).TotalPower()) * 0.011
	vs.Members[0].Power = uint64(float64(vs.Members[0].Power) - delta)
	pk.StoreValset(ctx, vs)

	// EndBlocker should set a new validator set
	EndBlocker(ctx, pk)
	require.NotNil(t, pk.GetValset(ctx, uint64(ctx.BlockHeight())))
}
