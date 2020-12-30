package peggy

import (
	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	// valsets are sorted so the most recent one is first
	valsets := k.GetValsets(ctx)
	if len(valsets) == 0 {
		k.SetValsetRequest(ctx)
	}

	for i, vs := range valsets {
		switch {
		// #1 condition
		// We look through the full bonded validator set (not just the active set, include unbonding validators)
		// and we slash users who haven't signed a valset that is currentHeight - signedBlocksWindow old
		case uint64(ctx.BlockHeight())-params.SignedValsetsWindow > vs.Height:

			// first we need to see which validators in the active set
			// haven't signed the valdiator set and slash them,
			var toSlash []stakingtypes.Validator
			confirms := k.GetValsetConfirms(ctx, vs.Nonce)
			for _, val := range currentBondedSet {
				found := false
				for _, conf := range confirms {
					if conf.EthAddress == k.GetEthAddress(ctx, val.GetOperator()) {
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
				k.StakingKeeper.Jail(ctx, cons)
			}

			// then we prune the valset from state
			k.DeleteValset(ctx, vs.Nonce)

		// on the latest validator set, check for change in power against
		// current, and emit a new validator set if the change in power >1%
		case i == 0:
			if types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(vs.Members) > 0.01 {
				k.SetValsetRequest(ctx)
			}
		}
	}

	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		if uint64(ctx.BlockHeight())-params.SignedBatchesWindow > batch.Block {
			var toSlash []stakingtypes.Validator
			confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
			for _, val := range currentBondedSet {
				found := false
				for _, conf := range confirms {
					// TODO: may need to look up actual validator address
					confVal, _ := sdk.AccAddressFromBech32(conf.Orchestrator)
					if confVal.Equals(val.GetOperator()) {
						found = true
					}
				}
				if !found {
					toSlash = append(toSlash, val)
				}
			}
			for _, val := range toSlash {
				cons, _ := val.GetConsAddr()
				// TODO: make this a different slash fraction in the params
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
				k.StakingKeeper.Jail(ctx, cons)
			}

			// clean up batches here
			k.DeleteBatch(ctx, *batch)
		}
	}

	// #3 condition
	// Oracle events MsgDepositClaim, MsgWithdrawClaim
	attmap := k.GetAttestationMapping(ctx)
	for _, atts := range attmap {
		// Conflicting votes should be slashed
		if len(atts) > 1 {
			var (
				toSlash []string
				unObs   []types.Attestation
			)
			oneObserved := false
			for _, att := range atts {
				if att.Observed == true {
					oneObserved = true
					continue
				}
				unObs = append(unObs, att)
			}
			if oneObserved {
				for _, att := range unObs {
					toSlash = append(toSlash, att.Votes...)
					k.DeleteAttestation(ctx, att)
				}
			}
			for _, valaddr := range toSlash {
				validator, _ := sdk.ValAddressFromBech32(valaddr)
				val := k.StakingKeeper.Validator(ctx, validator)
				cons, _ := val.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, validator), params.SlashFractionConflictingClaim)
				k.StakingKeeper.Jail(ctx, cons)
			}
		}

		// Pair tomorrow with Justin on slashing for not voting
		// TODO: time out attestations
	}
	// Blocked on storing of the claim

	// #4 condition (stretch goal)
	// TODO: lost eth key or delegate key
	// 1. submit a message signed by the priv key to the chain and it slashes the validator who delegated to that key
	// return

	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}
