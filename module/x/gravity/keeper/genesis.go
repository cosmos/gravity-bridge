package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)
	// reset signer sets in state
	for _, vs := range data.SignerSetTxs {
		// TODO: block height?
		k.StoreSignerSetTxUnsafe(ctx, vs)
	}

	// reset signer set tx signatures in state
	for _, sigMsg := range data.SignerSetTxSignatures {
		k.SetSignerSetTxSignature(ctx, *sigMsg)
	}

	// reset batches in state
	for _, batch := range data.BatchTxs {
		// TODO: block height?
		k.StoreBatchUnsafe(ctx, batch)
	}

	// reset batch tx signatures in state
	for _, sigMsg := range data.BatchTxSignatures {
		sigMsg := sigMsg
		k.SetBatchTxSignature(ctx, &sigMsg)
	}

	// reset logic calls in state
	for _, call := range data.ContractCallTxs {
		k.SetContractCallTx(ctx, call)
	}

	// reset contract call tx signatures in state
	for _, sigMsg := range data.ContractCallTxSignatures {
		sigMsg := sigMsg
		k.SetContractCallTxSignature(ctx, &sigMsg)
	}

	// reset pool transactions in state
	for _, tx := range data.UnbatchedTransfers {
		if err := k.setPoolEntry(ctx, tx); err != nil {
			panic(err)
		}
	}

	// reset ethereumEventVoteRecords in state
	for _, voteRecord := range data.EthereumEventVoteRecords {
		voteRecord := voteRecord
		event, err := k.UnpackEthereumEventVoteRecordEvent(&voteRecord)
		if err != nil {
			panic("couldn't cast to event")
		}

		// TODO: block height?
		k.SetEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.EventHash(), &voteRecord)
	}
	k.setLastObservedEventNonce(ctx, data.LastObservedNonce)

	// reset ethereumEventVoteRecord state of specific validators
	// this must be done after the above to be correct
	for _, voteRecord := range data.EthereumEventVoteRecords {
		voteRecord := voteRecord
		event, err := k.UnpackEthereumEventVoteRecordEvent(&voteRecord)
		if err != nil {
			panic("couldn't cast to event")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all ethereumEventVoteRecords
		// have been cleaned up GetLastEventNonceByValidator handles that case
		//
		// if we where to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all ethereumEventVoteRecords have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range voteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.GetLastEventNonceByValidator(ctx, val)
			if event.GetEventNonce() > last {
				k.setLastEventNonceByValidator(ctx, val, event.GetEventNonce())
			}
		}
	}

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		err := keys.ValidateBasic()
		if err != nil {
			panic("Invalid delegate key in Genesis!")
		}
		val, err := sdk.ValAddressFromBech32(keys.Validator)
		if err != nil {
			panic(err)
		}

		orch, err := sdk.AccAddressFromBech32(keys.Orchestrator)
		if err != nil {
			panic(err)
		}

		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		k.SetEthAddressForValidator(ctx, val, keys.EthAddress)
	}

	// populate state with cosmos originated denom-erc20 mapping
	for _, item := range data.Erc20ToDenoms {
		k.setCosmosOriginatedDenomToERC20(ctx, item.Denom, item.Erc20)
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p                        = k.GetParams(ctx)
		calls                    = k.GetContractCallTxs(ctx)
		batches                  = k.GetBatchTxs(ctx)
		valsets                  = k.GetSignerSetTxs(ctx)
		voteRecordMap            = k.GetEthereumEventVoteRecordMapping(ctx)
		signerSetSigs            = []*types.MsgSignerSetTxSignature{}
		batchSigs                = []types.MsgBatchTxSignature{}
		contractCallSigs         = []types.MsgContractCallTxSignature{}
		ethereumEventVoteRecords = []types.EthereumEventVoteRecord{}
		delegates                = k.GetDelegateKeys(ctx)
		lastobserved             = k.GetLastObservedEventNonce(ctx)
		erc20ToDenoms            = []*types.ERC20ToDenom{}
		unbatchedTransfers       = k.GetPoolTransactions(ctx)
	)

	// export signer set tx signatures from state
	for _, vs := range valsets {
		// TODO: set height = 0?
		signerSetSigs = append(signerSetSigs, k.GetSignerSetTxSignatures(ctx, vs.Nonce)...)
	}

	// export batch tx signatures from state
	for _, batch := range batches {
		// TODO: set height = 0?
		batchSigs = append(batchSigs,
			k.GetBatchTxSignaturesByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)...)
	}

	// export logic call tx signatures from state
	for _, call := range calls {
		// TODO: set height = 0?
		contractCallSigs = append(contractCallSigs,
			k.GetContractCallTxSignaturesByInvalidationIDAndNonce(ctx, call.InvalidationId, call.InvalidationNonce)...)
	}

	// export ethereumEventVoteRecords from state
	for _, voteRecords := range voteRecordMap {
		// TODO: set height = 0?
		ethereumEventVoteRecords = append(ethereumEventVoteRecords, voteRecords...)
	}

	// export erc20 to denom relations
	k.IterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
		return false
	})

	return types.GenesisState{
		Params:                   &p,
		LastObservedNonce:        lastobserved,
		SignerSetTxs:             valsets,
		SignerSetTxSignatures:    signerSetSigs,
		BatchTxs:                 batches,
		BatchTxSignatures:        batchSigs,
		ContractCallTxs:          calls,
		ContractCallTxSignatures: contractCallSigs,
		EthereumEventVoteRecords: ethereumEventVoteRecords,
		DelegateKeys:             delegates,
		Erc20ToDenoms:            erc20ToDenoms,
		UnbatchedTransfers:       unbatchedTransfers,
	}
}
