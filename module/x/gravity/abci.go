package gravity

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/keeper"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	// Question: what here can be epoched?
	slashing(ctx, k)
	ethereumEventVoteRecordTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	cleanupTimedOutLogicCalls(ctx, k)
	createSignerSetTxs(ctx, k)
	pruneSignerSetTxs(ctx, k, params)
	// TODO: prune events, ethereumEventVoteRecords when they pass in the handler
}

func createSignerSetTxs(ctx sdk.Context, k keeper.Keeper) {
	// Auto SignerSetTx Creation.
	// WARNING: do not use k.GetLastAcceptedSignerSetTx in this function, it *will* result in losing control of the bridge
	// 1. If there are no signer set txs, create a new one.
	// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
	//      This will make sure the unbonding validator has to provide an ethereumEventVoteRecord to a new SignerSetTx
	//	    that excludes him before he completely Unbonds.  Otherwise he will be slashed
	// 3. If power change between validators of CurrentSignerSetTx and latest signer set tx is > 5%
	latestSignerSetTx := k.GetLatestSignerSetTx(ctx)
	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)

	if (latestSignerSetTx == nil) || (lastUnbondingHeight == uint64(ctx.BlockHeight())) || (types.EthereumSigners(k.CreateSignerSetTx(ctx).Members).PowerDiff(latestSignerSetTx.Members) > 0.05) {
		// Store signer set tx
		k.SetSignerSetTx(ctx)
	}
}

func pruneSignerSetTxs(ctx sdk.Context, k keeper.Keeper, params types.Params) {
	// Validator set pruning
	// prune all validator sets with a nonce less than the
	// last accepted nonce, they can't be submitted any longer
	//
	// Only prune signer set txs after the signed signer set txs window has passed
	// so that slashing can occur the block before we remove them
	lastAccepted := k.GetLastAcceptedSignerSetTx(ctx)
	currentBlock := uint64(ctx.BlockHeight())
	tooEarly := currentBlock < params.SignedSignerSetTxsWindow
	if lastAccepted != nil && !tooEarly {
		earliestToPrune := currentBlock - params.SignedSignerSetTxsWindow
		sets := k.GetSignerSetTxs(ctx)
		for _, set := range sets {
			if set.Nonce < lastAccepted.Nonce && set.Height < earliestToPrune {
				k.DeleteSignerSetTx(ctx, set.Nonce)
			}
		}
	}
}

func slashing(ctx sdk.Context, k keeper.Keeper) {

	params := k.GetParams(ctx)

	// Slash validator for not signing signer set txs, batch txs
	SignerSetTxSlashing(ctx, k, params)
	BatchSlashing(ctx, k, params)
	// TODO slashing for contract call tx signatures is missing

}

