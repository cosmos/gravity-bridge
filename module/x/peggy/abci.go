package peggy

import (
	"github.com/althea-net/peggy/module/x/peggy/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	// valsets are sorted so the most recent one is first
	for i, vs := range k.GetValsets(ctx) {
		switch {
		// on the latest validator set, check for change in power against
		// current, and emit a new validator set if the change in power >1%
		case i == 0:
		// // currentValset := k.GetCurrentValset(ctx)
		// TODO: how to calculate the diff between two []*BridgeValidator
		// // k.SetValsetRequest(ctx)

		// #1 condition
		// We look through the full bonded validator set (not just the active set, include unbonding validators)
		// and we slash users who haven't signed a valset that is currentHeight - signedBlocksWindow old
		case uint64(ctx.BlockHeight())-params.SignedBlocksWindow > vs.Nonce:
			// first we need to see which validators in the active set
			// haven't signed the valdiator set and slash them,
			var toSlash []stakingtypes.Validator
			for _, val := range currentBondedSet {
				found := false
				for _, conf := range k.GetValsetConfirms(ctx, vs.Nonce) {
					if conf.EthAddress == k.GetEthAddress(ctx, sdk.AccAddress(val.GetOperator())) {
						found = true
						break
					}
				}
				if !found {
					toSlash = append(toSlash, val)
				}
			}

			for _, val := range toSlash {
				cons, _ := val.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionValset)
			}

			// then we prune the valset from state
			k.DeleteValset(ctx, vs.Nonce)
		}

	}

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
