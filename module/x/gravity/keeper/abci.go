package keeper

import (
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Question: what here can be epoched?
	params := k.GetParams(ctx)

	k.slash(ctx, params)
	k.tallyAttestations(ctx)
	k.timeoutTxs(ctx)
	k.createEthSignerSet(ctx, params)
}

// Iterate over all attestations and tally the current result.
func (k Keeper) tallyAttestations(ctx sdk.Context) {
	attmap := k.AttestationMap(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	keys := make([]uint64, 0, len(attmap))
	for k := range attmap {
		keys = append(keys, k)
	}

	// Then we sort it
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range keys {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attmap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last attestation that was
			// observed. If it is not, we just move on to the next nonce. This will skip over all
			// attestations that have already been observed.
			//
			// Once we hit an event nonce that is one higher than the last observed event, we stop
			// skipping over this conditional and start calling tryAttestation (counting votes)
			// Once an attestation at a given event nonce has enough votes and becomes observed,
			// every other attestation at that nonce will be skipped, since the lastObservedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the attestation mapping, if there is one. This
			// nonce will once again be one higher than the lastObservedEventNonce.
			// If there is an attestation at this event nonce which has enough votes to be observed,
			// we skip the other attestations and move on to the next nonce again.
			// If no attestation becomes observed, when we get to the next nonce, every attestation in
			// it will be skipped. The same will happen for every nonce after that.
			if nonce == uint64(k.GetLastObservedEventNonce(ctx))+1 {
				k.TryAttestation(ctx, att)
			}
		}
	}
}

// timeoutTxs deletes the batch and logic call transactions that have passed
// their expiration height on Ethereum.
func (k Keeper) timeoutTxs(ctx sdk.Context) {
	// get the latest ethereum height set after attesting the events
	// TODO: pass as argument?
	info, found := k.GetEthereumInfo(ctx)
	if !found {
		panic("ethereum observed info not found")
	}

	// TODO: start iteration in desc order from height = info.Height
	// TODO: can we iterate once for over a height range [0, info.Height] instead of
	// once for every tx type
	k.IterateBatchTxs(ctx, func(tokenContract common.Address, txID tmbytes.HexBytes, batchTx types.BatchTx) bool {
		if batchTx.Timeout < info.Height {
			k.CancelBatchTx(ctx, tokenContract, txID, batchTx)
		}

		return false
	})

	k.IterateLogicCallTxs(ctx, func(invalidationID tmbytes.HexBytes, invalidationNonce uint64, tx types.LogicCallTx) bool {
		if tx.Timeout < info.Height {
			k.CancelLogicCallTx(ctx, invalidationID, invalidationNonce)
		}

		return false
	})
}

// Auto ValsetRequest Creation.
// 1. If there are no valset requests, create a new one. TODO: why? is it necessary?
// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
//    This will make sure the unbonding validator has to provide an attestation to a new Valset
//      that excludes him before he completely Unbonds.  Otherwise he will be slashed
// 3. If power change between validators of CurrentValset and latest valset request is > 5% // TODO: define percentage on params?
func (k Keeper) createEthSignerSet(ctx sdk.Context, params types.Params) {
	latestSignerSet, nonce, found := k.GetLatestEthSignerSet(ctx)
	if !found || len(latestSignerSet.Signers) == 0 {
		k.Logger(ctx).Debug("no signer set", "nonce", strconv.FormatUint(nonce, 64))
		k.SetEthSignerSetRequest(ctx)
		return
	}

	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)

	if lastUnbondingHeight == uint64(ctx.BlockHeight()) {
		k.Logger(ctx).Debug("validator unbonding", "height", strconv.FormatInt(ctx.BlockHeight(), 64))
		k.SetEthSignerSetRequest(ctx)
		return
	}

	currentSignerSet := k.GetCurrentEthSignerSet(ctx)
	if len(currentSignerSet.Signers) == 0 {
		k.Logger(ctx).Debug("current signer set is empty", "height", strconv.FormatUint(currentSignerSet.Height, 64))
		return
	}

	currentPower := sdk.NewDec(currentSignerSet.Signers.TotalPower())
	latestPower := sdk.NewDec(latestSignerSet.Signers.TotalPower())

	powerDiff := latestPower.Sub(currentPower).Abs().QuoInt64(100)

	if powerDiff.LTE(params.MaxSignerSetPowerDiff) {
		// power difference is below threshold, don't submit request
		return
	}

	k.Logger(ctx).Debug("signer set power diff larger than threshold", "diff", powerDiff.String())
	k.SetEthSignerSetRequest(ctx)
}
