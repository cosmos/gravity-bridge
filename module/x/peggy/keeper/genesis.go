package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)
	// reset valsets in state
	for _, vs := range data.Valsets {
		// TODO: block height?
		k.StoreValsetUnsafe(ctx, vs)
	}

	// reset valset confirmations in state
	for _, conf := range data.ValsetConfirms {
		k.SetValsetConfirm(ctx, *conf)
	}

	// reset batches in state
	for _, batch := range data.Batches {
		// TODO: block height?
		k.StoreBatchUnsafe(ctx, batch)
	}

	// reset batch confirmations in state
	for _, conf := range data.BatchConfirms {
		k.SetBatchConfirm(ctx, &conf)
	}

	// reset attestations in state
	for _, att := range data.Attestations {
		// TODO: block height?
		k.SetAttestationUnsafe(ctx, &att)
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p            = k.GetParams(ctx)
		batches      = k.GetOutgoingTxBatches(ctx)
		valsets      = k.GetValsets(ctx)
		attmap       = k.GetAttestationMapping(ctx)
		vsconfs      = []*types.MsgValsetConfirm{}
		batchconfs   = []types.MsgConfirmBatch{}
		attestations = []types.Attestation{}
	)

	// export valset confirmations from state
	for _, vs := range valsets {
		// TODO: set height = 0?
		vsconfs = append(vsconfs, k.GetValsetConfirms(ctx, vs.Nonce)...)
	}

	// export batch confirmations from state
	for _, batch := range batches {
		// TODO: set height = 0?
		batchconfs = append(batchconfs, k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)...)
	}

	// export attestations from state
	for _, atts := range attmap {
		// TODO: set height = 0?
		attestations = append(attestations, atts...)
	}

	return types.GenesisState{
		Params:         &p,
		Valsets:        valsets,
		ValsetConfirms: vsconfs,
		Batches:        batches,
		BatchConfirms:  batchconfs,
		Attestations:   attestations,
	}
}
