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
	switch att.ClaimType {
	case types.CLAIM_TYPE_ETHEREUM_BRIDGE_DEPOSIT:
		ud, err := types.UnpackAttestationDetails(att.Details)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalid, "unpacking attestation details: %s", err)
		}
		deposit, ok := ud.(*types.BridgeDeposit)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}
		if !a.keeper.HasCounterpartDenominator(ctx, types.NewVoucherDenom(types.NewEthereumAddress(string(deposit.Erc_20Token.TokenContractAddress)), deposit.Erc_20Token.Symbol)) {
			a.keeper.StoreCounterpartDenominator(ctx, types.NewEthereumAddress(string(deposit.Erc_20Token.TokenContractAddress)), deposit.Erc_20Token.Symbol)
		}
		coin := deposit.Erc_20Token.AsVoucherCoin()
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
		ud, err := types.UnpackAttestationDetails(att.Details)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalid, "unpacking attestation details: %s", err)
		}
		details, ok := ud.(*types.WithdrawalBatch)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}

		a.keeper.OutgoingTxBatchExecuted(ctx, types.NewEthereumAddress(string(details.Erc_20Token.TokenContractAddress)), types.NewUInt64Nonce(details.BatchNonce))

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", att.ClaimType)
	}
	return nil
}
