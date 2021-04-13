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
	validatorAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	orchestratorAddr, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
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
			sdk.NewAttribute(types.AttributeKeySetOperatorAddr, msg.OrchestratorAddress),
		),
	)

	return &types.MsgDelegateKeyResponse{}, nil
}

func (k Keeper) SubmitConfirm(c context.Context, msg *types.MsgSubmitConfirm) (*types.MsgSubmitConfirmResponse, error) {

	// TODO:

	return &types.MsgSubmitConfirmResponse{}, nil
}

func (k Keeper) SubmitEvent(c context.Context, msg *types.MsgSubmitEvent) (*types.MsgSubmitEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	orchestratorAddr, _ := sdk.AccAddressFromBech32(msg.Signer)

	event, err := types.UnpackEvent(msg.Event)
	if err != nil {
		return nil, err
	}

	if err := k.HandleEthEvent(ctx, event, orchestratorAddr); err != nil {
		return nil, err
	}

	return &types.MsgSubmitEventResponse{}, nil
}

// RequestBatch handles MsgRequestBatch
func (k Keeper) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// orchestratorAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// FIXME: update logic

	var (
		tokenContract common.Address
		found         bool
	)

	if types.IsEthereumERC20Token(msg.Denom) {
		tokenContractHex := types.GravityDenomToERC20Contract(msg.Denom)
		tokenContract = common.HexToAddress(tokenContractHex)
	} else {
		// get contract from store
		tokenContract, found = k.GetERC20ContractFromCoinDenom(ctx, msg.Denom)
		if !found {
			// TODO: what if there is no corresponding contract? will it be "generated" on ethereum
			// upon receiving?
			// FIXME: this will fail if the cosmos tokens are relayed for the first time and they don't have a counterparty contract
			// Fede: Also there's the question of how do we handle IBC denominations from a security perspective. Do we assign them the same
			// contract? My guess is that each new contract assigned to a cosmos coin should be approved by governance
			return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "denom %s", msg.Denom)
		}
	}

	_, err := k.CreateBatchTx(ctx, tokenContract)
	if err != nil {
		return nil, err
	}

	// TODO later make sure that Demon matches a list of tokens already
	// in the bridge to send

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
		),
	)

	return &types.MsgRequestBatchResponse{}, nil
}

// SendToEth handles MsgSendToEth
func (k Keeper) SendToEth(c context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: errors checked on msg validation
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	ethereumAddr := common.HexToAddress(msg.EthRecipient)

	txID, err := k.AddToOutgoingPool(ctx, sender, ethereumAddr, msg.Amount, msg.BridgeFee)
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
