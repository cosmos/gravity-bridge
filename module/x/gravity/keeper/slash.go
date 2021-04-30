package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// Gravity slashes validator orchestrators for not confirming the current ethereum signer set or not
// building batch transactions for ERC20 tokens. The slash factor is defined, per type, on the module
// parameters.
func (k Keeper) slash(ctx sdk.Context, params types.Params) {
	// iterate available confirmations to check if the ethereum signer matches the validators ethereum
	// address.
	// map: <address>|<confirm_type> --> bool
	confirmsByAddressType := make(map[string]map[string]bool)
	k.IterateConfirmations(ctx, func(_ tmbytes.HexBytes, confirm types.Confirm) bool {
		ethereumAddr := confirm.GetEthSigner()
		_, ok := confirmsByAddressType[ethereumAddr]
		if !ok {
			confirmsByAddressType[ethereumAddr] = make(map[string]bool)
		}

		confirmsByAddressType[ethereumAddr][confirm.GetType()] = true
		return false
	})

	// iterate validators by power in DESC order and check if they signed the required confirmations
	k.IterateValidatorsByPower(ctx, func(validator stakingtypes.Validator) bool {
		// validators with unbonded status can't be slashed, so we continue with next one
		if validator.IsUnbonded() {
			return false
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			panic(fmt.Errorf("failed to get validator's %s consensus address: %w", validator.OperatorAddress, err))
		}

		_, exist := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
		if exist {
			k.Logger(ctx).Debug("signing info not found for validator", "consensus-address", consAddr.String())
			return false
		}

		// TODO: should we slash once per slashing event (signer set req, batch req and logic call) or just once?
		// TODO: if slashing once, do we slash with the max of the 3 slash fractions, min, avg?

		ethereumAddr := k.GetEthAddress(ctx, validator.GetOperator())

		// TODO: should we remove this check? if so, the validator will always be slashed if it hasn't submitted the ethereum key
		if (ethereumAddr == common.Address{}) {
			k.Logger(ctx).Debug("ethereum signing address not found for validator", "validator-address", validator.OperatorAddress)
			return false
		}

		// check the requests that the validator signed or not the confirms
		confirmsByType, hasConfirmed := confirmsByAddressType[ethereumAddr.String()]

		// slash up to 3 times in case the validator didn't sign the confirms

		if !hasConfirmed || confirmsByType[types.ConfirmTypeBatch] {
			k.Logger(ctx).Debug("slashing validator for not signing batch confirms", "consensus-address", consAddr.String())
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionBatch)
		}

		// FIXME: figure out the unslashed signer sets and only slash after the signer set nonce
		if !hasConfirmed || confirmsByType[types.ConfirmTypeSignerSet] {
			k.Logger(ctx).Debug("slashing validator for not signing signer set confirms", "consensus-address", consAddr.String())
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionSignerSet)
		}

		if !hasConfirmed || confirmsByType[types.ConfirmTypeLogicCall] {
			// TODO: create slash fraction for logic call
			k.Logger(ctx).Debug("slashing validator for not signing logic call confirms", "consensus-address", consAddr.String())
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionBatch)
		}

		// 	jail the validator if not already
		if !validator.IsJailed() {
			k.stakingKeeper.Jail(ctx, consAddr)
		}

		return false
	})
}
