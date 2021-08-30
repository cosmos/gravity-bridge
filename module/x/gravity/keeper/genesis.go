package keeper

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.setParams(ctx, *data.Params)

	// reset pool transactions in state
	for _, tx := range data.UnbatchedSendToEthereumTxs {
		k.setUnbatchedSendToEthereum(ctx, tx)
	}

	// reset ethereum event vote records in state
	for _, evr := range data.EthereumEventVoteRecords {
		event, err := types.UnpackEvent(evr.Event)
		if err != nil {
			panic("couldn't cast to event")
		}
		if err := event.Validate(); err != nil {
			panic("invalid event in genesis")
		}
		k.setEthereumEventVoteRecord(ctx, event.GetEventNonce(), event.Hash(), evr)
	}

	// reset last observed event nonce
	k.setLastObservedEventNonce(ctx, data.LastObservedEventNonce)

	// reset attestation state of all validators
	for _, eventVoteRecord := range data.EthereumEventVoteRecords {
		event, _ := types.UnpackEvent(eventVoteRecord.Event)
		for _, vote := range eventVoteRecord.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.getLastEventNonceByValidator(ctx, val)
			if event.GetEventNonce() > last {
				k.setLastEventNonceByValidator(ctx, val, event.GetEventNonce())
			}
		}
	}

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		if err := keys.ValidateBasic(); err != nil {
			panic("Invalid delegate key in Genesis!")
		}

		val, _ := sdk.ValAddressFromBech32(keys.ValidatorAddress)
		orch, _ := sdk.AccAddressFromBech32(keys.OrchestratorAddress)
		eth := common.HexToAddress(keys.EthereumAddress)

		// set the orchestrator address
		k.SetOrchestratorValidatorAddress(ctx, val, orch)
		// set the ethereum address
		k.setValidatorEthereumAddress(ctx, val, common.HexToAddress(keys.EthereumAddress))
		k.setEthereumOrchestratorAddress(ctx, eth, orch)
	}

	// populate state with cosmos originated denom-erc20 mapping
	for _, item := range data.Erc20ToDenoms {
		k.setCosmosOriginatedDenomToERC20(ctx, item.Denom, item.Erc20)
	}

	// reset outgoing txs in state
	for _, ota := range data.OutgoingTxs {
		otx, err := types.UnpackOutgoingTx(ota)
		if err != nil {
			panic("invalid outgoing tx any in genesis file")
		}
		k.SetOutgoingTx(ctx, otx)
	}

	// reset signatures in state
	for _, confa := range data.Confirmations {
		conf, err := types.UnpackConfirmation(confa)
		if err != nil {
			panic("invalid etheruem signature in genesis")
		}
		// TODO: not currently an easy way to get the validator address from the
		// etherum address here. once we implement the third index for keys
		// this will be easy.
		k.SetEthereumSignature(ctx, conf, sdk.ValAddress{})
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p                        = k.GetParams(ctx)
		outgoingTxs              []*cdctypes.Any
		ethereumTxConfirmations  []*cdctypes.Any
		attmap                   = k.GetEthereumEventVoteRecordMapping(ctx)
		ethereumEventVoteRecords []*types.EthereumEventVoteRecord
		delegates                = k.getDelegateKeys(ctx)
		lastobserved             = k.GetLastObservedEventNonce(ctx)
		erc20ToDenoms            []*types.ERC20ToDenom
		unbatchedTransfers       = k.getUnbatchedSendToEthereums(ctx)
	)

	// export ethereumEventVoteRecords from state
	for _, atts := range attmap {
		// TODO: set height = 0?
		ethereumEventVoteRecords = append(ethereumEventVoteRecords, atts...)
	}

	// export erc20 to denom relations
	k.iterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
		return false
	})

	// export signer set txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		sstx, _ := otx.(*types.SignerSetTx)
		k.iterateEthereumSignatures(ctx, sstx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.SignerSetTxConfirmation{sstx.Nonce, k.GetValidatorEthereumAddress(ctx, val).Hex(), sig})
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export batch txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.BatchTx)
		k.iterateEthereumSignatures(ctx, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.BatchTxConfirmation{btx.TokenContract, btx.BatchNonce, k.GetValidatorEthereumAddress(ctx, val).Hex(), sig})
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	// export contract call txs and sigs
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		ota, _ := types.PackOutgoingTx(otx)
		outgoingTxs = append(outgoingTxs, ota)
		btx, _ := otx.(*types.ContractCallTx)
		k.iterateEthereumSignatures(ctx, btx.GetStoreIndex(), func(val sdk.ValAddress, sig []byte) bool {
			siga, _ := types.PackConfirmation(&types.ContractCallTxConfirmation{btx.InvalidationScope, btx.InvalidationNonce, k.GetValidatorEthereumAddress(ctx, val).Hex(), sig})
			ethereumTxConfirmations = append(ethereumTxConfirmations, siga)
			return false
		})
		return false
	})

	return types.GenesisState{
		Params:                     &p,
		LastObservedEventNonce:     lastobserved,
		OutgoingTxs:                outgoingTxs,
		Confirmations:              ethereumTxConfirmations,
		EthereumEventVoteRecords:   ethereumEventVoteRecords,
		DelegateKeys:               delegates,
		Erc20ToDenoms:              erc20ToDenoms,
		UnbatchedSendToEthereumTxs: unbatchedTransfers,
	}
}
