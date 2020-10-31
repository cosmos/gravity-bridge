package keeper

import (
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddClaim starts the following process chain:
//  - persists the claim
//  - find or opens an attestation
//  - add weighted vote to the attestation
//  - calculates intermediate sum of all submitted votes for this claim
//  - when threshold is reached the attestation is marked as `observed`
//  - `observed` attestations are processed for state transition
//  - the process result is stored with the attestion
//  - an `observation` event is emitted
//
// Jehan's note: AddClaim is called every time an eth signer sends in a claim. Each time:
// - it adds the eth signer's vote to the claim using "tryAttestation". The votes are tallied in a struct called an
//   "attestation", and when they pass a threshold, an enum called "Certainty" is updated.
// - it checks "Certainty", and once the threshold is passed, "processAttestation" is called,
//   which checks the nonce, updates the nonce, and actually moves tokens.
//
// A problem exists, involving two perfectly legitimate, non duplicated bridge deposits. The first could have e.g. nonce 101, and the second nonce 102.
// If eth signers are not sending claims in to the Cosmos chain in strict nonce order, the second could happen to pass the vote threshold first, resulting
// in it being impossible to submit the first. This failure is non-obvious, and could happen if one validator's packet containing
// the claim for the first tx was dropped, or maybe just because of quirks of the Cosmos gossip network.
//
// To fix this here is a possible solution:
// - Enforce strict nonce ordering on claims from a given validator. Do not let a deposit be submitted unless the nonce is exactly one
//   higher than the last claim of the last deposit. This could be added with a nonce check in storeClaim.
//   Eth signers will need to know how to avoid this check failing, either by knowing how to retry earlier transactions,
//   or by not sending a claim before the last claim was successfully submitted.
// - Don't cause a state transition (e.g. sending tokens) immediately once a deposit has reached the threshold of votes (is "observed"). Every block,
//   attempt to play the "observed" deposit claims back and perform the state transitions. Stop if the nonce ever increases by more than 1.
//   The unused deposit claims stick around and can be retried next block. Eth signers still need to know how to retry claims, otherwise they will be slashed.
//
// Jehan's note: The nonce logic for the batches is completely different. As shown above, we have to reject deposit claims
// if all earlier deposit claims have not yet entered the system.
// Batch claims are totally different. For batches, we need to free the transactions of earlier batches when a later batch claims
// comes in. So given that this function is common to both types of claims, it is a bad place to implement any batch logic.
// We will have to think about how to split the code. We may want to make two versions of this whole function.
//
// Let's say there are two batches in the batch pool. Batches with nonce 101 and 102. Both are successfully submitted to the
// Ethereum chain. But due to quirks in the Cosmos gossip network, batch 102 is observed first. Batch 101's transactions are
// free and can be put in a new batch. This results in a double spend.
// So we need a similar mechanism to avoid out of order batch claims, just like the deposit claims. The difficulty is
// that if an earlier batch is never submitted onto Eth (a normal occurence), it will never be observed, and with a naive
// implementation of consecutive nonces, this will halt the bridge.
// We need to distinguish somehow between a batch that has not been observed because it never was submitted, and a batch that
// has not been observed because its claim got hung up in the Cosmos gossip network.
//
// This would not be an issue (AND deposit claims would not require a nonce at all), if we brought events in from Ethereum in blocks.
// Each ethereum block, which could contain events relating both to deposits and batches, and would have its own nonce (the block height)
// would be submitted as a claim. We would have the requirement that blocks be consecutive, as proposed for the deposit nonces above.
// This would eliminate the possibility of duplicate deposits, and eliminate the possibility of batches being "observed" out of order.
// It would also likely be more efficient.
func (k Keeper) AddClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) (*types.Attestation, error) {
	// Jehan's note: Seems like this just stores the individual claim for future reference (e.g. slashing)
	if err := k.storeClaim(ctx, claimType, nonce, validator, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
	// Jehan's note: This seems to be where votes for a given claim are counted. "Attestation" seems to be a claim
	// together with its total vote.
	att, err := k.tryAttestation(ctx, claimType, nonce, details, uint64(validatorPower))
	if err != nil {
		return nil, err
	}

	// this is a really strange conditional that needs to be simplified, just asking for trouble.
	// it is correct, but it's too difficult to read and there's too many ways for it to be true
	// or false that are non-obvious. Great way to have a double-spend bug TODO refactor
	// Jehan's note: This guards acceptance of a given claim. If this conditional does not return,
	// tokens are moved on the next line.
	if att.Certainty != types.CertaintyObserved || att.Status != types.ProcessStatusInit {
		k.storeAttestation(ctx, att)
		return att, nil
	}

	// next process Attestation
	// Jehan's note: The nonce is checked in here.
	if err := k.processAttestation(ctx, att); err != nil {
		return nil, err
	}

	// now store all updates
	k.storeAttestation(ctx, att)

	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, string(att.ClaimType)),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyAttestationID, string(att.ID())), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyNonce, nonce.String()),
	)
	ctx.EventManager().EmitEvent(observationEvent)
	return att, nil
}

