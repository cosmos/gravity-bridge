package keeper

import (
	"context"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) ValsetConfirm(goCtx context.Context, msg *types.MsgValsetConfirm) (*types.MsgValsetConfirmResponse, error) {
	return &types.MsgValsetConfirmResponse{}, nil
}
func (m msgServer) ValsetRequest(goCtx context.Context, msg *types.MsgValsetRequest) (*types.MsgValsetRequestResponse, error) {
	return &types.MsgValsetRequestResponse{}, nil
}
func (m msgServer) SetEthAddress(goCtx context.Context, msg *types.MsgSetEthAddress) (*types.MsgSetEthAddressResponse, error) {
	return &types.MsgSetEthAddressResponse{}, nil
}
func (m msgServer) SendToEth(goCtx context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	return &types.MsgSendToEthResponse{}, nil
}
func (m msgServer) RequestBatch(goCtx context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	return &types.MsgRequestBatchResponse{}, nil
}
func (m msgServer) ConfirmBatch(goCtx context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	return &types.MsgConfirmBatchResponse{}, nil
}
func (m msgServer) CreateEthereumClaims(goCtx context.Context, msg *types.MsgCreateEthereumClaims) (*types.MsgCreateEthereumClaimsResponse, error) {
	return &types.MsgCreateEthereumClaimsResponse{}, nil
}
func (m msgServer) BridgeSignatureSubmission(goCtx context.Context, msg *types.MsgBridgeSignatureSubmission) (*types.MsgBridgeSignatureSubmissionResponse, error) {
	return &types.MsgBridgeSignatureSubmissionResponse{}, nil
}
