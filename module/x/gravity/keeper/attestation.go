package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// AttestEvent signals an ethereum event as
// TODO: explain logic
func (k Keeper) AttestEvent(ctx sdk.Context, event types.EthereumEvent, validatorAddr sdk.ValAddress) error {
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two events with the same nonce
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validatorAddr)
	if event.GetNonce() != lastEventNonce+1 {
		return types.ErrNonContiguousEventNonce
	}

	eventHash := event.Hash()

	// Tries to get an attestation with the same hash as the event that has been submitted
	attestation, found := k.GetAttestation(ctx, eventHash)
	if !found {
		attestation = types.Attestation{
			EventID:       eventHash, // TODO: use hex bytes
			AttestedPower: 0,
			Height:        uint64(ctx.BlockHeight()),
		}
	}

	// Add the validator's vote to this attestation
	attestation.Votes = append(attestation.Votes, validatorAddr.String())

	k.SetAttestation(ctx, eventHash, attestation)
	// TODO: what is this for?
	k.setLastEventNonceByValidator(ctx, validatorAddr, event.GetNonce())
	return nil
}

// TallyAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TallyAttestation(ctx sdk.Context, attestation types.Attestation) {

	if attestation.Observed {
		panic("attempting to process observed attestation")
	}

	// Sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.

	// Sum the current powers of all validators who have voted and see if it passes the current threshold
	// TODO: The different integer types and math here needs a careful review
	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)
	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))

	attestationPower := sdk.ZeroInt()
	var thresholdMet bool

	for _, validator := range attestation.Votes {
		val, _ := sdk.ValAddressFromBech32(validator)

		validatorPower := k.stakingKeeper.GetLastValidatorPower(ctx, val)
		attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))

		// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
		// process the attestation, set Observed to true, and break
		if attestationPower.GTE(requiredPower) {
			thresholdMet = true
			break
		}
	}

	if !thresholdMet {
		k.Logger(ctx).Debug("attestation threshold not met for event", "type", event.GetType().String(), "hash", event.Hash())
		return
	}

	k.setLastObservedEventNonce(ctx, event.GetNonce())
	k.SetLastObservedEthereumBlockHeight(ctx, event.GetBlockHeight())

	attestation.Observed = true
	k.SetAttestation(ctx, event.GetNonce(), event.Hash(), attestation)

	k.processAttestation(ctx, attestation, event)
	k.emitObservedEvent(ctx, attestation, event)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, attestation types.Attestation, event types.EthereumEvent) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()

	if err := k.attestationHandler.OnAttestation(xCtx, attestation); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"event-type", event.GetType(),
			"event-hash", event.Hash(),
			"nonce", fmt.Sprint(event.GetNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, hash []byte) (types.Attestation, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	bz := store.Get(hash)
	if len(bz) == 0 {
		return types.Attestation{}, false
	}

	var attestation types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &attestation)

	return attestation, true
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, hash []byte, attestation types.Attestation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	store.Set(hash, k.cdc.MustMarshalBinaryBare(&attestation))
}

// DeleteAttestation deletes an attestation given an event hash
func (k Keeper) DeleteAttestation(ctx sdk.Context, hash []byte, attestation types.Attestation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	store.Delete(hash)
}
