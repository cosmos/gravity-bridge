package keeper

import (
	"context"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}
	return &types.QueryValsetConfirmResponse{Confirm: k.GetValsetConfirm(sdk.UnwrapSDKContext(c), req.Nonce, addr)}, nil
}

// ValsetConfirmsByNonce queries the ValsetConfirmsByNonce of the peggy module
func (k Keeper) ValsetConfirmsByNonce(c context.Context, req *types.QueryValsetConfirmsByNonceRequest) (*types.QueryValsetConfirmsByNonceResponse, error) {
	var confirms []*types.MsgValsetConfirm
	k.IterateValsetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(_ []byte, c types.MsgValsetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryValsetConfirmsByNonceResponse{Confirms: confirms}, nil
}

// LastValsetRequests queries the LastValsetRequests of the peggy module
func (k Keeper) LastValsetRequests(c context.Context, req *types.QueryLastValsetRequestsRequest) (*types.QueryLastValsetRequestsResponse, error) {
	var counter int
	var valReq []*types.Valset
	k.IterateValsetRequest(sdk.UnwrapSDKContext(c), func(_ []byte, val *types.Valset) bool {
		valReq = append(valReq, val)
		counter++
		return counter >= maxValsetRequestsReturned
	})
	return &types.QueryLastValsetRequestsResponse{Valsets: valReq}, nil
}

// LastPendingValsetRequestByAddr queries the LastPendingValsetRequestByAddr of the peggy module
func (k Keeper) LastPendingValsetRequestByAddr(c context.Context, req *types.QueryLastPendingValsetRequestByAddrRequest) (*types.QueryLastPendingValsetRequestByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingValsetReq *types.Valset
	k.IterateValsetRequest(sdk.UnwrapSDKContext(c), func(_ []byte, val *types.Valset) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := k.GetValsetConfirm(sdk.UnwrapSDKContext(c), val.Nonce, addr) != nil
		// if this valset has NOT been signed by operatorAddr, store it in pendingValsetReq
		// and exit the loop
		if !foundConfirm {
			pendingValsetReq = val
			return true
		}
		// return false to continue the loop
		return false
	})
	return &types.QueryLastPendingValsetRequestByAddrResponse{Valset: pendingValsetReq}, nil
}

// LastPendingBatchRequestByAddr queries the LastPendingBatchRequestByAddr of the peggy module
func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.OutgoingTxBatch) bool {
		foundConfirm := k.GetBatchConfirm(sdk.UnwrapSDKContext(c), batch.BatchNonce, batch.TokenContract, addr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})

	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: pendingBatchReq}, nil
}

// OutgoingTxBatches queries the OutgoingTxBatches of the peggy module
func (k Keeper) OutgoingTxBatches(c context.Context, req *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches}, nil
}

// BatchRequestByNonce queries the BatchRequestByNonce of the peggy module
func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateEthAddress(req.ContractAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(c), req.ContractAddress, req.Nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}
