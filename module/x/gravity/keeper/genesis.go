package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

// InitGenesis imports the gravity bridge state from a genesis.json and sets it
// to store
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
		k.SetLogicCallTx(ctx, tx.InvalidationID, tx.InvalidationNonce, tx.LogicCall)
	}

	for _, tx := range gs.TransferTxs {
		_ = k.SetTransferTx(ctx, tx)
	}

	for _, cc := range gs.Confirms {
		confirm, err := types.UnpackConfirm(cc.Confirm)
		if err != nil {
			panic(err)
		}

		k.SetConfirm(ctx, cc.ConfirmID, confirm)
	}

	// TODO: set last unbonding height?

	for _, att := range gs.Attestations {
		k.SetAttestation(ctx, att.AttestationID, att.Attestation)
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

// ExportGenesis exports the gravity bridge state into a genesis.json
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		BridgeID:          k.GetBridgeID(ctx),
		Params:            k.GetParams(ctx),
		LastObservedNonce: k.GetLastObservedEventNonce(ctx),
		SignerSets:        nil, // TODO:
		BatchTxs:          k.GetBatchTxs(ctx),
		LogicCallTxs:      k.GetIdentifiedLogicCalls(ctx),
		TransferTxs:       k.GetTransferTxs(ctx),
		Confirms:          nil, // k.GetIdentifiedConfirms(ctx),
		Attestations:      nil, // k.GetIdentifiedAttestations(ctx),
		DelegateKeys:      nil,
		Erc20ToDenoms:     k.GetERC20Denoms(ctx),
	}
}
