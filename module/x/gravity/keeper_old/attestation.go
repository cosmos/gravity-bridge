package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// Attest affirms a given ethereum claim as true
// TODO: explain logic
func (k Keeper) Attest(ctx sdk.Context, claim types.EthereumClaim) error {
	valAddr := k.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if valAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return types.ErrNonContiguousEventNonce
	}

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	attestation, found := k.GetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash())
	if !found {
		anyClaim, err := types.PackClaim(claim)
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
	attestation.Votes = append(attestation.Votes, valAddr.String())

	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), attestation)
	// TODO: what is this for?
	k.setLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())
	return nil
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, attestation types.Attestation) {
	claim, err := types.UnpackClaim(attestation.Claim)
	if err != nil {
		panic(err)
	}

	if attestation.Observed {
		panic("attempting to process observed attestation")
	}

	// Sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.

	// Sum the current powers of all validators who have voted and see if it passes the current threshold
	// TODO: The different integer types and math here needs a careful review
	totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))

	attestationPower := sdk.ZeroInt()
	var thresholdMet bool

	for _, validator := range attestation.Votes {
		val, _ := sdk.ValAddressFromBech32(validator)

		validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
		attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))

		// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
		// process the attestation, set Observed to true, and break
		if attestationPower.GTE(requiredPower) {
			thresholdMet = true
			break
		}
	}

	if !thresholdMet {
		k.logger(ctx).Debug("attestation threshold not met for claim", "type", claim.GetType().String(), "hash", claim.ClaimHash())
		return
	}

	k.setLastObservedEventNonce(ctx, claim.GetEventNonce())
	k.SetLastObservedEthereumBlockHeight(ctx, claim.GetBlockHeight())

	attestation.Observed = true
	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), attestation)

	k.processAttestation(ctx, attestation, claim)
	k.emitObservedEvent(ctx, attestation, claim)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, attestation types.Attestation, claim types.EthereumClaim) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, attestation, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, attestation types.Attestation, claim types.EthereumClaim) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, string(claim.GetType())),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(claim.GetEventNonce())),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, attestation types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAttestationKey(eventNonce, claimHash)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&attestation))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) (types.Attestation, bool) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return types.Attestation{}, false
	}

	var attestation types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &attestation)

	return attestation, true
}

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, attestation *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKeyWithHash(eventNonce, claimHash))
}

// IterateAttestationByNonce iterates through all attestations with a given event nonce
func (k Keeper) IterateAttestationByNonce(ctx sdk.Context, nonce uint64, cb func(types.Attestation) bool) {
	key := append([]byte(types.OracleAttestationKey), sdk.Uint64ToBigEndian(nonce)...)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), key)
	iter := store.Iterator(nil, nil)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var attestation types.Attestation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &attestation)

		if cb(attestation) {
			return
		}
	}
}

// ReverseIterateClaimsWindow iterates through all attestations in desc order within the current claims
// window [nonce - params.SignedClaimsWindow, nonce]
func (k Keeper) ReverseIterateClaimsWindow(ctx sdk.Context, nonce uint64, cb func(types.Attestation) bool) {
	params := k.GetParams(ctx)
	windowStartHeightBz := sdk.Uint64ToBigEndian(nonce - params.SignedClaimsWindow)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OracleAttestationKey))

	// TODO: test
	iter := store.ReverseIterator(sdk.Uint64ToBigEndian(nonce), windowStartHeightBz)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var attestation types.Attestation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &attestation)

		if cb(attestation) {
			return
		}
	}
}

// IterateAttestation iterates through all attestations
func (k Keeper) IterateAttestation(ctx sdk.Context, cb func(types.Attestation) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OracleAttestationKey))

	iter := store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var attestation types.Attestation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &attestation)

		if cb(attestation) {
			return
		}
	}
}

// GetLastObservedEthereumBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEthereumBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastObservedEthereumBlockHeightKey)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetLastObservedEthereumBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedEthereumBlockHeight(ctx sdk.Context, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEthereumBlockHeightKey, sdk.Uint64ToBigEndian(ethereumHeight))
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)
	if len(bytes) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bytes)
}

// setLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
// TODO: clean up
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bz) > 0 {
		return sdk.BigEndianToUint64(bz)
	}

	// in the case that we have no existing value this is the first
	// time a validator is submitting a claim. Since we don't want to force
	// them to replay the entire history of all events ever we can't start
	// at zero
	nonce := k.GetLastObservedEventNonce(ctx)
	if nonce == 0 {
		return 0
	}

	lastObservedNonce := nonce

	// We could return the LastObservedEventNonce but if we do that this
	// validator will be slashed, because they are responsible for making a claim
	// on any attestation that has not yet passed the slashing window.
	//
	// Therefore we need to return to them the latest attestation that is still within
	// the slashing window. Since we delete attestations after the slashing window that's
	// just the latest observed event in the store. If no claims have been submitted in for
	// params.SignedClaimsWindow we may have no attestations in our nonce. At which point
	// the last observed which is a persistant and never cleaned counter will suffice.

	k.ReverseIterateClaimsWindow(ctx, nonce-1, func(attestation types.Attestation) bool {
		if !attestation.Observed {
			// continue until either we encounter an observed attestation or we process a nonce less than
			// than the slashing window start height
			return false
		}

		claim, err := types.UnpackClaim(attestation.Claim)
		if err != nil {
			panic(err)
		}

		nonce = claim.GetEventNonce()
		return true // break iteration
	})

	if nonce == lastObservedNonce {
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		return nonce
	}

	// return the latest event minus one so that the validator
	// can submit that event and avoid slashing.

	return nonce - 1
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), sdk.Uint64ToBigEndian(nonce))
}
