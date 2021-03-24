package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

// ConfirmBatch handles MsgConfirmBatch
func (k msgServer) confirmBatch(c context.Context, msg *types.ConfirmBatch) error {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := k.GetOutgoingTXBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	peggyID := k.GetPeggyID(ctx)
	checkpoint, err := batch.GetCheckpoint(peggyID)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
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
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// check if we already have this confirm
	if k.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}
	key := k.SetBatchConfirm(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.GetType().String()),
			sdk.NewAttribute(types.AttributeKeyBatchConfirmKey, string(key)),
		),
	)

	return nil
}

// ConfirmLogicCall handles MsgConfirmLogicCall
func (k msgServer) confirmLogicCall(c context.Context, msg types.Confirm) error {
	ctx := sdk.UnwrapSDKContext(c)
	invalidationIdBytes, err := hex.DecodeString(msg.InvalidationId)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "invalidation id encoding")
	}

	// fetch the outgoing logic given the nonce
	logic := k.GetOutgoingLogicCall(ctx, invalidationIdBytes, msg.InvalidationNonce)
	if logic == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find logic")
	}

	peggyID := k.GetPeggyID(ctx)
	checkpoint, err := logic.GetCheckpoint(peggyID)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(msg.GetSignature())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(msg.GetOrchestratorAddress())
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
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// check if we already have this confirm
	if k.GetLogicCallConfirm(ctx, invalidationIdBytes, msg.InvalidationNonce, valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}

	k.SetLogicCallConfirm(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.GetType().String()),
		),
	)

	return nil
}

// // ValsetConfirm handles MsgValsetConfirm
// TODO: check msgValsetConfirm to have an Orchestrator field instead of a Validator field
func (k msgServer) valsetConfirm(c context.Context, msg types.Confirm) error {
	ctx := sdk.UnwrapSDKContext(c)
	valset := k.GetValset(ctx, msg.GetNonce())
	if valset == nil {
		return sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
	}

	peggyID := k.GetPeggyID(ctx)
	checkpoint := valset.GetCheckpoint(peggyID)

	sigBytes, err := hex.DecodeString(msg.GetSignature())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	valaddr, _ := sdk.AccAddressFromBech32(msg.GetOrchestratorAddress())
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
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// persist signature
	if k.GetValsetConfirm(ctx, msg.GetNonce(), valaddr) != nil {
		return sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := 1
	// k.SetValsetConfirm(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.GetType().String()),
			sdk.NewAttribute(types.AttributeKeyValsetConfirmKey, string(key)),
		),
	)

	return nil
}
