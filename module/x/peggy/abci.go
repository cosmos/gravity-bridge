package peggy

import (
	"sort"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Question: what here can be epoched?
	slashing(ctx, k)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	cleanupTimedOutLogicCalls(ctx, k)
}

func slashing(ctx sdk.Context, k keeper.Keeper) {

	params := k.GetParams(ctx)

	// Slash validator for not confirming valset requests, batch requests and not attesting claims rightfully
	ValsetSlashing(ctx, k, params)
	BatchSlashing(ctx, k, params)
	// TODO slashing for arbitrary logic is missing
	ClaimsSlashing(ctx, k, params)

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
				k.TryAttestation(ctx, &att)
			}
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

// cleanupTimedOutBatches deletes logic calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutLogicCalls(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumBlockHeight
	calls := k.GetOutgoingLogicCalls(ctx)
	for _, call := range calls {
		if call.Timeout < ethereumHeight {
			k.CancelOutgoingLogicCall(ctx, call.InvalidationId, call.InvalidationNonce)
		}
	}
}

func ValsetSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedValsetsWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedValsetsWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedValsetsWindow
	}

	unslashedValsets := k.GetUnSlashedValsets(ctx, maxHeight)

	// unslashedValsets are sorted by nonce in ASC order
	// Question: do we need to sort each time? See if this can be epoched
	for _, vs := range unslashedValsets {
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		confirms := k.GetValsetConfirms(ctx, vs.Nonce)
		for _, val := range currentBondedSet {

			// Don't slash validators who joined after valset is created
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if exist && valSigningInfo.StartHeight > int64(vs.Nonce) {
				continue
			}

			// Don't slash validators who are unbonded and UNBOND_SLASHING_WINDOW has passed
			if val.UnbondingHeight > 0 && vs.Nonce > uint64(val.UnbondingHeight)+params.UnbondSlashingValsetsWindow {
				continue
			}

			// Check if validator has confirmed valset or not
			found := false
			for _, conf := range confirms {
				if conf.EthAddress == k.GetEthAddress(ctx, val.GetOperator()) {
					found = true
					break
				}
			}

			// slash validators for not confirming valsets
			if !found {
				cons, _ := val.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionValset)
				k.StakingKeeper.Jail(ctx, cons)

			}
		}

		// then we set the latest slashed valset  nonce
		k.SetLastSlashedValsetNonce(ctx, vs.Nonce)
	}

	// on the latest validator set, check for change in power against
	// current, and emit a new validator set if the change in power >5%
	latestValset := k.GetLatestValset(ctx)
	if latestValset == nil {
		k.SetValsetRequest(ctx)
	} else if types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(latestValset.Members) > 0.05 {
		k.SetValsetRequest(ctx)
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
	}

	unslashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)
	for _, batch := range unslashedBatches {
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
		for _, val := range currentBondedSet {
			// Don't slash validators who joined after batch is created
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if exist && valSigningInfo.StartHeight > int64(batch.Block) {
				continue
			}

			// Don't slash validators who are unbonded and UNBOND_SLASHING_WINDOW has passed
			if val.UnbondingHeight > 0 && batch.Block > uint64(val.UnbondingHeight)+params.UnbondSlashingBatchWindow {
				continue
			}

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

		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.Block)

	}
}

func ClaimsSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {
	// #3 condition
	// Oracle events MsgDepositClaim, MsgWithdrawClaim

	attmap := k.GetAttestationMapping(ctx)
	for _, atts := range attmap {
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
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
}

// TestingEndBlocker is a second endblocker function only imported in the Gravity codebase itself
// if you are a consuming Cosmos chain DO NOT IMPORT THIS, it simulates a chain using the arbitrary
// logic API to request logic calls
func TestingEndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// if this is nil we have not set our test outgoing logic call yet
	if k.GetOutgoingLogicCall(ctx, []byte("GravityTesting"), 0).Payload == nil {
		// TODO this call isn't actually very useful for testing, since it always
		// throws, being just junk data that's expected. But it prevents us from checking
		// the full lifecycle of the call. We need to find some way for this to read data
		// and encode a simple testing call, probably to one of the already deployed ERC20
		// contracts so that we can get the full lifecycle.
		token := []*types.ERC20Token{&types.ERC20Token{
			Contract: "0x7580bfe88dd3d07947908fae12d95872a260f2d8",
			Amount:   sdk.NewIntFromUint64(5000),
		}}
		_ = types.OutgoingLogicCall{
			Transfers:            token,
			Fees:                 token,
			LogicContractAddress: "0x510ab76899430424d209a6c9a5b9951fb8a6f47d",
			Payload:              []byte("fake bytes"),
			Timeout:              10000,
			InvalidationId:       []byte("GravityTesting"),
			InvalidationNonce:    1,
		}
		//k.SetOutgoingLogicCall(ctx, &call)
	}
}
