package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/althea-net/peggy/module/x/peggy/types"
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

	return &types.MsgSubmitConfirmResponse{}, nil
}

func (k msgServer) SubmitClaim(c context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {

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

	// Check if the denom is a peggy coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID, err := k.BuildOutgoingTXBatch(ctx, tokenContract, OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}

	// question: is this right? If i can delegate my voting power to a different key then this would fail each time i call it
	valaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := k.GetOrchestratorValidator(ctx, valaddr)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
		if sval == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
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
