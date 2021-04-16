package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// DepositClaim handles MsgDepositClaim
// // TODO it is possible to submit an old msgDepositClaim (old defined as covering an event nonce that has already been
// // executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// // should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (k msgServer) depositClaim(c context.Context, claim types.EthereumClaim) error {
	ctx := sdk.UnwrapSDKContext(c)

	orch, _ := sdk.AccAddressFromBech32(claim.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, orch)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, validator)
	if val == nil || !val.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in active set")
	}

	msg, ok := claim.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}

	// Add the claim to the store
	_, err = k.Attest(ctx, claim, any)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, claim.Type().String()),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		),
	)

	return nil
}

// WithdrawClaim handles MsgWithdrawClaim
// TODO it is possible to submit an old msgWithdrawClaim (old defined as covering an event nonce that has already been
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (k msgServer) withdrawClaim(c context.Context, claim types.EthereumClaim) error {
	ctx := sdk.UnwrapSDKContext(c)

	orch, _ := sdk.AccAddressFromBech32(claim.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, orch)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, validator)
	if val == nil || !val.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	msg, ok := claim.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}

	// Add the claim to the store
	_, err = k.Attest(ctx, claim, any)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, claim.Type().String()),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		),
	)

	return nil
}

// ERC20Deployed handles MsgERC20Deployed
func (k msgServer) eRC20DeployedClaim(c context.Context, claim types.EthereumClaim) error {
	ctx := sdk.UnwrapSDKContext(c)

	orch, _ := sdk.AccAddressFromBech32(claim.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, orch)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, validator)
	if val == nil || !val.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	msg, ok := claim.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}

	// Add the claim to the store
	_, err = k.Attest(ctx, claim, any)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, claim.Type().String()),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		),
	)

	return nil
}

// LogicCallExecutedClaim handles claims for executing a logic call on Ethereum
func (k msgServer) logicCallExecutedClaim(c context.Context, claim types.EthereumClaim) error {
	ctx := sdk.UnwrapSDKContext(c)

	orch, _ := sdk.AccAddressFromBech32(claim.GetOrchestratorAddress())
	validator := k.GetOrchestratorValidator(ctx, orch)
	if validator == nil {
		sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(orch))
		if sval == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		validator = sval.GetOperator()
	}

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, validator)
	if val == nil || !val.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	msg, ok := claim.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}

	// Add the claim to the store
	_, err = k.Attest(ctx, claim, any)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, claim.Type().String()),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		),
	)

	return nil
}
