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
func (k Keeper) AddClaim(ctx sdk.Context, claimType types.ClaimType, eventNonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) (*types.Attestation, error) {
	// storeClaim stores the claim by an individual validator. It also makes sure that the event nonce is incremented by exactly 1
	if err := k.storeClaim(ctx, claimType, eventNonce, validator, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	// validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
	// Jehan's note: This seems to be where votes for a given claim are counted. "Attestation" seems to be a claim
	// together with its total vote.
	att, err := k.tryAttestation(ctx, claimType, eventNonce, details, validator)
	if err != nil {
		return nil, err
	}
	k.storeAttestation(ctx, att)

	// next process Attestation if it has been observed and has not already been processed.
	if att.Observed && !att.Processed {
		if err := k.processAttestation(ctx, att); err != nil {
			return nil, err
		}

		observationEvent := sdk.NewEvent(
			types.EventTypeObservation,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAttestationType, string(att.ClaimType)),
			sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx).String()),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(att.ID())), // todo: serialize with hex/ base64 ?
			sdk.NewAttribute(types.AttributeKeyNonce, eventNonce.String()),
		)
		ctx.EventManager().EmitEvent(observationEvent)
	}

	return att, nil
}

// end time check was handled in adding claim and would return an error early
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation) error {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		k.logger(ctx).Error("attestation failed", "cause", err.Error(), "claim type", att.ClaimType, "id", att.ID(), "nonce", att.EventNonce.String())
		// att.ProcessResult = types.ProcessResultFailure
	} else {
		// att.ProcessResult = types.ProcessResultSuccess
		commit() // persist transient storage
	}
	att.Processed = true
	// att.Status = types.ProcessStatusProcessed
	return nil
}

// storeClaim persists a claim. Fails when a claim submitted by an Eth signer does not increment the event nonce by exactly 1.
func (k Keeper) storeClaim(ctx sdk.Context, claimType types.ClaimType, eventNonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) error {
	store := ctx.KVStore(k.storeKey)
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator)
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	if eventNonce != lastEventNonce+1 {
		println("EVENT NONCE, LAST EVENT NONCE", eventNonce, lastEventNonce)
		return types.ErrNonContiguousEventNonce
	}
	// Store this nonce and the claim
	k.SetLastEventNonceByValidator(ctx, validator, eventNonce)
	// TODO: This is not actually storing the claim. It can only be used to check if we have stored the same claim before.
	// We need to think this through more. This will be used for slashing later.
	cKey := types.GetClaimKey(claimType, eventNonce, validator, details)
	store.Set(cKey, []byte{}) // empty as all payload is in the key already (no gas costs)
	return nil
}

var (
	hundred = sdk.NewUint(100)
	zero    = sdk.NewUint(0)
)

// tryAttestation loads an existing attestation for the given claim type and nonce and adds a vote.
// When none exists yet, a new attestation is instantiated (but not persisted here)
func (k Keeper) tryAttestation(ctx sdk.Context, claimType types.ClaimType, eventNonce types.UInt64Nonce, details types.AttestationDetails, validator sdk.ValAddress) (*types.Attestation, error) {
	att := k.GetAttestation(ctx, claimType, eventNonce)
	// TODO: We need to check here that the attestation that has been stored is about the same exact claim
	// that we currently processing. Right now it just appears that validators are voting on whether there has been a
	// claim with the same nonce.
	if att == nil {
		// count := len(k.StakingKeeper.GetBondedValidatorsByPower(ctx))
		// power := k.StakingKeeper.GetLastTotalPower(ctx)
		att = &types.Attestation{
			ClaimType:  claimType,
			EventNonce: eventNonce,
			Observed:   false,
			Processed:  false,
			Details:    details,
		}
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, validator)

	// Sum the current powers of all validators who have voted and see if it passes the current threshold
	// TODO: The different integer types and math here needs a careful review
	totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
	var attestationPower sdk.Int
	for _, validator := range att.Votes {
		// Get the power of the current validator
		validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
		// Add it to the attestation power's sum
		attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
		// If the power of all the validators that have voted on the attestation is higher than the threshold,
		// set Observed to true and break
		if attestationPower.GTE(requiredPower) {
			att.Observed = true
			break
		}
	}

	return att, nil
}

func (k Keeper) storeAttestation(ctx sdk.Context, att *types.Attestation) {
	aKey := types.GetAttestationKey(att.ClaimType, att.EventNonce)
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
		if att.Observed != true {
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

func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) types.UInt64Nonce {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		return types.NewUInt64Nonce(0)
	}
	return types.UInt64NonceFromBytes(bytes)
}

func (k Keeper) SetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce types.UInt64Nonce) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), nonce.Bytes())
}

func (k Keeper) SetBridgeObservedSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, signature []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeObservedSignatureKey(claimType, nonce, validator), signature)
}

func (k Keeper) GetBridgeObservedSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetBridgeObservedSignatureKey(claimType, nonce, validator))
}

func (k Keeper) HasBridgeObservedSignature(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeObservedSignatureKey(claimType, nonce, validator))
}

func (k Keeper) IterateBridgeObservedSignatures(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, cb func(_ []byte, sig []byte) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetBridgeObservedSignatureKeyPrefix(claimType))
	iter := prefixStore.Iterator(prefixRange(nonce.Bytes()))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(iter.Key(), iter.Value()) {
			break
		}
	}
}