// Iterate over all ethereumEventVoteRecords currently being voted on in order of nonce and
// "Accept" those who have passed the threshold. Break the loop once we see
// an ethereumEventVoteRecord that has not passed the threshold
func ethereumEventVoteRecordTally(ctx sdk.Context, k keeper.Keeper) {
	voteRecordMap := k.GetEthereumEventVoteRecordMapping(ctx)
	// We make a slice with all the event nonces that are in the ethereumEventVoteRecord mapping
	nonces := make([]uint64, 0, len(voteRecordMap))
	for k := range voteRecordMap {
		nonces = append(nonces, k)
	}
	// Then we sort it
	sort.Slice(nonces, func(i, j int) bool { return nonces[i] < nonces[j] })

	// This iterates over all keys (event nonces) in the ethereumEventVoteRecord mapping. Each value contains
	// a slice with one or more ethereumEventVoteRecords at that event nonce. There can be multiple ethereumEventVoteRecords
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range nonces {
		// This iterates over all ethereumEventVoteRecords at a particular event nonce.
		// They are ordered by when the first ethereumEventVoteRecord at the event nonce was received.
		// This order is not important.
		for _, voteRecord := range voteRecordMap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last ethereumEventVoteRecord that was
			// accepted. If it is not, we just move on to the next nonce. This will skip over all
			// ethereumEventVoteRecords that have already been accepted.
			//
			// Once we hit an event nonce that is one higher than the last accepted event, we stop
			// skipping over this conditional and start calling tryEthereumEventVoteRecord (counting votes)
			// Once an ethereumEventVoteRecord at a given event nonce has enough votes and becomes accepted,
			// every other ethereumEventVoteRecord at that nonce will be skipped, since the lastAcceptedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the ethereumEventVoteRecord mapping, if there is one. This
			// nonce will once again be one higher than the lastAcceptedEventNonce.
			// If there is an ethereumEventVoteRecord at this event nonce which has enough votes to be accepted,
			// we skip the other ethereumEventVoteRecords and move on to the next nonce again.
			// If no ethereumEventVoteRecord becomes accepted, when we get to the next nonce, every ethereumEventVoteRecord in
			// it will be skipped. The same will happen for every nonce after that.
			if nonce == uint64(k.GetLastAcceptedEventNonce(ctx))+1 {
				k.TryEthereumEventVoteRecord(ctx, &voteRecord)
			}
		}
	}
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be accepted. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatches(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLatestEthereumBlockHeight(ctx).EthereumBlockHeight
	batches := k.GetBatchTxs(ctx)
	for _, batch := range batches {
		if batch.BatchTimeout < ethereumHeight {
			k.CancelBatchTx(ctx, batch.TokenContract, batch.BatchNonce)
		}
	}
}

// cleanupTimedOutBatches deletes contract calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be accepted. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutLogicCalls(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLatestEthereumBlockHeight(ctx).EthereumBlockHeight
	calls := k.GetContractCallTxs(ctx)
	for _, call := range calls {
		if call.Timeout < ethereumHeight {
			k.CancelContractCallTx(ctx, call.InvalidationId, call.InvalidationNonce)
		}
	}
}

func SignerSetTxSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedSignerSetTxsWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedSignerSetTxsWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedSignerSetTxsWindow
	}

	unslashedSignerSetTxs := k.GetUnSlashedSignerSetTxs(ctx, maxHeight)

	// unslashedSignerSetTxs are sorted by nonce in ASC order
	// Question: do we need to sort each time? See if this can be epoched
	for _, vs := range unslashedSignerSetTxs {
		sigMsgs := k.GetSignerSetTxSignatures(ctx, vs.Nonce)

		// SLASH BONDED VALIDTORS who didn't sign signer set tx
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		for _, val := range currentBondedSet {
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)

			//  Slash validator ONLY if he joined before signer set tx is created
			if exist && valSigningInfo.StartHeight < int64(vs.Nonce) {
				// Check if validator has signed signer set tx or not
				found := false
				for _, sigMsg := range sigMsgs {
					if sigMsg.EthAddress == k.GetEthAddressByValidator(ctx, val.GetOperator()) {
						found = true
						break
					}
				}
				// slash validators for not signing signer set txs
				if !found {
					cons, _ := val.GetConsAddr()
					k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionSignerSetTx)
					if !val.IsJailed() {
						k.StakingKeeper.Jail(ctx, cons)
					}

				}
			}
		}

		// SLASH UNBONDING VALIDATORS who didn't sign signer set tx
		blockTime := ctx.BlockTime().Add(k.StakingKeeper.GetParams(ctx).UnbondingTime)
		blockHeight := ctx.BlockHeight()
		unbondingValIterator := k.StakingKeeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
		defer unbondingValIterator.Close()

		// All unbonding validators
		for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
			unbondingValidators := k.GetUnbondingvalidators(unbondingValIterator.Value())

			for _, valAddr := range unbondingValidators.Addresses {
				addr, err := sdk.ValAddressFromBech32(valAddr)
				if err != nil {
					panic(err)
				}
				validator, _ := k.StakingKeeper.GetValidator(ctx, sdk.ValAddress(addr))
				valConsAddr, _ := validator.GetConsAddr()
				valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, valConsAddr)

				// Only slash validators who joined after signer set tx is created and they are unbonding and UNBOND_SLASHING_WINDOW didn't passed
				if exist && valSigningInfo.StartHeight < int64(vs.Nonce) && validator.IsUnbonding() && vs.Nonce < uint64(validator.UnbondingHeight)+params.SlashingSignerSetUnbondWindow {
					// Check if validator has signed signer set tx or not
					found := false
					for _, sigMsg := range sigMsgs {
						if sigMsg.EthAddress == k.GetEthAddressByValidator(ctx, validator.GetOperator()) {
							found = true
							break
						}
					}

					// slash validators for not signing signer set txs
					if !found {
						k.StakingKeeper.Slash(ctx, valConsAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionSignerSetTx)
						if !validator.IsJailed() {
							k.StakingKeeper.Jail(ctx, valConsAddr)
						}
					}
				}
			}
		}
		// then we set the latest slashed signer set tx  nonce
		k.SetLastSlashedSignerSetTxNonce(ctx, vs.Nonce)
	}
}

func BatchSlashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {

	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch tx that is >15hrs in blocks old
	maxHeight := uint64(0)

	// don't slash in the beginning before there aren't even SignedBatchTxsWindow blocks yet
	if uint64(ctx.BlockHeight()) > params.SignedBatchTxsWindow {
		maxHeight = uint64(ctx.BlockHeight()) - params.SignedBatchTxsWindow
	}

	unslashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)
	for _, batch := range unslashedBatches {

		// SLASH BONDED VALIDTORS who didn't sign batch tx
		currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		sigMsgs := k.GetBatchTxSignaturesByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
		for _, val := range currentBondedSet {
			// Don't slash validators who joined after batch is created
			consAddr, _ := val.GetConsAddr()
			valSigningInfo, exist := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if exist && valSigningInfo.StartHeight > int64(batch.Block) {
				continue
			}

			found := false
			for _, sigMsg := range sigMsgs {
				// TODO: double check this logic
				validator, _ := sdk.AccAddressFromBech32(sigMsg.Orchestrator)
				if k.GetOrchestratorValidator(ctx, validator).Equals(val.GetOperator()) {
					found = true
					break
				}
			}
			if !found {
				cons, _ := val.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatchTx)
				if !val.IsJailed() {
					k.StakingKeeper.Jail(ctx, cons)
				}
			}
		}
		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.Block)

	}
}

// TestingEndBlocker is a second endblocker function only imported in the Gravity codebase itself
// if you are a consuming Cosmos chain DO NOT IMPORT THIS, it simulates a chain using the
// contract call API
func TestingEndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// if this is nil we have not set our test contract call yet
	if k.GetContractCallTx(ctx, []byte("GravityTesting"), 0).ContractCallPayload == nil {
		// TODO this call isn't actually very useful for testing, since it always
		// throws, being just junk data that's expected. But it prevents us from checking
		// the full lifecycle of the call. We need to find some way for this to read data
		// and encode a simple testing call, probably to one of the already deployed ERC20
		// contracts so that we can get the full lifecycle.
		token := []*types.ERC20Token{{
			Contract: "0x7580bfe88dd3d07947908fae12d95872a260f2d8",
			Amount:   sdk.NewIntFromUint64(5000),
		}}
		_ = types.ContractCallTx{
			Transfers:           token,
			Fees:                token,
			ContractCallAddress: "0x510ab76899430424d209a6c9a5b9951fb8a6f47d",
			ContractCallPayload: []byte("fake bytes"),
			Timeout:             10000,
			InvalidationId:      []byte("GravityTesting"),
			InvalidationNonce:   1,
		}
		//k.SetContractCallTx(ctx, &call)
	}
}
