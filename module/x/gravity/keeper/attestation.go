package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// AttestEvent signals an ethereum event as
// TODO: explain logic
func (k Keeper) AttestEvent(ctx sdk.Context, event types.EthereumEvent, orchestratorAddr sdk.AccAddress) error {
	validatorAddr := k.GetOrchestratorValidator(ctx, orchestratorAddr)
	if validatorAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two events with the same nonce
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validatorAddr)
	if event.GetEventNonce() != lastEventNonce+1 {
		return types.ErrNonContiguousEventNonce
	}

	// Tries to get an attestation with the same eventNonce and event as the event that was submitted.
	attestation, found := k.GetAttestation(ctx, event.GetEventNonce(), event.ClaimHash())
	if !found {
		anyClaim, err := types.PackClaim(event)
		if err != nil {
			return err
		}

		attestation = types.Attestation{
			Observed: false,
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	// Add the validator's vote to this attestation
	attestation.Votes = append(attestation.Votes, validatorAddr.String())

	k.SetAttestation(ctx, event.GetEventNonce(), event.ClaimHash(), attestation)
	// TODO: what is this for?
	k.setLastEventNonceByValidator(ctx, validatorAddr, event.GetEventNonce())
	return nil
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, attestation types.Attestation) {

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
		k.Logger(ctx).Debug("attestation threshold not met for event", "type", event.GetType().String(), "hash", event.ClaimHash())
		return
	}

	k.setLastObservedEventNonce(ctx, event.GetEventNonce())
	k.SetLastObservedEthereumBlockHeight(ctx, event.GetBlockHeight())

	attestation.Observed = true
	k.SetAttestation(ctx, event.GetEventNonce(), event.ClaimHash(), attestation)

	k.processAttestation(ctx, attestation, event)
	k.emitObservedEvent(ctx, attestation, event)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, attestation types.Attestation, event types.EthereumEvent) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.attestationHandler.HandleAttestation(xCtx, attestation); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"event type", event.GetType(),
			"id", types.GetAttestationKey(event.GetEventNonce(), event.ClaimHash()),
			"nonce", fmt.Sprint(event.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}
