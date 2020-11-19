package keeper

import (
	"context"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.QueryServer = Keeper{}

// Params satisfies the grpc.QueryServer interface
func (k Keeper) Params(goCtx context.Context, query *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{Params: k.GetParams(sdk.UnwrapSDKContext(goCtx))}, nil
}

// CurrentValset satisfies the grpc.QueryServer interface
func (k Keeper) CurrentValset(goCtx context.Context, query *types.QueryCurrentValsetRequest) (*types.QueryCurrentValsetResponse, error) {
	return &types.QueryCurrentValsetResponse{Valset: k.GetCurrentValset(sdk.UnwrapSDKContext(goCtx))}, nil
}

// ValsetRequest satisfies the grpc.QueryServer interface
func (k Keeper) ValsetRequest(goCtx context.Context, query *types.QueryValsetRequestRequest) (*types.QueryValsetRequestResponse, error) {
	return &types.QueryValsetRequestResponse{Valset: k.GetValsetRequest(sdk.UnwrapSDKContext(goCtx), query.Nonce)}, nil
}

// ValsetConfirm satisfies the grpc.QueryServer interface
func (k Keeper) ValsetConfirm(goCtx context.Context, query *types.QueryValsetConfirmRequest) (*types.QueryValsetConfirmResponse, error) {
	val, err := sdk.AccAddressFromBech32(query.Validator)
	if err != nil {
		return nil, err
	}
	return &types.QueryValsetConfirmResponse{Confirm: k.GetValsetConfirm(sdk.UnwrapSDKContext(goCtx), query.Nonce, val)}, nil
}

// ValsetConfirmsByNonce satisfies the grpc.QueryServer interface
func (k Keeper) ValsetConfirmsByNonce(goCtx context.Context, query *types.QueryValsetConfirmsByNonceRequest) (*types.QueryValsetConfirmsByNonceResponse, error) {
	return &types.QueryValsetConfirmsByNonceResponse{Confirms: k.AllValsetConfirmsByNonce(sdk.UnwrapSDKContext(goCtx), query.Nonce)}, nil
}

// LastValsetRequests satisfies the grpc.QueryServer interface
func (k Keeper) LastValsetRequests(goCtx context.Context, query *types.QueryLastValsetRequestsRequest) (*types.QueryLastValsetRequestsResponse, error) {
	return &types.QueryLastValsetRequestsResponse{Valsets: k.GetLastValsetRequests(sdk.UnwrapSDKContext(goCtx))}, nil
}

// LastPendingValsetRequestByAddr satisfies the grpc.QueryServer interface
func (k Keeper) LastPendingValsetRequestByAddr(goCtx context.Context, query *types.QueryLastPendingValsetRequestByAddrRequest) (*types.QueryLastPendingValsetRequestByAddrResponse, error) {
	valAddr, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryLastPendingValsetRequestByAddrResponse{Valset: k.GetLastPendingValsetRequest(sdk.UnwrapSDKContext(goCtx), valAddr)}, nil
}

// Batch staisfies the grpc.QueryServer interface
func (k Keeper) Batch(goCtx context.Context, query *types.QueryBatchRequest) (*types.QueryBatchResponse, error) {
	if err := types.ValidateEthAddress(query.Contract); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(goCtx), query.Contract, query.Nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	return &types.QueryBatchResponse{Batch: foundBatch}, nil
}

// BatchConfirms staisfies the grpc.QueryServer interface
func (k Keeper) BatchConfirms(goCtx context.Context, query *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := types.ValidateEthAddress(query.Contract); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	return &types.QueryBatchConfirmsResponse{Batches: k.GetAllBatchBatchConfirmsByNonceAndTokenContract(sdk.UnwrapSDKContext(goCtx), query.Nonce, query.Contract)}, nil
}

// LastPendingBatchRequestByAddr staisfies the grpc.QueryServer interface
func (k Keeper) LastPendingBatchRequestByAddr(goCtx context.Context, query *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	valAddr, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: k.GetLastPendingBatch(sdk.UnwrapSDKContext(goCtx), valAddr)}, nil
}

// OutgoingTxBatches staisfies the grpc.QueryServer interface
func (k Keeper) OutgoingTxBatches(goCtx context.Context, query *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(goCtx), func(_ []byte, batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	if len(batches) == 0 {
		return nil, nil
	}
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches}, nil
}
