package peggy

import (
	"github.com/althea-net/peggy/module/x/peggy/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// params := k.GetParams(ctx)

	// downtimeBlock := ctx.BlockHeight() - int64(params.SignedBlocksWindow)

	// #1 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a valset that is >15hrs in blocks old
	// k.IterateValsets()
	// find last valset, k.GetCurrentValset.Diff(latest valset) if > %1 then we k.SetValsetRequest
	// prune old valsets
	// // k.IterateValsetConfirmByNonce()
	// // Slash here

	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
	// k.IterateOutgoingTXBatches()
	// if there are batches older than 15h that are confirmed, prune them from state
	// // k.IterateBatchConfirmByNonceAndTokenContract()

	// #3 condition
	// Oracle events MsgDepositClaim, MsgWithdrawClaim
	// Blocked on storing of the claim

	// #4 condition (stretch goal)
	// TODO: lost eth key or delegate key
	// 1. submit a message signed by the priv key to the chain and it slashes the validator who delegated to that key
	// return

	// stretch goal
	// Trigger valset creation

	// TODO: prune valsets older than one month
	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}
