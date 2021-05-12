package keeper

import (
	"fmt"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO-JT: carefully look at atomicity of this function
func (k Keeper) Attest(
	ctx sdk.Context,
	claim types.EthereumEvent,
	anyClaim *codectypes.Any,
) (*types.EthereumEventVoteRecord, error) {
	valAddr := k.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if valAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processEthereumEventVoteRecord as well,
	// but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two ethereumEventVoteRecords with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return nil, types.ErrNonContiguousEventNonce
	}

	// Tries to get an ethereumEventVoteRecord with the same eventNonce and claim as the claim that was submitted.
	att := k.GetEthereumEventVoteRecord(ctx, claim.GetEventNonce(), claim.ClaimHash())

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.EthereumEventVoteRecord{
			Accepted: false,
			Height:   uint64(ctx.BlockHeight()),
			Event:    anyClaim,
		}
	}

	// Add the validator's vote to this ethereumEventVoteRecord
	att.Votes = append(att.Votes, valAddr.String())

	k.SetEthereumEventVoteRecord(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
	k.setLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())

	return att, nil
}

// TryEthereumEventVoteRecord checks if an ethereumEventVoteRecord has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processEthereumEventVoteRecord to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryEthereumEventVoteRecord(ctx sdk.Context, att *types.EthereumEventVoteRecord) {
	claim, err := k.UnpackEthereumEventVoteRecordClaim(att)
	if err != nil {
		panic("could not cast to claim")
	}
	// If the ethereumEventVoteRecord has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the ethereumEventVoteRecord from accidentally being applied twice.
	if !att.Accepted {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.EthereumEventVoteRecordPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		ethereumEventVoteRecordPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the ethereumEventVoteRecord power's sum
			ethereumEventVoteRecordPower = ethereumEventVoteRecordPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the ethereumEventVoteRecord is higher or equal to the threshold,
			// process the ethereumEventVoteRecord, set Observed to true, and break
			if ethereumEventVoteRecordPower.GTE(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if claim.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.setLastObservedEventNonce(ctx, claim.GetEventNonce())
				k.SetLatestEthereumBlockHeight(ctx, claim.GetBlockHeight())

				att.Accepted = true
				k.SetEthereumEventVoteRecord(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)

				k.processEthereumEventVoteRecord(ctx, att, claim)
				k.emitObservedEvent(ctx, att, claim)
				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed ethereumEventVoteRecord")
	}
}

// processEthereumEventVoteRecord actually applies the ethereumEventVoteRecord to the consensus state
func (k Keeper) processEthereumEventVoteRecord(ctx sdk.Context, att *types.EthereumEventVoteRecord, claim types.EthereumEvent) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.EthereumEventVoteRecordHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the ethereumEventVoteRecord fails, something has gone wrong and we can't recover it. Log and move on
		// The ethereumEventVoteRecord will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("ethereumEventVoteRecord failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetEthereumEventVoteRecordKey(claim.GetEventNonce(), claim.ClaimHash()),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// emitObservedEvent emits an event with information about an ethereumEventVoteRecord that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.EthereumEventVoteRecord, claim types.EthereumEvent) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordType, string(claim.GetType())),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		// todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordID,
			string(types.GetEthereumEventVoteRecordKey(claim.GetEventNonce(), claim.ClaimHash()))),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(claim.GetEventNonce())),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// SetEthereumEventVoteRecord sets the ethereumEventVoteRecord in the store
func (k Keeper) SetEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.EthereumEventVoteRecord) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetEthereumEventVoteRecordKey(eventNonce, claimHash)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// GetEthereumEventVoteRecord return an ethereumEventVoteRecord given a nonce
func (k Keeper) GetEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.EthereumEventVoteRecord {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetEthereumEventVoteRecordKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.EthereumEventVoteRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &att)
	return &att
}

