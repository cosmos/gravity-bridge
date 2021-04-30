package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

// TODO:
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetBridgeID(ctx, gs.BridgeID)
	k.SetParams(ctx, gs.Params)
	k.SetLastObservedEventNonce(ctx, gs.LastObservedNonce)

	for _, ss := range gs.SignerSets {
		// store signer set and latest height
		k.StoreEthSignerSet(ctx, ss)
	}

	for _, tx := range gs.BatchTxs {
		k.SetBatchTx(ctx, tx)
	}

	for _, tx := range gs.LogicCallTxs {
		// FIXME: id and nonce
		k.SetLogicCallTx(ctx, nil, 0, tx)
	}

	for _, tx := range gs.TransferTxs {
		_ = k.SetTransferTx(ctx, tx)
	}

	// for _, confirm := range gs.SignerSetConfirms {

	// }

	// for _, confirm := range gs.BatchConfirms {

	// }

	// for _, confirm := range gs.LogicCallConfirms {

	// }

	// TODO: set last unbonding height?

	// TODO: export attestations as map<hash, Attestation>?
	for _, attestation := range gs.Attestations {
		// FIXME: attestation id
		k.SetAttestation(ctx, nil, &attestation)
	}

	// TODO: events

	for _, delegation := range gs.DelegateKeys {
		validatorAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		orchestratorAddr, _ := sdk.AccAddressFromBech32(delegation.ValidatorAddress)
		ethereumAddr := common.HexToAddress(delegation.EthAddress)

		// set the orchestrator and ethereum addresses
		k.SetOrchestratorValidator(ctx, validatorAddr, orchestratorAddr)
		k.SetEthAddress(ctx, validatorAddr, ethereumAddr)
	}

	for _, e := range gs.Erc20ToDenoms {
		tokenAddress := common.HexToAddress(e.Erc20Address)
		k.setERC20DenomMap(ctx, e.Denom, tokenAddress)
	}
}

// TODO:
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		BridgeID:          k.GetBridgeID(ctx),
		Params:            k.GetParams(ctx),
		LastObservedNonce: k.GetLastObservedEventNonce(ctx),
		SignerSets:        nil, // TODO:
		BatchTxs:          k.GetBatchTxs(ctx),
		LogicCallTxs:      k.GetOutgoingLogicCalls(ctx),
		TransferTxs:       k.GetTransferTxs(ctx),
		SignerSetConfirms: nil,
		BatchConfirms:     nil,
		LogicCallConfirms: nil,
		Attestations:      nil,
		DelegateKeys:      nil,
		Erc20ToDenoms:     k.GetERC20Denoms(ctx),
	}
}
