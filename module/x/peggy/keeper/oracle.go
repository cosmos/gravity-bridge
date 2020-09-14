package keeper

import (
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	EventTypeObservation        = "observation"
	AttributeKeyAttestationType = "attestation_type"
	AttributeKeyContract        = "bridge_contract"
	AttributeKeyNonce           = "nonce"
	AttributeKeyBridgeChainID   = "bridge_chain_id"
)

func (k Keeper) AddClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.Nonce, validator sdk.ValAddress) error {
	if err := k.storeClaim(ctx, claimType, nonce, validator); err != nil {
		return sdkerrors.Wrap(err, "claim")
	}
	validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, validator)
	att, err := k.tryAttestation(ctx, claimType, nonce, uint64(validatorPower))
	if err != nil {
		return err
	}
	if att.Certainty != types.CertaintyObserved || att.Status != types.ProcessStatusInit {
		return nil
	}

	// next process Attestation
	xCtx, commit := ctx.CacheContext()

	// end time was handled in adding claim and would return an error early
	// no need to re-check here
	if err := k.AttestationHandler.Handle(xCtx, *att); err != nil { // execute with a transient storage
		// log
		att.ProcessResult = types.ProcessResultFailure
	} else {
		att.ProcessResult = types.ProcessResultSuccess
		commit() // persist transient storage
	}
	att.Status = types.ProcessStatusProcessed

	// now store all updates
	k.storeAttestation(ctx, att)

	observationEvent := sdk.NewEvent(
		EventTypeObservation, // todo: revisit types for clients
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(AttributeKeyAttestationType, fmt.Sprintf("%X", att.ClaimType)), // todo: map to string
		sdk.NewAttribute(AttributeKeyContract, types.BridgeContractAddress.String()),
		sdk.NewAttribute(AttributeKeyBridgeChainID, types.BridgeContractChainID),
		sdk.NewAttribute(AttributeKeyNonce, string(nonce)), // todo: serialize with hex/ base64 ?

	)
	ctx.EventManager().EmitEvent(observationEvent)
	return nil
}

func (k Keeper) storeClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.Nonce, validator sdk.ValAddress) error {
	store := ctx.KVStore(k.storeKey)
	cKey := types.GetClaimKey(claimType, nonce, validator)
	r := store.Get(cKey)
	if r != nil {
		return types.ErrDuplicate
	}
	store.Set(cKey, []byte{}) // empty as all payload is in the key already (no gas costs)
	return nil
}

func (k Keeper) tryAttestation(ctx sdk.Context, claimType types.ClaimType, nonce types.Nonce, power uint64) (*types.Attestation, error) {
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
			Tally: types.AttestationTally{
				RequiredVotesPower: types.AttestationVotesPowerThreshold.Mul(power.Uint64()),
				RequiredVotesCount: types.AttestationVotesCountThreshold.Mul(uint64(count)),
			},
			SubmitTime:          now,
			ConfirmationEndTime: now.Add(types.AttestationPeriod),
		}
	}
	if err := att.AddConfirmation(now, power); err != nil {
		return nil, err
	}
	return att, nil
}

func (k Keeper) storeAttestation(ctx sdk.Context, att *types.Attestation) {
	aKey := types.GetAttestationKey(att.ClaimType, att.Nonce)
	store := ctx.KVStore(k.storeKey)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

func (k Keeper) GetAttestation(ctx sdk.Context, claimType types.ClaimType, nonce types.Nonce) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(claimType, nonce)
	bz := store.Get(aKey)
	if bz == nil {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &att)
	return &att
}

func (k Keeper) HasClaim(ctx sdk.Context, claimType types.ClaimType, nonce types.Nonce, validator sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetClaimKey(claimType, nonce, validator))
}

func (k Keeper) IterateClaims(ctx sdk.Context, cb func(key []byte, claimType types.ClaimType, nonce types.Nonce, validator sdk.ValAddress) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleClaimKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		rawKey := iter.Key()
		claimType, validator, nonce := types.SplitClaimKey(rawKey)
		if cb(rawKey, claimType, nonce, validator) {
			break
		}
	}
}