// DeleteEthereumEventVoteRecord deletes an ethereumEventVoteRecord given an event nonce and claim
func (k Keeper) DeleteEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.EthereumEventVoteRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetEthereumEventVoteRecordKeyWithHash(eventNonce, claimHash))
}

// GetEthereumEventVoteRecordMapping returns a mapping of eventnonce -> ethereumEventVoteRecords at that nonce
func (k Keeper) GetEthereumEventVoteRecordMapping(ctx sdk.Context) (out map[uint64][]types.EthereumEventVoteRecord) {
	out = make(map[uint64][]types.EthereumEventVoteRecord)
	k.IterateAttestaions(ctx, func(_ []byte, att types.EthereumEventVoteRecord) bool {
		claim, err := k.UnpackEthereumEventVoteRecordClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		if val, ok := out[claim.GetEventNonce()]; !ok {
			out[claim.GetEventNonce()] = []types.EthereumEventVoteRecord{att}
		} else {
			out[claim.GetEventNonce()] = append(val, att)
		}
		return false
	})
	return
}

// IterateAttestaions iterates through all ethereumEventVoteRecords
func (k Keeper) IterateAttestaions(ctx sdk.Context, cb func([]byte, types.EthereumEventVoteRecord) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.OracleEthereumEventVoteRecordKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.EthereumEventVoteRecord{}
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

// GetLatestEthereumBlockHeight height gets the block height to of the last observed ethereumEventVoteRecord from
// the store
func (k Keeper) GetLatestEthereumBlockHeight(ctx sdk.Context) types.LatestEthereumBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LatestEthereumBlockHeightKey)

	if len(bytes) == 0 {
		return types.LatestEthereumBlockHeight{
			CosmosBlockHeight:   0,
			EthereumBlockHeight: 0,
		}
	}
	height := types.LatestEthereumBlockHeight{}
	k.cdc.MustUnmarshalBinaryBare(bytes, &height)
	return height
}

// SetLatestEthereumBlockHeight sets the block height in the store.
func (k Keeper) SetLatestEthereumBlockHeight(ctx sdk.Context, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LatestEthereumBlockHeight{
		EthereumBlockHeight: ethereumHeight,
		CosmosBlockHeight:   uint64(ctx.BlockHeight()),
	}
	store.Set(types.LatestEthereumBlockHeightKey, k.cdc.MustMarshalBinaryBare(&height))
}

// GetLastObservedSignerSetTx retrieves the last observed validator set from the store
// WARNING: This value is not an up to date validator set on Ethereum, it is a validator set
// that AT ONE POINT was the one in the Gravity bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedSignerSetTxKey)

	if len(bytes) == 0 {
		return nil
	}
	valset := types.SignerSetTx{}
	k.cdc.MustUnmarshalBinaryBare(bytes, &valset)
	return &valset
}

// SetLastObservedSignerSetTx updates the last observed validator set in the store
func (k Keeper) SetLastObservedSignerSetTx(ctx sdk.Context, valset types.SignerSetTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedSignerSetTxKey, k.cdc.MustMarshalBinaryBare(&valset))
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
		// on any ethereumEventVoteRecord that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest ethereumEventVoteRecord that is still within
		// the slashing window. Since we delete ethereumEventVoteRecords after the slashing window that's
		// just the lowest observed event in the store. If no claims have been submitted in for
		// params.SignedClaimsWindow we may have no ethereumEventVoteRecords in our nonce. At which point
		// the last observed which is a persistent and never cleaned counter will suffice.
		lowestObserved := k.GetLastObservedEventNonce(ctx)
		attmap := k.GetEthereumEventVoteRecordMapping(ctx)
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(attmap) == 0 {
			return lowestObserved
		}
		for nonce, atts := range attmap {
			for att := range atts {
				if atts[att].Accepted && nonce < lowestObserved {
					lowestObserved = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowestObserved > 0 {
			return lowestObserved - 1
		}
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), types.UInt64Bytes(nonce))
}
