package peggy

import (
	"sort"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	slashing(ctx, k)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
}

func slashing(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	// valsets are sorted so the most recent one is first
	valsets := k.GetValsets(ctx)
	if len(valsets) == 0 {
		k.SetValsetRequest(ctx)
	}

	for i, vs := range valsets {
		signedWithinWindow := uint64(ctx.BlockHeight()) > params.SignedValsetsWindow && uint64(ctx.BlockHeight())-params.SignedValsetsWindow > vs.Height
		switch {
		// #1 condition
		// We look through the full bonded validator set (not just the active set, include unbonding validators)
		// and we slash users who haven't signed a valset that is currentHeight - signedBlocksWindow old
		case signedWithinWindow:

			// first we need to see which validators in the active set
			// haven't signed the valdiator set and slash them,
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
					cons, _ := val.GetConsAddr()
					k.StakingKeeper.Slash(ctx, cons,
						ctx.BlockHeight(), val.ConsensusPower(),
						params.SlashFractionValset)
					k.StakingKeeper.Jail(ctx, cons)
				}
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
		signedWithinWindow := uint64(ctx.BlockHeight()) > params.SignedBatchesWindow && uint64(ctx.BlockHeight())-params.SignedBatchesWindow > batch.Block
		if signedWithinWindow {
			confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
			for _, val := range currentBondedSet {
				found := false
				for _, conf := range confirms {
					// TODO: double check this logic
					confVal, _ := sdk.AccAddressFromBech32(conf.Orchestrator)
					if k.GetOrchestratorValidator(ctx, confVal).Equals(val.GetOperator()) {
						found = true
						break
					}
				}
				if !found {
					cons, _ := val.GetConsAddr()
					k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
					k.StakingKeeper.Jail(ctx, cons)
				}
			}

			// clean up batches here
			k.DeleteBatch(ctx, *batch)
		}
	}

	// #3 condition
	// Oracle events MsgDepositClaim, MsgWithdrawClaim
	attmap := k.GetAttestationMapping(ctx)
	for _, atts := range attmap {
		// slash conflicting votes
		if len(atts) > 1 {
			var unObs []types.Attestation
			oneObserved := false
			for _, att := range atts {
				if att.Observed {
					oneObserved = true
					continue
				}
				unObs = append(unObs, att)
			}
			// if one is observed delete the *other* attestations, do not delete
			// the original as we will need it later.
			if oneObserved {
				for _, att := range unObs {
					for _, valaddr := range att.Votes {
						validator, _ := sdk.ValAddressFromBech32(valaddr)
						val := k.StakingKeeper.Validator(ctx, validator)
						cons, _ := val.GetConsAddr()
						k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, validator), params.SlashFractionConflictingClaim)
						k.StakingKeeper.Jail(ctx, cons)
					}
					claim, err := k.UnpackAttestationClaim(&att)
					if err != nil {
						panic("couldn't cast to claim")
					}

					k.DeleteAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
				}
			}
		}

		if len(atts) == 1 {
			att := atts[0]
			// TODO-JT: Review this
			windowPassed := uint64(ctx.BlockHeight()) > params.SignedClaimsWindow &&
				uint64(ctx.BlockHeight())-params.SignedClaimsWindow > att.Height

			// if the signing window has passed and the attestation is still unobserved wait.
			if windowPassed && att.Observed {
				for _, bv := range currentBondedSet {
					found := false
					for _, val := range att.Votes {
						confVal, _ := sdk.ValAddressFromBech32(val)
						if confVal.Equals(bv.GetOperator()) {
							found = true
							break
						}
					}
					if !found {
						cons, _ := bv.GetConsAddr()
						k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, bv.GetOperator()), params.SlashFractionClaim)
						k.StakingKeeper.Jail(ctx, cons)
					}
				}
				claim, err := k.UnpackAttestationClaim(&att)
				if err != nil {
					panic("couldn't cast to claim")
				}

				k.DeleteAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
			}
		}
	}

	// #4 condition (stretch goal)
	// TODO: lost eth key or delegate key
	// 1. submit a message signed by the priv key to the chain and it slashes the validator who delegated to that key
	// return

	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k keeper.Keeper) {
	attmap := k.GetAttestationMapping(ctx)
	keys := make([]uint64, 0, len(attmap))
	for k := range attmap {
		keys = append(keys, k)
	}
	// the go standard library does not include a way to sort by uint64s...
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		oneObserved := false
		skipped := false
		for _, att := range attmap[key] {
			// If it was already observed we don't need to do anything
			if att.Observed {
				skipped = true
				continue
			}
			k.TryAttestation(ctx, &att)

			// if we observe any attestation for this nonce no others can pass
			// so it is safe to continue
			if att.Observed {
				oneObserved = true
				continue
			}
		}
		// If any one attestation for a nonce was not skipped and also does
		// not become observed no later events can possibly pass. Meaning it is
		// safe to return early.
		if !oneObserved && !skipped {
			return
		}
	}
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatches(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumBlockHeight
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		if batch.BatchTimeout < ethereumHeight {
			k.CancelOutgoingTXBatch(ctx, batch.TokenContract, batch.BatchNonce)
		}
	}
}
