package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// Gravity slashes validator orchestrators for not confirming the current validator set or not
// building batch transactions for ERC20 tokens. The slash factor is defined on the module parameters
// (default 0.1%).
func (k Keeper) slash(ctx sdk.Context, params types.Params) {
	// params := k.GetParams(ctx)

	// iterate available confirmations to check if the ethereum signer matches the validators ethereum
	// address.
	// map: <address>|<confirm_type> --> bool
	confirmsByAddressType := make(map[string]map[string]bool)
	k.IterateConfirmations(ctx, func(confirmID tmbytes.HexBytes, confirm types.Confirm) bool {
		ethereumAddr := confirm.GetEthSigner()
		_, ok := confirmsByAddressType[ethereumAddr]
		if !ok {
			confirmsByAddressType[ethereumAddr] = make(map[string]bool)
		}

		confirmsByAddressType[ethereumAddr][confirm.GetType()] = true
		return false
	})

	// iterate validators by power in DESC order
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
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionBatch)
		}

		if !hasConfirmed || confirmsByType[types.ConfirmTypeSignerSet] {
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionSignerSet)
		}

		if !hasConfirmed || confirmsByType[types.ConfirmTypeLogicCall] {
			// TODO: create slash fraction for logic call
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), validator.ConsensusPower(), params.SlashFractionBatch)
		}

		// 	jail the validator if not already
		if !validator.IsJailed() {
			k.stakingKeeper.Jail(ctx, consAddr)
		}

		return false
	})

	// TODO: prune validator sets, older than 6 months, this time is chosen out of an abundance of caution
	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune events, attestations
}

// func (k Keeper) valsetSlashing(ctx sdk.Context, params types.Params) {
// 	maxHeight := uint64(0)

// 	// don't slash in the beginning before there aren't even SignersetWindow blocks yet
// 	// TODO: I don't understand the purpose of this logic
// 	if uint64(ctx.BlockHeight()) > params.SignersetWindow {
// 		maxHeight = uint64(ctx.BlockHeight()) - params.SignersetWindow
// 	}

// 	// what's an unslashed valset?
// 	unslashedValsets := k.GetUnslashedValsets(ctx, maxHeight)

// 	// unslashedValsets are sorted by nonce in ASC order
// 	for _, valset := range unslashedValsets {
// 		// TODO: use iterator here
// 		confirms := k.GetValsetConfirms(ctx, valset.Nonce)

// 		// slash bonded validators who didn't attest valset request events
// 		k.slashBondedValidators(ctx, valset.Nonce, confirms, params.SlashFractionSignerset)

// 		// slash unbonding validators who didn't attest valset request events
// 		k.slashUnbondingValidators(ctx, valset.Nonce, confirms, params.SlashFractionSignerset, params.SignedValsetsWindow)

// 		// set the latest slashed valset nonce
// 		// TODO: why every time tho??
// 		k.SetLastSlashedValsetNonce(ctx, valset.Nonce)
// 	}
// }

// func (k Keeper) slashBondedValidators(ctx sdk.Context, nonce uint64, confirms []types.Confirm, slashFraction sdk.Dec) {
// 	currentBondedSet := k.stakingKeeper.GetBondedValidatorsByPower(ctx)

// 	for _, val := range currentBondedSet {
// 		consAddr, _ := val.GetConsAddr()

// 		//  Slash validator ONLY if he joined after valset is created
// 		valSigningInfo, exist := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
// 		if !exist || valSigningInfo.StartHeight >= int64(nonce) {
// 			continue
// 		}

// 		// Check if validator has confirmed valset or not
// 		found := false

// 		for _, conf := range confirms {
// 			if conf.EthAddress == k.GetEthAddress(ctx, val.GetOperator()) {
// 				found = true
// 				break
// 			}
// 		}

// 		// slash validators for not confirming valsets
// 		if found {
// 			continue
// 		}

// 		// NOTE: this shouldn't panic because we are iterating over the bonded valset
// 		k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), val.ConsensusPower(), slashFraction)

// 		// jail the validator if not already
// 		if !val.IsJailed() {
// 			k.stakingKeeper.Jail(ctx, consAddr)
// 		}
// 	}
// }

// func (k Keeper) slashUnbondingValidators(ctx sdk.Context, nonce uint64, confirms []types.Confirm, slashFraction sdk.Dec, slashingWindow uint64) {
// 	blockTime := ctx.BlockTime().Add(k.stakingKeeper.GetParams(ctx).UnbondingTime)
// 	blockHeight := ctx.BlockHeight()

