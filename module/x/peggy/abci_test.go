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
	vs.Nonce = uint64(ctx.BlockHeight()) - (params.SignedBlocksWindow + 1)
	pk.StoreValset(ctx, vs)
	for i, val := range keeper.AccAddrs {
		if i == 0 {
			// don't sign the first validator
			continue
		}
		conf := types.NewMsgValsetConfirm(vs.Nonce, keeper.EthAddrs[i].String(), val, "dummysig")
		pk.SetValsetConfirm(ctx, *conf)
	}
	// Set the current valset to avoid the first switch
	pk.SetValsetRequest(ctx)
	EndBlocker(ctx, pk)

	val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	require.True(t, val.IsJailed())

	// persist some validator sets
	// persist some validator set confirmations for each validator
	// find test cases for this
}
