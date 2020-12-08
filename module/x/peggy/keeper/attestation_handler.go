package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation) error {
	ud, err := types.UnpackEthereumClaim(att.Details)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "unpacking attestation details: %s", err)
	}
	switch ud.GetType() {
	case types.CLAIM_TYPE_ETHEREUM_BRIDGE_DEPOSIT:
		deposit, ok := ud.(*types.EthereumBridgeDepositClaim)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		coin := deposit.Erc20Token.PeggyCoin()
		vouchers := sdk.Coins{coin}
		if err = a.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}

		addr, err := sdk.AccAddressFromBech32(deposit.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid reciever address")
		}
		if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, vouchers); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}

	case types.CLAIM_TYPE_ETHEREUM_BRIDGE_WITHDRAWAL_BATCH:
		details, ok := ud.(*types.EthereumBridgeWithdrawalBatchClaim)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}

		a.keeper.OutgoingTxBatchExecuted(ctx, details.Erc20Token.Contract, details.BatchNonce)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", ud.GetType())
	}
	return nil
}
