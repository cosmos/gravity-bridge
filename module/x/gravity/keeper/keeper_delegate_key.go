package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"
)

/////////////////////////////
//    ADDRESS DELEGATION   //
/////////////////////////////

// SetOrchestratorValidator sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidator(ctx sdk.Context, val sdk.ValAddress, orch sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOrchestratorAddressKey(orch), val.Bytes())
}

// GetOrchestratorValidator returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidator(ctx sdk.Context, orch sdk.AccAddress) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	valAddr := store.Get(types.GetOrchestratorAddressKey(orch))
	if valAddr == nil {
		return stakingtypes.Validator{}, false
	}
	validator, found = k.StakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, false
	}

	return validator, true
}

/////////////////////////////
//       ETH ADDRESS       //
/////////////////////////////

// SetEthAddress sets the ethereum address for a given validator
func (k Keeper) SetEthAddressForValidator(ctx sdk.Context, validator sdk.ValAddress, ethAddr string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEthAddressByValidatorKey(validator), []byte(ethAddr))
	store.Set(types.GetValidatorByEthAddressKey(ethAddr), []byte(validator))
}

// GetEthAddressByValidator returns the eth address for a given gravity validator
func (k Keeper) GetEthAddressByValidator(ctx sdk.Context, validator sdk.ValAddress) (ethAddress string, found bool) {
	store := ctx.KVStore(k.storeKey)
	ethAddr := store.Get(types.GetEthAddressByValidatorKey(validator))
	if ethAddr == nil {
		return "", false
	} else {
		return string(ethAddr), true
	}
}

// GetValidatorByEthAddress returns the validator for a given eth address
func (k Keeper) GetValidatorByEthAddress(ctx sdk.Context, ethAddr string) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	valAddr := store.Get(types.GetValidatorByEthAddressKey(ethAddr))
	if valAddr == nil {
		return stakingtypes.Validator{}, false
	}
	validator, found = k.StakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, false
	}

	return validator, true
}

// GetCurrentValset gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'gravity power' is computed as
// Cosmos power for that validator / total cosmos power = x / uint32 Max
// where x is the voting power on the gravity contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
//
// 'total cosmos power' has an edge case, if a validator has not set their
// Ethereum key they are not included in the total. If they where control
// of the bridge could be lost in the following situation.
//
// If we have 100 total power, and 100 total power joins the validator set
// the new validators hold more than 33% of the bridge power, if we generate
// and submit a valset and they don't have their eth keys set they can never
// update the validator set again and the bridge and all its' funds are lost.
// For this reason we exclude validators with unset eth keys from validator sets
func (k Keeper) GetCurrentValset(ctx sdk.Context) *types.Valset {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	// allocate enough space for all validators, but len zero, we then append
	// so that we have an array with extra capacity but the correct length depending
	// on how many validators have keys set.
	bridgeValidators := make([]*types.BridgeValidator, 0, len(validators))
	var totalPower uint64
	// TODO someone with in depth info on Cosmos staking should determine
	// if this is doing what I think it's doing
	for _, validator := range validators {
		val := validator.GetOperator()

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, val))

		if ethAddr, found := k.GetEthAddressByValidator(ctx, val); found {
			bv := &types.BridgeValidator{Power: p, EthereumAddress: ethAddr}
			bridgeValidators = append(bridgeValidators, bv)
			totalPower += p
		}
	}
	// normalize power values
	for i := range bridgeValidators {
		bridgeValidators[i].Power = sdk.NewUint(bridgeValidators[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	// get the reward from the params store
	reward := k.GetParams(ctx).ValsetReward
	var rewardToken string
	var rewardAmount sdk.Int
	if !reward.IsValid() || reward.IsZero() {
		// the case where a validator has 'no reward'. The 'no reward' value is interpreted as having a zero
		// address for the ERC20 token and a zero value for the reward amount. Since we store a coin with the
		// params, a coin with a blank denom and/or zero amount is interpreted in this way.
		rewardToken = "0x0000000000000000000000000000000000000000"
		rewardAmount = sdk.NewIntFromUint64(0)

	} else {
		rewardToken, rewardAmount = k.RewardToERC20Lookup(ctx, reward)
	}

	// TODO: make the nonce an incrementing one (i.e. fetch last nonce from state, increment, set here)
	return types.NewValset(uint64(ctx.BlockHeight()), uint64(ctx.BlockHeight()), bridgeValidators, rewardAmount, rewardToken)
}
