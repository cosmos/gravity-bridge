package keeper

import (
	"context"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params queries the params of the peggy module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.QueryParamsResponse{Params: params}, nil

}

// CurrentValset queries the CurrentValset of the peggy module
func (k Keeper) CurrentValset(c context.Context, req *types.QueryCurrentValsetRequest) (*types.QueryCurrentValsetResponse, error) {
	return &types.QueryCurrentValsetResponse{Valset: k.GetCurrentValset(sdk.UnwrapSDKContext(c))}, nil
}

// ValsetRequest queries the ValsetRequest of the peggy module
func (k Keeper) ValsetRequest(c context.Context, req *types.QueryValsetRequestRequest) (*types.QueryValsetRequestResponse, error) {
	return &types.QueryValsetRequestResponse{Valset: k.GetValsetRequest(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

// ValsetConfirm queries the ValsetConfirm of the peggy module
func (k Keeper) ValsetConfirm(c context.Context, req *types.QueryValsetConfirmRequest) (*types.QueryValsetConfirmResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryValsetConfirmResponse{Confirm: k.GetValsetConfirm(sdk.UnwrapSDKContext(c), req.Nonce, addr)}, nil
}

// ValsetConfirmsByNonce queries the ValsetConfirmsByNonce of the peggy module
func (k Keeper) ValsetConfirmsByNonce(c context.Context, req *types.QueryValsetConfirmsByNonceRequest) (*types.QueryValsetConfirmsByNonceResponse, error) {
	return nil, nil
}

// LastValsetRequests queries the LastValsetRequests of the peggy module
func (k Keeper) LastValsetRequests(c context.Context, req *types.QueryLastValsetRequestsRequest) (*types.QueryLastValsetRequestsResponse, error) {
	return nil, nil
}

// LastPendingValsetRequestByAddr queries the LastPendingValsetRequestByAddr of the peggy module
func (k Keeper) LastPendingValsetRequestByAddr(c context.Context, req *types.QueryLastPendingValsetRequestByAddrRequest) (*types.QueryLastPendingValsetRequestByAddrResponse, error) {
	return nil, nil
}

// LastPendingBatchRequestByAddr queries the LastPendingBatchRequestByAddr of the peggy module
func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	return nil, nil
}

// OutgoingTxBatches queries the OutgoingTxBatches of the peggy module
func (k Keeper) OutgoingTxBatches(c context.Context, req *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	return nil, nil
}

// BatchRequestByNonce queries the BatchRequestByNonce of the peggy module
func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	return nil, nil
}