// end time check was handled in adding claim and would return an error early
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation) error {
	// nonce > last one of same claim type
	if !att.Nonce.GreaterThan(k.GetLastAttestedNonce(ctx, att.ClaimType)) {
		return sdkerrors.Wrap(types.ErrOutdated, "nonce")
	}
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att); err != nil { // execute with a transient storage
		k.logger(ctx).Error("attestation failed", "cause", err.Error(), "claim type", att.ClaimType, "id", att.ID(), "nonce", att.Nonce.String())
		att.ProcessResult = types.ProcessResultFailure
	} else {
		att.ProcessResult = types.ProcessResultSuccess
		commit() // persist transient storage
	}
	att.Status = types.ProcessStatusProcessed
	k.setLastAttestedNonce(ctx, att.ClaimType, att.Nonce)
	return nil
}

// storeClaim persists a claim. Fails when a claim of given type and nonce was was submitted by the validator before
func (k Keeper) storeClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) error {
	store := ctx.KVStore(k.storeKey)
	cKey := types.GetClaimKey(claimType, nonce, validator, details)
	if store.Has(cKey) {
		return types.ErrDuplicate
	}
	store.Set(cKey, []byte{}) // empty as all payload is in the key already (no gas costs)
	return nil
}

var (
	hundred = sdk.NewUint(100)
	zero    = sdk.NewUint(0)
)

// tryAttestation loads an existing attestation for the given claim type and nonce and adds a vote.
// When none exists yet, a new attestation is instantiated (but not persisted here)
func (k Keeper) tryAttestation(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, details types.AttestationDetails, power uint64) (*types.Attestation, error) {
	now := ctx.BlockTime()
	att := k.GetAttestation(ctx, claimType, nonce)
	if att == nil {
		count := len(k.StakingKeeper.GetBondedValidatorsByPower(ctx))
		power := k.StakingKeeper.GetLastTotalPower(ctx)
		att = &types.Attestation{
			ClaimType:     claimType,
			Nonce:         nonce,
			Certainty:     types.CertaintyRequested,
			Status:        types.ProcessStatusInit,
			ProcessResult: types.ProcessResultUnknown,
			Details:       details,
			Tally: types.AttestationTally{
				TotalVotesPower:    zero,
				RequiredVotesPower: types.AttestationVotesPowerThreshold.MulUint64(power.Uint64()).Quo(hundred),
				RequiredVotesCount: types.AttestationVotesCountThreshold.MulUint64(uint64(count)).Quo(hundred).Uint64(),
			},
			SubmitTime:          now,
			ConfirmationEndTime: now.Add(types.AttestationPeriod),
		}
	}
	// Jehan's note: Adds the validators vote to a given claim
	if err := att.AddVote(now, power); err != nil {
		return nil, err
	}
	return att, nil
}

func (k Keeper) storeAttestation(ctx sdk.Context, att *types.Attestation) {
	aKey := types.GetAttestationKey(att.ClaimType, att.Nonce)
	store := ctx.KVStore(k.storeKey)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// GetAttestation loads an attestation for the given claim type and nonce. Returns nil when none exists
func (k Keeper) GetAttestation(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(claimType, nonce)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &att)
	return &att
}

func (k Keeper) HasClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetClaimKey(claimType, nonce, validator, details))
}

func (k Keeper) IterateAttestationByClaimTypeDesc(ctx sdk.Context, claimType types.ClaimType, cb func([]byte, types.Attestation) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleAttestationKey)
	iter := prefixStore.ReverseIterator(prefixRange(claimType.Bytes()))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var att types.Attestation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &att)
		if cb(iter.Key(), att) { // cb returns true to stop early
			return
		}
	}
	return
}

// GetLastProcessedAttestation returns attestation for given claim type or nil when none found
func (k Keeper) GetLastProcessedAttestation(ctx sdk.Context, claimType types.ClaimType) *types.Attestation {
	var result *types.Attestation
	k.IterateAttestationByClaimTypeDesc(ctx, claimType, func(_ []byte, att types.Attestation) bool {
		if att.Certainty != types.CertaintyObserved {
			return false
		}
		result = &att
		return true
	})
	return result
}

func (k Keeper) setLastAttestedNonce(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastNonceByClaimTypeSecondIndexKey(claimType, nonce), []byte{}) // store payload in key only for gas optimization
}

func (k Keeper) GetLastAttestedNonce(ctx sdk.Context, claimType types.ClaimType) *types.UInt64Nonce {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetLastNonceByClaimTypeSecondIndexKeyPrefix(claimType))
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	if !iter.Valid() {
		return nil
	}
	v := types.UInt64NonceFromBytes(iter.Key())
	return &v
}
