package gravity

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/keeper"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Question: what here can be epoched?
	slashing(ctx, k)
	attestationTally(ctx, k)
	cleanupTimedOutBatchTxs(ctx, k)
	cleanupTimedOutContractCallTxs(ctx, k)
	createSignerSetTxs(ctx, k)
	pruneSignerSetTxs(ctx, k)
}

func createSignerSetTxs(ctx sdk.Context, k keeper.Keeper) {
	// Auto ValsetRequest Creation.
	// 1. If there are no valset requests, create a new one.
	// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
	//      This will make sure the unbonding validator has to provide an ethereum signature to a new signer set tx
	//	    that excludes him before he completely Unbonds.  Otherwise he will be slashed
	// 3. If power change between validators of CurrentValset and latest valset request is > 5%
	latestValset := k.GetLatestSignerSetTx(ctx)
	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)
	powerDiff := types.EthereumSigners(k.NewSignerSetTx(ctx).Signers).PowerDiff(latestValset.Signers)
	if (latestValset == nil) || (lastUnbondingHeight == uint64(ctx.BlockHeight())) || (powerDiff > 0.05) {
		k.SetOutgoingTx(ctx, k.NewSignerSetTx(ctx))
	}
}

func pruneSignerSetTxs(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	// Validator set pruning
	// prune all validator sets with a nonce less than the
	// last observed nonce, they can't be submitted any longer
	//
	// Only prune valsets after the signed valsets window has passed
	// so that slashing can occur the block before we remove them
	lastObserved := k.GetLastObservedSignerSetTx(ctx)
	currentBlock := uint64(ctx.BlockHeight())
	tooEarly := currentBlock < params.SignedSignerSetTxsWindow
	if lastObserved != nil && !tooEarly {
		earliestToPrune := currentBlock - params.SignedSignerSetTxsWindow
		for _, set := range k.GetSignerSetTxs(ctx) {
			if set.Nonce < lastObserved.Nonce && set.Height < earliestToPrune {
				k.DeleteOutgoingTx(ctx, set.GetStoreIndex())
			}
		}
	}
}

func slashing(ctx sdk.Context, k keeper.Keeper) {

	params := k.GetParams(ctx)

	// Slash validator for not confirming valset requests, batch requests and not attesting claims rightfully
	SignerSetTxSlashing(ctx, k, params)
	BatchSlashing(ctx, k, params)
	// TODO slashing for arbitrary logic is missing

	// TODO: prune validator sets, older than 6 months, this time is chosen out of an abundance of caution
	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k keeper.Keeper) {
	attmap := k.GetEthereumEventVoteRecordMapping(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	keys := make([]uint64, 0, len(attmap))
	for k := range attmap {
		keys = append(keys, k)
	}
	// Then we sort it
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range keys {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attmap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last attestation that was
			// observed. If it is not, we just move on to the next nonce. This will skip over all
			// attestations that have already been observed.
			//
			// Once we hit an event nonce that is one higher than the last observed event, we stop
			// skipping over this conditional and start calling tryAttestation (counting votes)
			// Once an attestation at a given event nonce has enough votes and becomes observed,
			// every other attestation at that nonce will be skipped, since the lastObservedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the attestation mapping, if there is one. This
			// nonce will once again be one higher than the lastObservedEventNonce.
			// If there is an attestation at this event nonce which has enough votes to be observed,
			// we skip the other attestations and move on to the next nonce again.
			// If no attestation becomes observed, when we get to the next nonce, every attestation in
			// it will be skipped. The same will happen for every nonce after that.
			if nonce == uint64(k.GetLastObservedEventNonce(ctx))+1 {
				k.TryEventVoteRecord(ctx, &att)
			}
		}
	}
}

// cleanupTimedOutBatchTxs deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatchTxs(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumHeight
	batches := k.GetBatchTxes(ctx)
	for _, batch := range batches {
		if batch.Timeout < ethereumHeight {
			k.CancelBatchTx(ctx, common.HexToAddress(batch.TokenContract), batch.Nonce)
		}
	}
}

