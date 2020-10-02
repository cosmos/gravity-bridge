package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper       Keeper
	supplyKeeper types.SupplyKeeper
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation) error {
	switch att.ClaimType {
	case types.ClaimTypeEthereumBridgeDeposit:
		deposit, ok := att.Details.(types.BridgeDeposit)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		if !a.keeper.HasCounterpartDenominator(ctx, types.NewVoucherDenom(deposit.ERC20Token.TokenContractAddress, deposit.ERC20Token.Symbol)) {
			a.keeper.StoreCounterpartDenominator(ctx, deposit.ERC20Token.TokenContractAddress, deposit.ERC20Token.Symbol)
		}
		coin := deposit.ERC20Token.AsVoucherCoin()
		vouchers := sdk.Coins{coin}
		err := a.supplyKeeper.MintCoins(ctx, types.ModuleName, vouchers)
		if err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}
		err = a.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, deposit.CosmosReceiver, vouchers)
		if err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}
	case types.ClaimTypeEthereumBridgeWithdrawalBatch:
		batchID, err := types.Uint64FromNonce(att.Nonce)
		if err != nil {
			return sdkerrors.Wrap(err, "nonce")
		}
		b := a.keeper.GetOutgoingTXBatch(ctx, batchID)
		if b == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "nonce")
		}
		if err := b.Observed(); err != nil {
			return err
		}
		a.keeper.storeBatch(ctx, batchID, *b)
		// cleanup outgoing TX pool
		for i := range b.Elements {
			a.keeper.removePoolEntry(ctx, b.Elements[i].ID)
		}
		return nil
	case types.ClaimTypeEthereumBridgeMultiSigUpdate:
		if !att.Nonce.GreaterThan(a.keeper.GetLastValsetObservedNonce(ctx)) {
			return types.ErrOutdated
		}
		a.keeper.setLastValsetObservedNonce(ctx, att.Nonce)

		if !a.keeper.HasValsetRequest(ctx, att.Nonce) {
			return sdkerrors.Wrap(types.ErrUnknown, "nonce")
		}

		// todo: is there any cleanup for us like:
		a.keeper.IterateValsetRequest(ctx, func(key []byte, _ types.Valset) bool {
			nonce := types.UInt64NonceFromBytes(key)
			if att.Nonce.GreaterThan(nonce) {
				ctx.Logger().Info("TODO: let's remove valset request", "nonce", nonce)
			}
			// todo: also remove all confirmations < height
			return false
		})
		return nil
	case types.ClaimTypeEthereumBootstrap:
		// todo: here and others that we do not work with an outdated nonce:
		// todo: can be a confirmed valset update already > current nonce or another bootstrap attestation
		// todo: abort when we had a previous successful processed bootstrap
		bootstrap, ok := att.Details.(types.BridgeBootstrap)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		// storing the bootstrap data here to avoid the gov process in MVY. TODO: remove this
		// todo: verify that StartThreshold == params.StartThreshold ? or before when accepting message?
		// todo: verify that PeggyID == params.PeggyID ? or before when accepting message?

		a.keeper.setPeggyID(ctx, bootstrap.PeggyID)
		a.keeper.setStartThreshold(ctx, bootstrap.StartThreshold)

		initialMultisigSet := types.Valset{
			Nonce:        att.Nonce,
			Powers:       bootstrap.ValidatorPowers,
			EthAddresses: bootstrap.AllowedValidatorSet,
		}
		// todo: do we want to do a sanity check that these validator addresses exits already?
		return a.keeper.SetBootstrapValset(ctx, att.Nonce, initialMultisigSet)
	case types.ClaimTypeOrchestratorSignedMultiSigUpdate:
		if !att.Nonce.GreaterThan(a.keeper.GetLastValsetApprovedNonce(ctx)) {
			return types.ErrOutdated
		}
		a.keeper.setLastValsetApprovedNonce(ctx, att.Nonce)

		signedCheckpoint, ok := att.Details.(types.SignedCheckpoint)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		_ = signedCheckpoint
		// todo: any cleanup to do? delete all valsets with nonce < last observed one?
		return nil
	case types.ClaimTypeOrchestratorSignedWithdrawBatch:
		signedCheckpoint, ok := att.Details.(types.SignedCheckpoint)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		_ = signedCheckpoint
		// todo: any cleanup to do? delete all withdraw batches with nonce < last observed one?
		return nil
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", att.ClaimType)
	}
	return nil
}
