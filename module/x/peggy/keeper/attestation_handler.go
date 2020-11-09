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
		details, ok := att.Details.(types.WithdrawalBatch)
		if !ok {
			return sdkerrors.Wrapf(types.ErrInvalid, "unexpected type: %T", att.Details)
		}

		b := a.keeper.GetOutgoingTXBatch(ctx, details.BatchNonce)
		if b == nil {
			return sdkerrors.Wrap(types.ErrUnknown, "nonce")
		}

		// cleanup outgoing TX pool
		for i := range b.Elements {
			a.keeper.removePoolEntry(ctx, b.Elements[i].ID)
		}

		// Iterate through remaining batches
		a.keeper.IterateOutgoingTXBatches(ctx, func(key []byte, iter_batch types.OutgoingTxBatch) bool {
			// If the iterated batches nonce is lower than the one that was just executed, cancel it
			// TODO: iterate only over batches we need to iterate over
			if iter_batch.Nonce < b.Nonce {
				a.keeper.CancelOutgoingTXBatch(ctx, iter_batch.Nonce)
			}
			return false
		})
		a.keeper.deleteBatch(ctx, *b)
		return nil

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", att.ClaimType)
	}
	return nil
}
