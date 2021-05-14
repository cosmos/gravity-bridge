package keeper

import (
	"fmt"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO-JT: carefully look at atomicity of this function
func (k Keeper) Vote(
	ctx sdk.Context,
	event types.EthereumEvent,
	anyEvent *codectypes.Any,
) (*types.EthereumEventVoteRecord, error) {
	valAddr := k.GetOrchestratorValidator(ctx, event.GetValidator())
	if valAddr == nil {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processEthereumEventVoteRecord as well,
	// but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two events with the same nonce.
	// This prevents there being two ethereumEventVoteRecords with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if event.GetEventNonce() != lastEventNonce+1 {
		return nil, types.ErrNonContiguousEventNonce
	}

	// Tries to get an ethereumEventVoteRecord with the same eventNonce and event as the event that was submitted.
	voteRecord := k.GetEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.EventHash())

	// If it does not exist, create a new one.
	if voteRecord == nil {
		voteRecord = &types.EthereumEventVoteRecord{
			Accepted: false,
			Height:   uint64(ctx.BlockHeight()),
			Event:    anyEvent,
		}
	}

	// Add the validator's vote to this ethereumEventVoteRecord
	voteRecord.Votes = append(voteRecord.Votes, valAddr.String())

	k.SetEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.EventHash(), voteRecord)
	k.setLastEventNonceByValidator(ctx, valAddr, event.GetEventNonce())

	return voteRecord, nil
}

// TryEthereumEventVoteRecord checks if an ethereumEventVoteRecord has enough votes to be applied to the consensus state
// and has not already been marked Accepted, then calls processEthereumEventVoteRecord to actually apply it to the state,
// and then marks it Accepted and emits an event.
func (k Keeper) TryEthereumEventVoteRecord(ctx sdk.Context, voteRecord *types.EthereumEventVoteRecord) {
	event, err := k.UnpackEthereumEventVoteRecordEvent(voteRecord)
	if err != nil {
		panic("could not cast to event")
	}
	// If the ethereumEventVoteRecord has not yet been Accepted, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the ethereumEventVoteRecord from accidentally being applied twice.
	if !voteRecord.Accepted {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.EthereumEventVoteRecordPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		ethereumEventVoteRecordPower := sdk.NewInt(0)
		for _, validator := range voteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the ethereumEventVoteRecord power's sum
			ethereumEventVoteRecordPower = ethereumEventVoteRecordPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the ethereumEventVoteRecord is higher or equal to the threshold,
			// process the ethereumEventVoteRecord, set Accepted to true, and break
			if ethereumEventVoteRecordPower.GTE(requiredPower) {
				lastEventNonce := k.GetLastAcceptedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if event.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.setLastAcceptedEventNonce(ctx, event.GetEventNonce())
				k.SetLatestEthereumBlockHeight(ctx, event.GetBlockHeight())

				voteRecord.Accepted = true
				k.SetEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.EventHash(), voteRecord)

				k.processEthereumEventVoteRecord(ctx, voteRecord, event)
				k.emitAcceptedEvent(ctx, voteRecord, event)
				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process accepted ethereumEventVoteRecord")
	}
}

// processEthereumEventVoteRecord actually applies the ethereumEventVoteRecord to the consensus state
func (k Keeper) processEthereumEventVoteRecord(ctx sdk.Context, voteRecord *types.EthereumEventVoteRecord, event types.EthereumEvent) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.EthereumEventVoteRecordHandler.Handle(xCtx, *voteRecord, event); err != nil { // execute with a transient storage
		// If the ethereumEventVoteRecord fails, something has gone wrong and we can't recover it. Log and move on
		// The ethereumEventVoteRecord will still be marked "Accepted", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("ethereumEventVoteRecord failed",
			"cause", err.Error(),
			"event type", event.GetType(),
			"id", types.GetEthereumEventVoteRecordKey(event.GetEventNonce(), event.EventHash()),
			"nonce", fmt.Sprint(event.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// emitAcceptedEvent emits an event with information about an ethereumEventVoteRecord that has been applied to
// consensus state.
func (k Keeper) emitAcceptedEvent(ctx sdk.Context, voteRecord *types.EthereumEventVoteRecord, event types.EthereumEvent) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordType, string(event.GetType())),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		// todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordID,
			string(types.GetEthereumEventVoteRecordKey(event.GetEventNonce(), event.EventHash()))),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(event.GetEventNonce())),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// SetEthereumEventVoteRecord sets the ethereumEventVoteRecord in the store
func (k Keeper) SetEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, eventHash []byte, voteRecord *types.EthereumEventVoteRecord) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetEthereumEventVoteRecordKey(eventNonce, eventHash)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(voteRecord))
}

