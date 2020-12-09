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
	switch claim := ud.(type) {
	case *types.MsgDepositClaim:
		token := types.ERC20Token{
			claim.Amount,
			claim.TokenContract,
		}
		coin := token.PeggyCoin()
		vouchers := sdk.Coins{coin}
		if err = a.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}

		addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid reciever address")
		}
		if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, vouchers); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}
	case *types.DepositClaim:
		coin := claim.Erc20Token.PeggyCoin()
		vouchers := sdk.Coins{coin}
		if err = a.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}

		addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid reciever address")
		}
		if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, vouchers); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}
	case *types.WithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.Erc20Token.Contract, claim.BatchNonce)
	case *types.MsgWithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", ud.GetType())
	}
	return nil
}
