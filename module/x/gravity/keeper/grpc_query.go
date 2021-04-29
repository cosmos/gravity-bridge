package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the params of the gravity module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// GetDelegateKeyByValidator
func (k Keeper) GetDelegateKeyByValidator(c context.Context, req *types.QueryDelegateKeysByValidatorAddress) (*types.QueryDelegateKeysByValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.New("invalid address"), val.String())
	}
	eth := k.GetEthAddress(ctx, val)
	orc := k.GetEthOrchAddress(ctx, eth)
	return &types.QueryDelegateKeysByValidatorAddressResponse{EthAddress: eth.Hex(), OrchestratorAddress: orc.String()}, nil
}

// GetDelegateKeyByEth
func (k Keeper) GetDelegateKeyByEth(c context.Context, req *types.QueryDelegateKeysByEthAddress) (*types.QueryDelegateKeysByEthAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := types.ValidateEthAddress(req.EthAddress); err != nil {
		return nil, sdkerrors.Wrap(errors.New("invalid address"), req.EthAddress)
	}
	orc := k.GetEthOrchAddress(ctx, common.HexToAddress(req.EthAddress))
	val := k.GetOrchestratorValidator(ctx, orc)

	return &types.QueryDelegateKeysByEthAddressResponse{OrchestratorAddress: orc.String(), ValidatorAddress: val.String()}, nil
}

// GetDelegateKeyByOrchestrator
func (k Keeper) GetDelegateKeyByOrchestrator(c context.Context, req *types.QueryDelegateKeysByOrchestratorAddress) (*types.QueryDelegateKeysByOrchestratorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orc, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.New("invalid address"), req.OrchestratorAddress)
	}
	val := k.GetOrchestratorValidator(ctx, orc)
	eth := k.GetEthAddress(ctx, val)

	return &types.QueryDelegateKeysByOrchestratorAddressResponse{EthAddress: eth.Hex(), ValidatorAddress: val.String()}, nil
}

func (k Keeper) QuerySignerSetConfirmsRequest(c context.Context, req *types.QuerySignerSetConfirmsRequest) (*types.QueryConfirmsResponse, error) {
	// TODO: more validation here
	ctx := sdk.UnwrapSDKContext(c)
	// if the validator isn't specified return all the signatures and their associated validator
	if req.Validator == "" {
		return &types.QueryConfirmsResponse{Confirms: k.GetSignerSetConfirms(ctx, req.Nonce)}, nil
	}
	val, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrAddressInvalid, req.Validator)
	}
	return &types.QueryConfirmsResponse{Confirms: map[string][]byte{val.String(): k.GetSignerSetConfirm(ctx, req.Nonce, val)}}, nil
}

func (k Keeper) QueryBatchConfirmsRequest(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryConfirmsResponse, error) {
	// TODO: more validation here
	ctx := sdk.UnwrapSDKContext(c)
	// if the validator isn't specified return all the signatures and their associated validator
	if req.Validator == "" {
		return &types.QueryConfirmsResponse{Confirms: k.GetBatchConfirms(ctx, req.Nonce, common.HexToAddress(req.TokenAddress))}, nil
	}
	val, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrAddressInvalid, req.Validator)
	}
	return &types.QueryConfirmsResponse{Confirms: map[string][]byte{val.String(): k.GetConfirmBatch(ctx, req.Nonce, val, common.HexToAddress(req.TokenAddress))}}, nil
}

func (k Keeper) QueryLogicCallConfirmsRequest(c context.Context, req *types.QueryLogicCallConfirmsRequest) (*types.QueryConfirmsResponse, error) {
	// TODO: more validation here
	ctx := sdk.UnwrapSDKContext(c)
	// if the validator isn't specified return all the signatures and their associated validator
	if req.Validator == "" {
		return &types.QueryConfirmsResponse{Confirms: k.GetLogicCallConfirms(ctx, req.InvalidationId, req.InvalidationNonce)}, nil

	}
	val, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrAddressInvalid, req.Validator)
	}
	return &types.QueryConfirmsResponse{Confirms: map[string][]byte{val.String(): k.GetConfirmLogicCall(ctx, req.InvalidationId, req.InvalidationNonce, val)}}, nil

}