// GetEthereumEventVoteRecord return an ethereumEventVoteRecord given a nonce
func (k Keeper) GetEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, eventHash []byte) *types.EthereumEventVoteRecord {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetEthereumEventVoteRecordKey(eventNonce, eventHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var voteRecord types.EthereumEventVoteRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &voteRecord)
	return &voteRecord
}

// DeleteEthereumEventVoteRecord deletes an ethereumEventVoteRecord given an event nonce and event
// func (k Keeper) DeleteEthereumEventVoteRecord(ctx sdk.Context, eventNonce uint64, eventHash []byte, voteRecord *types.EthereumEventVoteRecord) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Delete(types.GetEthereumEventVoteRecordKeyWithHash(eventNonce, eventHash))
// }

// GetEthereumEventVoteRecordMapping returns a mapping of eventnonce -> ethereumEventVoteRecords at that nonce
func (k Keeper) GetEthereumEventVoteRecordMapping(ctx sdk.Context) (out map[uint64][]types.EthereumEventVoteRecord) {
	out = make(map[uint64][]types.EthereumEventVoteRecord)
	k.IterateEthereumVoteRecords(ctx, func(_ []byte, voteRecord types.EthereumEventVoteRecord) bool {
		event, err := k.UnpackEthereumEventVoteRecordEvent(&voteRecord)
		if err != nil {
			panic("couldn't cast to event")
		}

		if val, ok := out[event.GetEventNonce()]; !ok {
			out[event.GetEventNonce()] = []types.EthereumEventVoteRecord{voteRecord}
		} else {
			out[event.GetEventNonce()] = append(val, voteRecord)
		}
		return false
	})
	return
}

// IterateEthereumVoteRecords iterates through all ethereumEventVoteRecords
func (k Keeper) IterateEthereumVoteRecords(ctx sdk.Context, cb func([]byte, types.EthereumEventVoteRecord) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.EthereumEventVoteRecordKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		voteRecord := types.EthereumEventVoteRecord{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &voteRecord)
		// cb returns true to stop early
		if cb(iter.Key(), voteRecord) {
			return
		}
	}
}

// GetLastAcceptedEventNonce returns the latest accepted event nonce
func (k Keeper) GetLastAcceptedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastAcceptedEventNonceKey)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetLatestEthereumBlockHeight height gets the block height to of the last accepted ethereumEventVoteRecord from
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

// GetLastAcceptedSignerSetTx retrieves the last accepted validator set from the store
// WARNING: This value is not an up to date validator set on Ethereum, it is a validator set
// that AT ONE POINT was the one in the Gravity bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastAcceptedSignerSetTx(ctx sdk.Context) *types.SignerSetTx {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastAcceptedSignerSetTxKey)

	if len(bytes) == 0 {
		return nil
	}
	signerSetTx := types.SignerSetTx{}
	k.cdc.MustUnmarshalBinaryBare(bytes, &signerSetTx)
	return &signerSetTx
}

// SetLastAcceptedSignerSetTx updates the last accepted validator set in the store
func (k Keeper) SetLastAcceptedSignerSetTx(ctx sdk.Context, signerSetTx types.SignerSetTx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastAcceptedSignerSetTxKey, k.cdc.MustMarshalBinaryBare(&signerSetTx))
}

// setLastAcceptedEventNonce sets the latest accepted event nonce
func (k Keeper) setLastAcceptedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastAcceptedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		// TODO: I believe that params.SignedClaimsWindow is deprecated... how does this impact the comment below?
		// in the case that we have no existing value this is the first
		// time a validator is submitting a event. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		//
		// We could start at the LastAcceptedEventNonce but if we do that this
		// validator will be slashed, because they are responsible for making a event
		// on any ethereumEventVoteRecord that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest ethereumEventVoteRecord that is still within
		// the slashing window. Since we delete ethereumEventVoteRecords after the slashing window that's
		// just the lowest accepted event in the store. If no events have been submitted in for
		// params.SignedClaimsWindow we may have no ethereumEventVoteRecords in our nonce. At which point
		// the last accepted which is a persistent and never cleaned counter will suffice.
		lowestAccepted := k.GetLastAcceptedEventNonce(ctx)
		voteRecordMap := k.GetEthereumEventVoteRecordMapping(ctx)
		// no new events in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(voteRecordMap) == 0 {
			return lowestAccepted
		}
		for nonce, voteRecords := range voteRecordMap {
			for voteRecord := range voteRecords {
				if voteRecords[voteRecord].Accepted && nonce < lowestAccepted {
					lowestAccepted = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowestAccepted > 0 {
			return lowestAccepted - 1
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
