package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

type ClaimHandler struct {
	Keeper
}

func NewClaimhandler(k Keeper) *ClaimHandler {
	return &ClaimHandler{Keeper: k}
}

func (ch *ClaimHandler) HandleClaim(claim *sdk.Any) (proto.Message, error) {
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

// DepositClaim handles MsgDepositClaim
// // TODO it is possible to submit an old msgDepositClaim (old defined as covering an event nonce that has already been
// // executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// // should not be a security risk as 'old' events can never execute but it does store spam in the chain.
// func (k msgServer) DepositClaim(c context.Context, msg *types.MsgDepositClaim) (*types.MsgDepositClaimResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)

// 	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, orch)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	// return an error if the validator isn't in the active set
// 	val := k.StakingKeeper.Validator(ctx, validator)
// 	if val == nil || !val.IsBonded() {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in active set")
// 	}

// 	any, err := codectypes.NewAnyWithValue(msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add the claim to the store
// 	_, err = k.Attest(ctx, msg, any)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(err, "create attestation")
// 	}

// 	// Emit the handle message event
// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			// TODO: maybe return something better here? is this the right string representation?
// 			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(msg.EventNonce, msg.ClaimHash()))),
// 		),
// 	)

// 	return &types.MsgDepositClaimResponse{}, nil
// }

// WithdrawClaim handles MsgWithdrawClaim
// TODO it is possible to submit an old msgWithdrawClaim (old defined as covering an event nonce that has already been
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
// func (k msgServer) WithdrawClaim(c context.Context, msg *types.MsgWithdrawClaim) (*types.MsgWithdrawClaimResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)

// 	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, orch)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	// return an error if the validator isn't in the active set
// 	val := k.StakingKeeper.Validator(ctx, validator)
// 	if val == nil || !val.IsBonded() {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
// 	}

// 	any, err := codectypes.NewAnyWithValue(msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add the claim to the store
// 	_, err = k.Attest(ctx, msg, any)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(err, "create attestation")
// 	}

// 	// Emit the handle message event
// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			// TODO: maybe return something better here? is this the right string representation?
// 			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(msg.EventNonce, msg.ClaimHash()))),
// 		),
// 	)

// 	return &types.MsgWithdrawClaimResponse{}, nil
// }

// ERC20Deployed handles MsgERC20Deployed
// func (k msgServer) ERC20DeployedClaim(c context.Context, msg *types.MsgERC20DeployedClaim) (*types.MsgERC20DeployedClaimResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)

// 	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, orch)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	// return an error if the validator isn't in the active set
// 	val := k.StakingKeeper.Validator(ctx, validator)
// 	if val == nil || !val.IsBonded() {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
// 	}

// 	any, err := codectypes.NewAnyWithValue(msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add the claim to the store
// 	_, err = k.Attest(ctx, msg, any)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(err, "create attestation")
// 	}

// 	// Emit the handle message event
// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			// TODO: maybe return something better here? is this the right string representation?
// 			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(msg.EventNonce, msg.ClaimHash()))),
// 		),
// 	)

// 	return &types.MsgERC20DeployedClaimResponse{}, nil
// }

// LogicCallExecutedClaim handles claims for executing a logic call on Ethereum
// func (k msgServer) LogicCallExecutedClaim(c context.Context, msg *types.MsgLogicCallExecutedClaim) (*types.MsgLogicCallExecutedClaimResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(c)

// 	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
// 	validator := k.GetOrchestratorValidator(ctx, orch)
// 	if validator == nil {
// 		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
// 		if sval == nil {
// 			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
// 		}
// 		validator = sval.GetOperator()
// 	}

// 	// return an error if the validator isn't in the active set
// 	val := k.StakingKeeper.Validator(ctx, validator)
// 	if val == nil || !val.IsBonded() {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
// 	}

// 	any, err := codectypes.NewAnyWithValue(msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add the claim to the store
// 	_, err = k.Attest(ctx, msg, any)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(err, "create attestation")
// 	}

// 	// Emit the handle message event
// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
// 			// TODO: maybe return something better here? is this the right string representation?
// 			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(msg.EventNonce, msg.ClaimHash()))),
// 		),
// 	)

// 	return &types.MsgLogicCallExecutedClaimResponse{}, nil
// }
