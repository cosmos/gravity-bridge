package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// ConfirmBatch handles ConfirmBatch
func (k Keeper) confirmBatch(c context.Context, confirm types.Confirm) error {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := k.GetOutgoingTXBatch(ctx, confirm.GetTokenContract(), confirm.GetNonce())
	if batch == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	peggyID := k.GetGravityID(ctx)
	checkpoint, err := batch.GetCheckpoint(peggyID)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(confirm.GetSignature())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(confirm.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, valaddr)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	ethAddress := k.GetEthAddress(ctx, validator)
	if ethAddress == "" {
		return sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}

	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), confirm.GetSignature()))
	}

	// check if we already have this confirm
	if k.GetBatchConfirm(ctx, confirm.GetNonce(), confirm.GetTokenContract(), valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}

	msg, ok := confirm.(*types.ConfirmBatch)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	key := k.SetBatchConfirm(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, confirm.GetType().String()),
			sdk.NewAttribute(types.AttributeKeyBatchConfirmKey, string(key)),
		),
	)

	return nil
}

// ConfirmLogicCall handles MsgConfirmLogicCall
func (k Keeper) confirmLogicCall(c context.Context, confirm types.Confirm) error {
	ctx := sdk.UnwrapSDKContext(c)
	invalidationIdBytes, err := hex.DecodeString(confirm.GetInvalidationId())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "invalidation id encoding")
	}

	// fetch the outgoing logic given the nonce
	logic := k.GetOutgoingLogicCall(ctx, invalidationIdBytes, confirm.GetInvalidationNonce())
	if logic == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find logic")
	}

	peggyID := k.GetGravityID(ctx)
	checkpoint, err := logic.GetCheckpoint(peggyID)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(confirm.GetSignature())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(confirm.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, valaddr)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	ethAddress := k.GetEthAddress(ctx, validator)
	if ethAddress == "" {
		return sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}

	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), confirm.GetSignature()))
	}

	// check if we already have this confirm
	if k.GetLogicCallConfirm(ctx, invalidationIdBytes, confirm.GetInvalidationNonce(), valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}

	msg, ok := confirm.(*types.ConfirmLogicCall)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	k.SetLogicCallConfirm(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, confirm.GetType().String()),
		),
	)

	return nil
}

// // ValsetConfirm handles MsgValsetConfirm
// TODO: check msgValsetConfirm to have an Orchestrator field instead of a Validator field
func (k Keeper) valsetConfirm(c context.Context, confirm types.Confirm) error {
	ctx := sdk.UnwrapSDKContext(c)
	valset := k.GetValset(ctx, confirm.GetNonce())
	if valset == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
	}

	peggyID := k.GetGravityID(ctx)
	checkpoint := valset.GetCheckpoint(peggyID)

	sigBytes, err := hex.DecodeString(confirm.GetSignature())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(confirm.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, valaddr)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	ethAddress := k.GetEthAddress(ctx, validator)
	if ethAddress == "" {
		return sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}

	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress); err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), confirm.GetSignature()))
	}

	// persist signature
	if k.GetValsetConfirm(ctx, confirm.GetNonce(), valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}

	msg, ok := confirm.(*types.ValsetConfirm)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", confirm)
	}

	key := k.SetValsetConfirm(ctx, *msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, confirm.GetType().String()),
			sdk.NewAttribute(types.AttributeKeyValsetConfirmKey, string(key)),
		),
	)

	return nil
}
