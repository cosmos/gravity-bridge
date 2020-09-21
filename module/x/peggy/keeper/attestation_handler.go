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
			return sdkerrors.Wrapf(types.ErrInternal, "unexpected type: %T", att.Details)
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
		batchID := att.Nonce.AsUint64()
		b := a.keeper.GetOutgoingTXBatch(ctx, batchID)
		if b == nil {
			return types.ErrUnknown
		}
		if err := b.Observed(); err != nil {
			return err
		}
		a.keeper.storeBatch(ctx, batchID, *b)
		if err := a.keeper.UpdateLastObservedBatchID(ctx, batchID); err != nil {
			return nil
		}
		// cleanup outgoing TX pool
		for i := range b.Elements {
			a.keeper.removePoolEntry(ctx, b.Elements[i].ID)
		}
		return nil
	case types.ClaimTypeEthereumBridgeMultiSigUpdate:
		// todo: update nonce for "MultiSig Set"
	default:
		return sdkerrors.Wrapf(types.ErrDuplicate, "event type: %s", att.ClaimType)
	}
	return nil
}