// 	// TODO: double-check this iterator on the staking module
// 	iterator := k.stakingKeeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
// 	defer iterator.Close()

// 	// All unbonding validators
// 	for ; iterator.Valid(); iterator.Next() {
// 		// TODO: ?? why is the value an array of addresses?
// 		unbondingValidators := k.GetUnbondingValidators(iterator.Value())

// 		for _, bechValAddr := range unbondingValidators.Addresses {
// 			validatorAddr, err := sdk.ValAddressFromBech32(bechValAddr)
// 			if err != nil {
// 				panic(err)
// 			}

// 			validator, found := k.stakingKeeper.GetValidator(ctx, validatorAddr)
// 			if !found {
// 				panic("validator not found")
// 			}

// 			valConsAddr, _ := validator.GetConsAddr()
// 			valSigningInfo, exist := k.slashingKeeper.GetValidatorSigningInfo(ctx, valConsAddr)

// 			// Only slash validators who joined after valset is created and they are unbonding and UNBOND_SLASHING_WINDOW didn't passed
// 			if !exist ||
// 				!validator.IsUnbonding() ||
// 				valSigningInfo.StartHeight >= int64(nonce) ||
// 				nonce >= uint64(validator.UnbondingHeight)+slashingWindow {
// 				continue
// 			}

// 			// Check if validator has confirmed valset or not
// 			found = false
// 			for _, conf := range confirms {
// 				if conf.EthAddress == k.GetEthAddress(ctx, validator.GetOperator()) {
// 					found = true
// 					break
// 				}
// 			}

// 			// slash validators for not confirming valsets
// 			if found {
// 				continue
// 			}

// 			k.stakingKeeper.Slash(ctx, valConsAddr, ctx.BlockHeight(), validator.ConsensusPower(), slashFraction)
// 			// jail the validator if not already
// 			if !validator.IsJailed() {
// 				k.stakingKeeper.Jail(ctx, valConsAddr)
// 			}
// 		}

// 	}
// }

// func (k Keeper) batchTxSlashing(ctx sdk.Context, params types.Params) {
// 	// #2 condition
// 	// We look through the full bonded set (not just the active set, include unbonding validators)
// 	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
// 	maxHeight := uint64(0)

// 	// don't slash in the beginning before there aren't even BatchTxWindow blocks yet
// 	if uint64(ctx.BlockHeight()) > params.BatchTxWindow {
// 		maxHeight = uint64(ctx.BlockHeight()) - params.BatchTxWindow
// 	}

// 	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)

// 	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
// 		// Don't slash validators who joined after batch is created
// 		consAddr, _ := validator.GetConsAddr()
// 		valSigningInfo, exist := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
// 		if exist && valSigningInfo.StartHeight > int64(batch.Block) {
// 			return false
// 		}

// 		return false
// 	})

// 	//
// 	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock+1, maxHeight, func(txID tmbytes.HexBytes, batch types.BatchTx) bool {
// 		confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.Nonce, batch.TokenContract)

// 		return false
// 	})

// 	unslashedBatches := k.GetUnslashedBatches(ctx, maxHeight)
// 	for _, batch := range unslashedBatches {

// 		// SLASH BONDED VALIDTORS who didn't attest batch requests
// 		currentBondedSet := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
// 		confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.Nonce, batch.TokenContract)

// 		for _, val := range currentBondedSet {
// 			// Don't slash validators who joined after batch is created
// 			consAddr, _ := val.GetConsAddr()
// 			valSigningInfo, exist := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
// 			if exist && valSigningInfo.StartHeight > int64(batch.Block) {
// 				continue
// 			}

// 			found := false
// 			for _, conf := range confirms {
// 				// TODO: double check this logic
// 				confVal, _ := sdk.AccAddressFromBech32(conf.OrchestratorAddress)
// 				validatorAddr := k.GetOrchestratorValidator(ctx, confVal)
// 				if validatorAddr.Equals(val.GetOperator()) {
// 					found = true
// 					break
// 				}
// 			}

// 			if found {
// 				continue
// 			}

// 			cons, _ := val.GetConsAddr()
// 			k.stakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
// 			if !val.IsJailed() {
// 				k.stakingKeeper.Jail(ctx, cons)
// 			}
// 		}

// 		// then we set the latest slashed batch block
// 		k.SetLastSlashedBatchBlock(ctx, batch.Block)
// 	}
// }
