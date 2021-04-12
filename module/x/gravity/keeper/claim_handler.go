package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// HandleClaim handles a given claim by attesting it
// TODO: it's not clear the utility of this from the code. Explain what it does,
// provice example and where this is executed on the step-by-step incoming logic.
func (k Keeper) HandleClaim(ctx sdk.Context, claim types.EthereumClaim) error {
	orch, _ := sdk.AccAddressFromBech32(claim.GetOrchestratorAddress())
	validatorAddr := k.GetOrchestratorValidator(ctx, orch)
	if validatorAddr == nil {
		validatorAddr = sdk.ValAddress(orch)
	}

	// return an error if the validator isn't in the active set
	validator := k.stakingKeeper.Validator(ctx, validatorAddr)
	if validator == nil {
		return sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, validatorAddr.String())
	} else if !validator.IsBonded() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "validator %s not in active set", validatorAddr)
	}

	// Add the claim to the store
	if err := k.Attest(ctx, claim); err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, claim.GetType().String()),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		),
	)

	return nil
}
