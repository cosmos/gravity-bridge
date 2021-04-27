package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// AttestEvent signals an ethereum event as
// TODO: explain logic
func (k Keeper) AttestEvent(ctx sdk.Context, event types.EthereumEvent, validator stakingtypes.ValidatorI) (tmbytes.HexBytes, error) {
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two events with the same nonce

	// TODO: why is this per validator?
	// lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator.GetOperator())
	// if event.GetNonce() != lastEventNonce+1 {
	// 	return types.ErrNonContiguousEventNonce
	// }

	eventID := event.Hash()

	// Tries to get an attestation with the same hash as the event that has been submitted
	attestation, found := k.GetAttestation(ctx, eventID)
	if !found {
		attestation = types.Attestation{
			EventID:       eventID,
			Votes:         []string{},
			AttestedPower: 0,
			Height:        uint64(ctx.BlockHeight()),
		}

		k.SetEthereumEvent(ctx, eventID, event)
	}

	// Add the validator's vote to this attestation
	attestation.Votes = append(attestation.Votes, validator.GetOperator().String())
	attestation.AttestedPower += validator.GetConsensusPower()

	k.SetAttestation(ctx, eventID, attestation)
	// TODO: what is this for?
	// k.setLastEventNonceByValidator(ctx, validatorAddr, event.GetNonce())
	return nil, nil
}

// TallyAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TallyAttestation(ctx sdk.Context, hash tmbytes.HexBytes, attestation types.Attestation) {
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
		k.Logger(ctx).Debug("attestation threshold not met for event", "event-id", attestation.EventID.String())
		return
	}

	// fetch the event to set the ethereum info
	event, found := k.GetEthereumEvent(ctx, attestation.EventID)
	if !found {
		panic(fmt.Errorf("event with ID %s not found for observed attestation", attestation.EventID))
	}

	// TODO: figure nonces
	// k.setLastObservedEventNonce(ctx, event.GetNonce())

	// now that the the event is attested (observed), we set the ethereum info to
	// the store

	info, found := k.GetEthereumInfo(ctx)
	if !found {
		panic("ethereum info not found")
	}

	// we only override the latest ethereum info if the block height from the
	// event is greater than the latest seen height
	if info.Height < event.GetEthereumHeight() {
		info = types.EthereumInfo{
			Timestamp: ctx.BlockTime(),
			Height:    event.GetEthereumHeight(),
		}

		k.SetEthereumInfo(ctx, info)
	}

	// FIXME: define an attestation key that is not dependent on the event ID
	// TODO: Ideally we should have multiple events attested at the same time?
	k.SetAttestation(ctx, event.Hash(), attestation)

	k.processAttestation(ctx, attestation)
	// k.emitObservedEvent(ctx, attestation, event)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, attestation types.Attestation) {
	// then execute in a new Tx so that we can store state on failure
	cacheCtx, commit := ctx.CacheContext()

	if err := k.attestationHandler.OnAttestation(cacheCtx, attestation); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed",
			"event-id", attestation.EventID.String(),
			"error", err.Error(),
		)
	} else {
		commit() // persist transient storage
	}
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, hash tmbytes.HexBytes) (types.Attestation, bool) {
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
func (k Keeper) SetAttestation(ctx sdk.Context, hash tmbytes.HexBytes, attestation types.Attestation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	store.Set(hash, k.cdc.MustMarshalBinaryBare(&attestation))
}

func (k Keeper) IterateAttestations(ctx sdk.Context, cb func(hash tmbytes.HexBytes, attestation types.Attestation) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {

		hash := tmbytes.HexBytes(iterator.Key()[1:])

		var attestation types.Attestation
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &attestation)
		if cb(hash, attestation) {
			break
		}
	}

}

// DeleteAttestation deletes an attestation given an event hash
func (k Keeper) DeleteAttestation(ctx sdk.Context, hash tmbytes.HexBytes, attestation types.Attestation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AttestationKeyPrefix)
	store.Delete(hash)
}
