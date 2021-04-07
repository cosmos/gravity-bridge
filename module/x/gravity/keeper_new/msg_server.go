package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

var _ types.MsgServer = &Keeper{}

// SetDelegateKey implements MsgServer.SetDelegateKey. The
func (k Keeper) SetDelegateKey(c context.Context, msg *types.MsgDelegateKey) (*types.MsgDelegateKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: address checked on msg validation
	validatorAddr, _ := sdk.ValAddressFromBech32(msg.Validator)
	orchestratorAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	ethereumAddr := common.HexToAddress(msg.EthAddress)

	// ensure that the validator exists
	if k.stakingKeeper.Validator(ctx, validatorAddr) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, validatorAddr.String())
	}

	// TODO consider impact of maliciously setting duplicate delegate
	// addresses since no signatures from the private keys of these addresses
	// are required for this message it could be sent in a hostile way.

	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, validatorAddr, orchestratorAddr)
	// set the ethereum address
	k.SetEthAddress(ctx, validatorAddr, ethereumAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeySetOperatorAddr, msg.Orchestrator),
		),
	)

	return &types.MsgDelegateKeyResponse{}, nil
}

func (k Keeper) SubmitConfirm(c context.Context, msg *types.MsgSubmitConfirm) (*types.MsgSubmitConfirmResponse, error) {

	confirm := msg.GetConfirm()

	switch msg.ConfirmType {
	case types.ConfirmType_CONFIRM_TYPE_BATCH:
	case types.ConfirmType_CONFIRM_TYPE_LOGIC:
	case types.ConfirmType_CONFIRM_TYPE_VALSET:
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidConfirm, confirm.GetType().String())
	}

	return &types.MsgSubmitConfirmResponse{}, nil
}

func (k Keeper) SubmitClaim(c context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	claim, err := types.UnpackClaim(msg.Claim)
	if err != nil {
		return nil, err
	}

	if err := k.HandleClaim(ctx, claim); err != nil {
		return nil, err
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

// RequestBatch handles MsgRequestBatch
func (k Keeper) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// question: is this right? If i can delegate my voting power to a different key then this would fail each time i call it
	valaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// FIXME: update logic

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID, err := k.BuildOutgoingTxBatch(ctx, tokenContract, OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}

	validator := k.GetOrchestratorValidator(ctx, valaddr)

	// a validator request can be sent from a delegate key or a validator
	// key directly
	if validator == nil {
		sval := k.stakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
		if sval == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
	}

	// TODO later make sure that Demon matches a list of tokens already
	// in the bridge to send

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyBatchNonce, fmt.Sprint(batchID.BatchNonce)),
		),
	)

	return &types.MsgRequestBatchResponse{}, nil
}

// SendToEth handles MsgSendToEth
func (k Keeper) SendToEth(c context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: errors checked on msg validation
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	ethreumAddr := common.HexToAddress(msg.EthDest)

	txID, err := k.AddToOutgoingPool(ctx, sender, ethreumAddr, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(txID)),
		),
	)

	return &types.MsgSendToEthResponse{}, nil
}

// CancelSendToEth
func (k Keeper) CancelSendToEth(c context.Context, msg *types.MsgCancelSendToEth) (*types.MsgCancelSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	if err := k.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(msg.TransactionId)),
		),
	)

	return &types.MsgCancelSendToEthResponse{}, nil
}
