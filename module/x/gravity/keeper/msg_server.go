package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetDelegateKey(c context.Context, msg *types.MsgDelegateKey) (*types.MsgDelegateKeyResponse, error) {
	// ensure that this passes validation
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	val, _ := sdk.ValAddressFromBech32(msg.Validator)
	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// ensure that the validator exists
	if k.Keeper.StakingKeeper.Validator(ctx, val) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	}

	// TODO consider impact of maliciously setting duplicate delegate
	// addresses since no signatures from the private keys of these addresses
	// are required for this message it could be sent in a hostile way.

	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, val, orch)
	// set the ethereum address
	k.SetEthAddress(ctx, val, msg.EthAddress)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeySetOperatorAddr, orch.String()),
		),
	)

	return &types.MsgDelegateKeyResponse{}, nil
}

func (k msgServer) SubmitConfirm(c context.Context, msg *types.MsgSubmitConfirm) (*types.MsgSubmitConfirmResponse, error) {

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

func (k msgServer) SubmitClaim(c context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {

	claim := msg.GetClaim()

	switch msg.ClaimType {
	case types.ClaimType_DEPOSIT:
		if err := k.depositClaim(c, claim); err != nil {
			return nil, err
		}
	case types.ClaimType_WITHDRAW:
		if err := k.withdrawClaim(c, claim); err != nil {
			return nil, err
		}
	case types.ClaimType_ERC20_DEPLOYED:
		if err := k.eRC20DeployedClaim(c, claim); err != nil {
			return nil, err
		}
	case types.ClaimType_LOGIC_CALL_EXECUTED:
		if err := k.logicCallExecutedClaim(c, claim); err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidClaim, claim.GetType().String())
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

// SendToEth handles MsgSendToEth
func (k msgServer) SendToEth(c context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	txID, err := k.AddToOutgoingPool(ctx, sender, msg.EthDest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)

	return &types.MsgSendToEthResponse{}, nil
}

// RequestBatch handles MsgRequestBatch
func (k msgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID, err := k.BuildOutgoingTXBatch(ctx, tokenContract, OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyBatchNonce, fmt.Sprint(batchID.BatchNonce)),
		),
	)

	return &types.MsgRequestBatchResponse{}, nil
}

func (k msgServer) CancelSendToEth(c context.Context, msg *types.MsgCancelSendToEth) (*types.MsgCancelSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(msg.TransactionId)),
		),
	)

	return &types.MsgCancelSendToEthResponse{}, nil
}
