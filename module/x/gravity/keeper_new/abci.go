package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Question: what here can be epoched?
	k.slash(ctx)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	cleanupTimedOutLogicCalls(ctx, k)
	// createValsets(ctx, k)
}

// func createValsets(ctx sdk.Context, k Keeper) {
// 	// Auto ValsetRequest Creation.
// 	/*
// 			1. If there are no valset requests, create a new one.
// 			2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
// 			   This will make sure the unbonding validator has to provide an attestation to a new Valset
// 		       that excludes him before he completely Unbonds.  Otherwise he will be slashed
// 			3. If power change between validators of CurrentValset and latest valset request is > 5%
// 		**/
// 	latestValset := k.GetLatestValset(ctx)
// 	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)

// 	if (latestValset == nil) || (lastUnbondingHeight == uint64(ctx.BlockHeight())) || (types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(latestValset.Members) > 0.05) {
// 		k.SetValsetRequest(ctx)
// 	}
// }

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k Keeper) {
	// We check the attestations that haven't been observed, i.e nonce is exactly 1 higher than the last attestation
	nonce := uint64(k.GetLastObservedEventNonce(ctx)) + 1

	k.IterateAttestationByNonce(ctx, nonce, func(attestation types.Attestation) bool {
		// try unobserved attestations
		// TODO: rename. "Try" is too ambiguous
		k.TryAttestation(ctx, attestation)
		return false
	})
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
func cleanupTimedOutBatches(ctx sdk.Context, k Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx)
	if ethereumHeight == 0 {
		panic("ethereum observed height cannot be 0")
	}

	// TODO: use iterator
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		if batch.Timeout < ethereumHeight {
			k.CancelOutgoingTXBatch(ctx, batch.TokenContract, batch.BatchNonce)
		}
	}
}

// cleanupTimedOutBatches deletes logic calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutLogicCalls(ctx sdk.Context, k Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx)
	if ethereumHeight == 0 {
		panic("ethereum observed height cannot be 0")
	}

	// TODO: use iterator
	calls := k.GetOutgoingLogicCalls(ctx)
	for _, call := range calls {
		if call.Timeout < ethereumHeight {
			k.CancelOutgoingLogicCall(ctx, call.InvalidationId, call.InvalidationNonce)
		}
	}
}
