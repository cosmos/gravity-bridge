package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)
	// reset valsets in state
	for _, vs := range data.SignerSetTxs {
		// TODO: block height?
		k.StoreSignerSetTxUnsafe(ctx, vs)
	}

	// reset valset confirmations in state
	for _, conf := range data.SignerSetTxSignatures {
		k.SetEthereumSignature(ctx, conf, nil)
	}

	// reset batches in state
	for _, batch := range data.BatchTxs {
		// TODO: block height?
		k.StoreBatchUnsafe(ctx, batch)
	}

	// reset batch confirmations in state
	for _, conf := range data.BatchTxSignatures {
		conf := conf
		k.SetBatchConfirm(ctx, &conf)
	}

	// reset logic calls in state
	for _, call := range data.ContractCallTxs {
		k.SetContractCallTx(ctx, call)
	}

	// reset batch confirmations in state
	for _, conf := range data.ContractCallTxSignatures {
		conf := conf
		k.SetContractCallTxSignature(ctx, &conf)
	}

	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEthereumTxs {
		if err := k.setPoolEntry(ctx, tx); err != nil {
			panic(err)
		}
	}

	// reset attestations in state
	for _, eventVoteRecord := range data.EthereumEventVoteRecords {
		att := eventVoteRecord
		event, err := types.UnpackEvent(att.Event)
		if err != nil {
			panic("couldn't cast to event")
		}

		// TODO: block height?
		k.SetEthereumEventVoteRecord(ctx, event.GetNonce(), event.Hash(), &att)
	}
	k.setLastObservedEventNonce(ctx, data.LastObservedEventNonce)

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for _, eventVoteRecord := range data.EthereumEventVoteRecords {
		event, err := types.UnpackEvent(eventVoteRecord.Event)
		if err != nil {
			panic("couldn't cast to event")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all attestations
		// have been cleaned up GetLastEventNonceByValidator handles that case
		//
		// if we where to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all attestations have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range eventVoteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.GetLastEventNonceByValidator(ctx, val)
			if event.GetNonce() > last {
				k.setLastEventNonceByValidator(ctx, val, event.GetNonce())
			}
		}
	}

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		err := keys.ValidateBasic()
		if err != nil {
			panic("Invalid delegate key in Genesis!")
		}
		val, err := sdk.ValAddressFromBech32(keys.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		orch, err := sdk.AccAddressFromBech32(keys.OrchestratorAddress)
		if err != nil {
			panic(err)
		}

		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		k.SetEthAddress(ctx, val, keys.EthereumAddress)
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
		p                           = k.GetParams(ctx)
		contractCallTxs             = k.GetContractCallTxs(ctx)
		batchTxs                    = k.GetBatchTxes(ctx)
		updateSignerSetTxs          = k.GetSignerSetTxs(ctx)
		attmap                      = k.GetEthereumEventVoteRecordMapping(ctx)
		updateSignerSetTxSignatures []*types.SignerSetTxSignature
		batchTxSignatures           []types.BatchTxSignature
		contractCallTxSignatures    []types.ContractCallTxSignature
		ethereumEventVoteRecords    []types.EthereumEventVoteRecord
		delegates                   = k.GetDelegateKeys(ctx)
		lastobserved                = k.GetLastObservedEventNonce(ctx)
		erc20ToDenoms               []*types.ERC20ToDenom
		unbatchedTransfers          = k.GetPoolTransactions(ctx)
	)

	// export valset confirmations from state
	for _, updateSignerSetTx := range updateSignerSetTxs {
		// TODO: set height = 0?

		updateSignerSetTxSignatures = append(updateSignerSetTxSignatures, k.GetEthereumSignatures(ctx, updateSignerSetTx)...)
	}

	// export batch confirmations from state
	for _, batch := range batchTxs {
		// TODO: set height = 0?
		batchTxSignatures = append(batchTxSignatures,
			k.GetBatchTxSignatureByNonceAndTokenContract(ctx, batch.Nonce, batch.TokenContract)...)
	}

	// export logic call confirmations from state
	for _, call := range contractCallTxs {
		// TODO: set height = 0?
		contractCallTxSignatures = append(contractCallTxSignatures,
			k.GetContractCallTxSignatureByInvalidationIDAndNonce(ctx, call.InvalidationScope, call.InvalidationNonce)...)
	}

	// export ethereumEventVoteRecords from state
	for _, atts := range attmap {
		// TODO: set height = 0?
		ethereumEventVoteRecords = append(ethereumEventVoteRecords, atts...)
	}

	// export erc20 to denom relations
	k.IterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
		return false
	})

	return types.GenesisState{
		Params:                      &p,
		LastObservedEventNonce:      lastobserved,
		SignerSetTxs:          updateSignerSetTxs,
		SignerSetTxSignatures: updateSignerSetTxSignatures,
		BatchTxs:                    batchTxs,
		BatchTxSignatures:           batchTxSignatures,
		ContractCallTxs:             contractCallTxs,
		ContractCallTxSignatures:    contractCallTxSignatures,
		EthereumEventVoteRecords:    ethereumEventVoteRecords,
		DelegateKeys:                delegates,
		Erc20ToDenoms:               erc20ToDenoms,
		UnbatchedSendToEthereumTxs:  unbatchedTransfers,
	}
}
