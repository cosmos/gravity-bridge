package peggy

import (
	"github.com/althea-net/peggy/module/x/peggy/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// params := k.GetParams(ctx)

	// #1 condition
	// TODO: slashing for lack of valset update signature submission
	// 1. Get params and calculate (S = ctx.current_block - params.num_blocks_downtime)
	// 2. Get the last Validator set request (ethereum signed validator set) with blockheight S< this is V
	// 3. Check if any active validator at blockheight S or still bonded but not active validator at height S has not signed V
	// 4. Slash if true

	// *GOAL OF VALSET SLASHING* unbonding validators need to sign 1 update that does not include their key

	// #2 condition
	// TODO: slashing for lack of batch signature submission
	// 1. Get params and calculate (S = ctx.current_block - params.num_blocks_downtime)
	// 2. Get the last batch set request (etherum signed batch set) with blockheight S< this is V
	// 3. Check if any active validator at blockheight S or still bonded but not active valdiator at height S has not signed V
	// 4. Slash if true

	// *GOAL OF BATCH SLASHING* ensure that batch can be submitted, and if this doesn't occur, this doesn't
	// trigger a correctness violation so long as the validator set gets updated

	// #3 condition
	// TODO: oracle downtime slashing
	// 1. Get params and calculate (S = ctx.current_block - params.num_blocks_downtime)

	// Stretch Goal
	// TODO: lost eth key or delegate key
	// 1. submit a message signed by the priv key to the chain and it slashes the validator who delegated to that key
	return
}
