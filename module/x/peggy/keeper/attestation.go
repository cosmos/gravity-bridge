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
func (k Keeper) AddClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.UInt64Nonce, validator sdk.ValAddress, details types.AttestationDetails) (*types.Attestation, error) {
	if err := k.storeClaim(ctx, claimType, nonce, validator, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
	att, err := k.tryAttestation(ctx, claimType, nonce, details, uint64(validatorPower))
	if err != nil {
		return nil, err
	}
	if att.Certainty != types.CertaintyObserved || att.Status != types.ProcessStatusInit {
		return att, nil
	}

	// next process Attestation
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
