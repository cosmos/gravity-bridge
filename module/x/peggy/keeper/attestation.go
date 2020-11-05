package keeper

import (
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddClaim starts the following process chain:
// - Records that a given validator has made a claim about a given ethereum event, checking that the event nonce is contiguous
//   (non contiguous eventNonces indicate out of order events which can cause double spends)
// - Either creates a new attestation or adds the validator's vote to the existing attestation for this event
// - Checks if the attestation has enough votes to be considered "Observed", then attempts to apply it to the
//   consensus state (e.g. minting tokens for a deposit event)
// - If so, marks it "Observed" and emits an event
func (k Keeper) AddClaim(ctx sdk.Context, claimType types.ClaimType, eventNonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) (*types.Attestation, error) {
	if err := k.storeClaim(ctx, claimType, eventNonce, validator, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	att := k.voteForAttestation(ctx, claimType, eventNonce, details, validator)

	k.tryAttestation(ctx, att)

	k.SetAttestation(ctx, att)

	return att, nil
}

// storeClaim persists a claim. Fails when a claim submitted by an Eth signer does not increment the event nonce by exactly 1.
func (k Keeper) storeClaim(ctx sdk.Context, claimType types.ClaimType, eventNonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) error {
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator)
	if eventNonce != lastEventNonce+1 {
		return types.ErrNonContiguousEventNonce
	}
	k.setLastEventNonceByValidator(ctx, validator, eventNonce)

	// Store this nonce and the claim
	// TODO: This is not actually storing the claim. It can only be used to check if we have stored the same claim before.
	// We need to think this through more. This will be used for slashing later.
	store := ctx.KVStore(k.storeKey)
	cKey := types.GetClaimKey(claimType, eventNonce, validator, details)
	store.Set(cKey, []byte{}) // empty as all payload is in the key already (no gas costs)
	return nil
}

// voteForAttestation either gets the attestation for this claim from storage, or creates one if this is the first time a validator
// has submitted a claim for this exact event
func (k Keeper) voteForAttestation(
	ctx sdk.Context,
	claimType types.ClaimType,
	eventNonce types.UInt64Nonce,
	details types.AttestationDetails,
	validator sdk.ValAddress,
) *types.Attestation {
	// Tries to get an attestation with the same eventNonce and details as the claim that was submitted.
	att := k.GetAttestation(ctx, eventNonce, details)
	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			ClaimType:  claimType,
			EventNonce: eventNonce,
			Observed:   false,
			Details:    details,
		}
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, validator)

	return att
}

// tryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) tryAttestation(ctx sdk.Context, att *types.Attestation) {
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidently being applied twice.
	if !att.Observed {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			// Get the power of the current validator
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				k.processAttestation(ctx, att)
				att.Observed = true
				k.emitObservedEvent(ctx, att)
				break
			}
		}
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.Attestation) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, string(att.ClaimType)),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(att.EventNonce, att.Details))), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyNonce, att.EventNonce.String()),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation) {
	lastEventNonce := k.GetLastObservedEventNonce(ctx)
	if att.EventNonce != lastEventNonce+1 {
		// TODO: We need to figure out how to handle this situation, and whether it is even possible.
		// I'm panicking here because if attestations are applied to the consensus state out of order, it WILL cause a
		// double spend.
		// In theory, the fact that all votes on attestations are strictly ordered when the claim is submitted should mean
		// that this is impossible, but we should know for sure before removing the check. If it is possible, we need to
		// figure out how to recover.
		panic("attempting to apply events to state out of order")
	}
	k.setLastObservedEventNonce(ctx, att.EventNonce)

	// then execute in a new Tx so that we can store state on failure
	// TODO: It seems that the validator who puts an attestation over the threshold of votes will also
	// be charged for the gas of applying it to the consensus state. We should figure out a way to avoid this.
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", att.ClaimType,
			"id", types.GetAttestationKey(att.EventNonce, att.Details),
			"nonce", att.EventNonce.String(),
		)
	} else {
		commit() // persist transient storage
	}
}

func (k Keeper) SetAttestation(ctx sdk.Context, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(att.EventNonce, att.Details)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce types.UInt64Nonce, details types.AttestationDetails) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, details)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &att)
	return &att
}

func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) types.UInt64Nonce {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastObservedEventNonceKey())

	if len(bytes) == 0 {
		return types.NewUInt64Nonce(0)
	}
	return types.UInt64NonceFromBytes(bytes)
}

func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, nonce types.UInt64Nonce) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastObservedEventNonceKey(), nonce.Bytes())
}

func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) types.UInt64Nonce {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		return types.NewUInt64Nonce(0)
	}
	return types.UInt64NonceFromBytes(bytes)
}

func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce types.UInt64Nonce) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), nonce.Bytes())
}

func (k Keeper) HasClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetClaimKey(claimType, nonce, validator, details))
}

// func (k Keeper) IterateClaims(ctx sdk.Context, cb func(key []byte, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress) bool) {
// 	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleClaimKey)
// 	iter := prefixStore.Iterator(nil, nil)
// 	for ; iter.Valid(); iter.Next() {
// 		rawKey := iter.Key()
// 		claimType, validator, nonce := types.SplitClaimKey(rawKey)
// 		if cb(rawKey, claimType, nonce, validator) {
// 			break
// 		}
// 	}
// }
