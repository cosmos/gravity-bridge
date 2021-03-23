package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

type ConfirmHandler struct {
	Keeper
}

func NewConfirmhandler(k Keeper) *ConfirmHandler {
	return &ConfirmHandler{Keeper: k}
}

func (ch *ConfirmHandler) HandleClaim(claim *sdk.Any) (proto.Message, error) {
	switch cl := claim.(type) {
	case types.DepositClaim:
	case types.WithdrawClaim:
	case types.ERC20DeployedClaim:
	case types.LogicCallExecutedClaim:
	default:
		return nil, fmt.Errorf("claim type %T is not supported", cl)
	}
	return &types.MsgSubmitClaimResponse{}, nil
}

// // ConfirmBatch handles MsgConfirmBatch
// func (k msgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)

// 	// fetch the outgoing batch given the nonce
// 	batch := k.GetOutgoingTXBatch(ctx, msg.TokenContract, msg.Nonce)
// 	if batch == nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
// 	}

// 	peggyID := k.GetPeggyID(ctx)
// 	checkpoint, err := batch.GetCheckpoint(peggyID)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
// 	}

// 	sigBytes, err := hex.DecodeString(msg.Signature)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
// 	}

// 	valaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, valaddr)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	ethAddress := k.GetEthAddress(ctx, validator)
// 	if ethAddress == "" {
// 		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
// 	}

// 	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
// 	}

// 	// check if we already have this confirm
// 	if k.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, valaddr) != nil {
// 		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
// 	}
// 	key := k.SetBatchConfirm(ctx, msg)

// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			sdk.NewAttribute(types.AttributeKeyBatchConfirmKey, string(key)),
// 		),
// 	)

// 	return nil, nil
// }

// ConfirmLogicCall handles MsgConfirmLogicCall
// func (k msgServer) ConfirmLogicCall(c context.Context, msg *types.MsgConfirmLogicCall) (*types.MsgConfirmLogicCallResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)
// 	invalidationIdBytes, err := hex.DecodeString(msg.InvalidationId)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "invalidation id encoding")
// 	}

// 	// fetch the outgoing logic given the nonce
// 	logic := k.GetOutgoingLogicCall(ctx, invalidationIdBytes, msg.InvalidationNonce)
// 	if logic == nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find logic")
// 	}

// 	peggyID := k.GetPeggyID(ctx)
// 	checkpoint, err := logic.GetCheckpoint(peggyID)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
// 	}

// 	sigBytes, err := hex.DecodeString(msg.Signature)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
// 	}

// 	valaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, valaddr)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	ethAddress := k.GetEthAddress(ctx, validator)
// 	if ethAddress == "" {
// 		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
// 	}

// 	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
// 	}

// 	// check if we already have this confirm
// 	if k.GetLogicCallConfirm(ctx, invalidationIdBytes, msg.InvalidationNonce, valaddr) != nil {
// 		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
// 	}

// 	k.SetLogicCallConfirm(ctx, msg)

// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 		),
// 	)

// 	return nil, nil
// }

// // ValsetConfirm handles MsgValsetConfirm
// // TODO: check msgValsetConfirm to have an Orchestrator field instead of a Validator field
// func (k msgServer) ValsetConfirm(c context.Context, msg *types.MsgValsetConfirm) (*types.MsgValsetConfirmResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)
// 	valset := k.GetValset(ctx, msg.Nonce)
// 	if valset == nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
// 	}

// 	peggyID := k.GetPeggyID(ctx)
// 	checkpoint := valset.GetCheckpoint(peggyID)

// 	sigBytes, err := hex.DecodeString(msg.Signature)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
// 	}

// 	valaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, valaddr)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(valaddr))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	ethAddress := k.GetEthAddress(ctx, validator)
// 	if ethAddress == "" {
// 		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
// 	}

// 	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress); err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
// 	}

// 	// persist signature
// 	if k.GetValsetConfirm(ctx, msg.Nonce, valaddr) != nil {
// 		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
// 	}
// 	key := k.SetValsetConfirm(ctx, *msg)

// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			sdk.NewAttribute(types.AttributeKeyValsetConfirmKey, string(key)),
// 		),
// 	)

// 	return &types.MsgValsetConfirmResponse{}, nil
// }
