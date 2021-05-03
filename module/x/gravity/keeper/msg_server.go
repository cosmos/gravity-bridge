package keeper

import (
	"context"
	"strconv"

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

	k.SetOrchestratorValidator(ctx, validatorAddr, orchestratorAddr)
	k.SetEthAddress(ctx, validatorAddr, ethereumAddr)

	k.Logger(ctx).Debug(
		"orchestrator key delegated",
		"validator-address", msg.ValidatorAddress,
		"orchestrator-address", msg.OrchestratorAddress,
		"ethereum-address", msg.EthAddress,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegateKey,
			sdk.NewAttribute(types.AttributeKeyOrchestratorAddr, msg.OrchestratorAddress),
			sdk.NewAttribute(types.AttributeKeyEthereumAddr, msg.EthAddress),
			sdk.NewAttribute(types.AttributeKeyValidatorAddr, msg.ValidatorAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgDelegateKeyResponse{}, nil
}

func (k Keeper) SubmitEvent(c context.Context, msg *types.MsgSubmitEvent) (*types.MsgSubmitEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	event, _ := types.UnpackEvent(msg.Event)

	val, err := k.GetValidatorFromSigner(ctx, signer)
	if err != nil {
		return nil, err
	}

	eventID, err := k.AttestEvent(ctx, event, val)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitEvent,
			sdk.NewAttribute(types.AttributeKeyEventID, eventID.String()),
			sdk.NewAttribute(types.AttributeKeyEventType, event.GetType()),
			sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(event.GetNonce(), 64)),
			sdk.NewAttribute(types.AttributeKeyOrchestratorAddr, val.GetOperator().String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgSubmitEventResponse{EventID: eventID}, nil
}

func (k Keeper) SubmitConfirm(c context.Context, msg *types.MsgSubmitConfirm) (*types.MsgSubmitConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	confirm, _ := types.UnpackConfirm(msg.Confirm)

	val, err := k.GetValidatorFromSigner(ctx, signer)
	if err != nil {
		return nil, err
	}

	ethAddress := k.GetEthAddress(ctx, val.GetOperator())
	if (ethAddress == common.Address{}) {
		return nil, sdkerrors.Wrap(types.ErrValidatorEthAddressNotFound, val.GetOperator().String())
	}

	confirmID, err := k.ProcessConfirm(ctx, confirm, ethAddress)
	if err != nil {
		return nil, err
	}

	k.Logger(ctx).Debug(
		"confirm",
		"id", confirmID.String(),
		"type", confirm.GetType(),
		"ethereum-address", ethAddress.String(),
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
		sdk.NewEvent(
			types.EventTypeConfirm,
			sdk.NewAttribute(types.AttributeKeyConfirmID, confirmID.String()),
			sdk.NewAttribute(types.AttributeKeyConfirmType, confirm.GetType()),
			sdk.NewAttribute(types.AttributeKeyEthereumAddr, ethAddress.String()),
		),
	})

	return &types.MsgSubmitConfirmResponse{ConfirmID: confirmID}, nil
}

// RequestBatch handles MsgRequestBatch
func (k Keeper) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// orchestratorAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

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

	// TODO: later make sure that Demon matches a list of tokens already
	// in the bridge to send

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgRequestBatchResponse{}, nil
}

// Transfer handles MsgTransfer
func (k Keeper) Transfer(c context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: errors checked on msg validation
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	ethereumAddr := common.HexToAddress(msg.EthRecipient)

	txID, err := k.AddTransferToOutgoingPool(ctx, sender, ethereumAddr, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return &types.MsgTransferResponse{
		TxID: txID,
	}, nil
}

// CancelTransfer
func (k Keeper) CancelTransfer(c context.Context, msg *types.MsgCancelTransfer) (*types.MsgCancelTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	if err := k.RemoveFromOutgoingPoolAndRefund(ctx, msg.TxID, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	return &types.MsgCancelTransferResponse{}, nil
}

func (k Keeper) GetValidatorFromSigner(ctx sdk.Context, signer sdk.AccAddress) (stakingtypes.ValidatorI, error) {
	validatorAddr := k.GetOrchestratorValidator(ctx, signer)
	if validatorAddr == nil {
		validatorAddr = sdk.ValAddress(signer)
	}

	// return an error if the validator isn't in the active set
	validator := k.stakingKeeper.Validator(ctx, validatorAddr)
	if validator == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, validatorAddr.String())
	} else if !validator.IsBonded() {
		return nil, sdkerrors.Wrapf(types.ErrValidatorNotBonded, "validator %s not in active set", validatorAddr)
	}
	return validator, nil
}
