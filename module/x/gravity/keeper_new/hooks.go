package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks create new gravity staking hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterValidatorBeginUnbonding persists the block height in the store
// Later in endblocker, check if there is at least one validator who started unbonding and create a valset request.
// The reason for creating valset requests in endblock is to create only one valset request per block if multiple validators starts unbonding at same block.
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {
	// TODO: this should just create the valset request instead of setting the height for each unbonding operation
	h.k.SetLastUnbondingBlockHeight(ctx, uint64(ctx.BlockHeight()))
}

// BeforeDelegationCreated performs a no-op.
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}

// AfterValidatorCreated performs a no-op.
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {}

// BeforeValidatorModified performs a no-op.
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) {}

// AfterValidatorBonded performs a no-op.
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}

// BeforeDelegationRemoved performs a no-op.
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}

// AfterValidatorRemoved performs a no-op.
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {}

// BeforeValidatorSlashed performs a no-op.
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {}

// BeforeDelegationSharesModified performs a no-op.
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}

// AfterDelegationModified performs a no-op.
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
