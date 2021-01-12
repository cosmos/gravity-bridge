package keeper

import (
	"fmt"
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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
func (k Keeper) AddClaim(ctx sdk.Context, details types.EthereumClaim) (*types.Attestation, error) {
	if err := k.storeClaim(ctx, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	att := k.voteForAttestation(ctx, details)

	k.tryAttestation(ctx, att, details)

	k.SetAttestation(ctx, att, details)

	return att, nil
}

// storeClaim persists a claim. Fails when a claim submitted by an Eth signer does not increment the event nonce by exactly 1.
func (k Keeper) storeClaim(ctx sdk.Context, details types.EthereumClaim) error {
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, sdk.ValAddress(details.GetClaimer()))
	if details.GetEventNonce() != lastEventNonce+1 {
		return types.ErrNonContiguousEventNonce
	}
	valAddr := k.GetOrchestratorValidator(ctx, details.GetClaimer())
	if valAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	k.setLastEventNonceByValidator(ctx, valAddr, details.GetEventNonce())
	// Store the claim
	genericClaim, _ := types.GenericClaimfromInterface(details)
	store := ctx.KVStore(k.storeKey)
	cKey := types.GetClaimKey(details)
	store.Set(cKey, k.cdc.MustMarshalBinaryBare(genericClaim))
	return nil
}

// voteForAttestation either gets the attestation for this claim from storage, or creates one if this is the first time a validator
// has submitted a claim for this exact event
func (k Keeper) voteForAttestation(
	ctx sdk.Context,
	details types.EthereumClaim,
) *types.Attestation {
	// Tries to get an attestation with the same eventNonce and details as the claim that was submitted.
	att := k.GetAttestation(ctx, details.GetEventNonce(), details)

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			EventNonce: details.GetEventNonce(),
			Observed:   false,
		}
	}

	valAddr := k.GetOrchestratorValidator(ctx, details.GetClaimer())
	if valAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, valAddr.String())

	return att
}

// tryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) tryAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				k.processAttestation(ctx, att, claim)
				att.Observed = true
				k.emitObservedEvent(ctx, att, claim)
				break
			}
		}
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, string(claim.GetType())),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(att.EventNonce, claim))), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(att.EventNonce)),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	lastEventNonce := k.GetLastObservedEventNonce(ctx)
	if att.EventNonce != uint64(lastEventNonce)+1 {
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
	if err := k.AttestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(att.EventNonce, claim),
			"nonce", fmt.Sprint(att.EventNonce),
		)
	} else {
		commit() // persist transient storage

		// TODO: after we commit, delete the outgoingtxbatch that this claim references
	}
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	store := ctx.KVStore(k.storeKey)
	att.ClaimHash = claim.ClaimHash()
	att.Height = uint64(ctx.BlockHeight())
	aKey := types.GetAttestationKey(att.EventNonce, claim)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// SetAttestationUnsafe sets the attestation w/o setting height and claim hash
func (k Keeper) SetAttestationUnsafe(ctx sdk.Context, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKeyWithHash(att.EventNonce, att.ClaimHash)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, details types.EthereumClaim) *types.Attestation {
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

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKeyWithHash(att.EventNonce, att.ClaimHash))
}

// GetAttestationMapping returns a mapping of eventnonce -> attestations at that nonce
func (k Keeper) GetAttestationMapping(ctx sdk.Context) (out map[uint64][]types.Attestation) {
	out = make(map[uint64][]types.Attestation)
	k.IterateAttestaions(ctx, func(_ []byte, att types.Attestation) bool {
		if val, ok := out[att.EventNonce]; !ok {
			out[att.EventNonce] = []types.Attestation{att}
		} else {
			out[att.EventNonce] = append(val, att)
		}
		return false
	})
	return
}

// IterateAttestaions iterates through all attestations
func (k Keeper) IterateAttestaions(ctx sdk.Context, cb func([]byte, types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(types.OracleAttestationKey)
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// setLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		//
		// We could start at the LastObservedEventNonce but if we do that this
		// validator will be slashed, because they are responsible for making a claim
		// on any attestation that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest attestation that is still within
		// the slashing window. Since we delete attestations after the slashing window that's
		// just the lowest observed event in the store. If no claims have been submitted in for
		// params.SignedClaimsWindow we may have no attestations in our nonce. At which point
		// the last observed which is a persistant and never cleaned counter will suffice.
		lowest_observed := k.GetLastObservedEventNonce(ctx)
		attmap := k.GetAttestationMapping(ctx)
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(attmap) == 0 {
			return lowest_observed
		}
		for nonce, atts := range attmap {
			for att := range atts {
				if atts[att].Observed && nonce < lowest_observed {
					lowest_observed = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowest_observed > 0 {
			return lowest_observed - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytes(bytes)
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), types.UInt64Bytes(nonce))
}

// HasClaim returns true if a claim exists
func (k Keeper) HasClaim(ctx sdk.Context, details types.EthereumClaim) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetClaimKey(details))
}

// IterateClaimsByValidatorAndType takes a validator key and a claim type and then iterates over these claims
func (k Keeper) IterateClaimsByValidatorAndType(ctx sdk.Context, claimType types.ClaimType, validatorKey sdk.ValAddress, cb func([]byte, types.EthereumClaim) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleClaimKey)
	prefix := []byte(validatorKey)
	iter := prefixStore.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		genericClaim := types.GenericClaim{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &genericClaim)
		// cb returns true to stop early
		if cb(iter.Key(), &genericClaim) {
			break
		}
	}
}

// GetClaimsByValidatorAndType returns the list of claims a validator has signed for
func (k Keeper) GetClaimsByValidatorAndType(ctx sdk.Context, claimType types.ClaimType, val sdk.ValAddress) (out []types.EthereumClaim) {
	k.IterateClaimsByValidatorAndType(ctx, claimType, val, func(_ []byte, claim types.EthereumClaim) bool {
		out = append(out, claim)
		return false
	})
	return
}