// cleanupTimedOutBatchTxs deletes logic calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutContractCallTxs(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumHeight
	k.IterateOutgoingTxs(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		cctx, _ := otx.(*types.ContractCallTx)
		if cctx.Timeout < ethereumHeight {
			k.DeleteOutgoingTx(ctx, cctx.GetStoreIndex())
		}
		return true
	})
}

func SignerSetTxSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedValsetsWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedSignerSetTxsWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedSignerSetTxsWindow
	}

	unslashedsignerSetTxs := k.GetUnSlashedSignerSetTxs(ctx, maxHeight)

	// unslashedsignerSetTxs are sorted by nonce in ASC order
	for _, sstx := range unslashedsignerSetTxs {
		signatures := k.GetEthereumSignatures(ctx, sstx.GetStoreIndex())

		// SLASH BONDED VALIDTORS who didn't attest valset request
		for _, val := range k.StakingKeeper.GetBondedValidatorsByPower(ctx) {
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)

			//  Slash validator ONLY if he joined after valset is created
			if exist && valSigningInfo.StartHeight < int64(sstx.Height) {
				if _, found := signatures[val.GetOperator().String()]; !found {
					if !val.IsJailed() {
						// TODO: do we want to slash jailed validators?
						k.StakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
						k.StakingKeeper.Jail(ctx, consAddr)
					}
				}
			}
		}

		// SLASH UNBONDING VALIDATORS who didn't attest valset request
		blockTime := ctx.BlockTime().Add(k.StakingKeeper.GetParams(ctx).UnbondingTime)
		blockHeight := ctx.BlockHeight()
		unbondingValIterator := k.StakingKeeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
		defer unbondingValIterator.Close()

		// All unbonding validators
		for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
			unbondingValidators := k.GetUnbondingvalidators(unbondingValIterator.Value())

			for _, valAddr := range unbondingValidators.Addresses {
				addr, _ := sdk.ValAddressFromBech32(valAddr)
				validator, _ := k.StakingKeeper.GetValidator(ctx, sdk.ValAddress(addr))
				valConsAddr, _ := validator.GetConsAddr()
				valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, valConsAddr)

				// Only slash validators who joined after valset is created and they are unbonding and UNBOND_SLASHING_WINDOW didn't pass
				if exist && valSigningInfo.StartHeight < int64(sstx.Nonce) && validator.IsUnbonding() && sstx.Height < uint64(validator.UnbondingHeight)+params.UnbondSlashingSignerSetTxsWindow {
					// Check if validator has confirmed valset or not
					if _, found := signatures[validator.GetOperator().String()]; !found {
						if !validator.IsJailed() {
							// TODO: do we want to slash jailed validators
							k.StakingKeeper.Slash(ctx, valConsAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionSignerSetTx)
							k.StakingKeeper.Jail(ctx, valConsAddr)
						}
					}
				}
			}
		}
		// then we set the latest slashed valset  nonce
		k.SetLastSlashedSignerSetTxNonce(ctx, sstx.Nonce)
	}
}

func BatchSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedBatchesWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedBatchesWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedBatchesWindow
	} else {
		return
	}

	unslashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)
	for _, batch := range unslashedBatches {

		// SLASH BONDED VALIDTORS who didn't attest batch requests
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		signatures := k.GetEthereumSignatures(ctx, batch.GetStoreIndex())
		for _, val := range currentBondedSet {
			// Don't slash validators who joined after batch is created
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if exist && valSigningInfo.StartHeight > int64(batch.Height) {
				continue
			}

			if _, ok := signatures[val.GetOperator().String()]; !ok {
				k.StakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
				if !val.IsJailed() {
					k.StakingKeeper.Jail(ctx, consAddr)
				}
			}
		}
		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.Height)

	}
}
