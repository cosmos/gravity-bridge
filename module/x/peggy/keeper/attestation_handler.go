package keeper

import (
	"encoding/binary"

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
		batchID := att.Nonce.Uint64()
		b := a.keeper.GetOutgoingTXBatch(ctx, batchID)
		if b == nil {
			return types.ErrUnknown
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
		height := att.Nonce.Uint64()
		if !a.keeper.HasValsetRequest(ctx, height) {
			return types.ErrUnknown
		}

		// todo: is there any cleanup for us like:
		a.keeper.IterateValsetRequest(ctx, func(key []byte, _ types.Valset) bool {
			id := binary.BigEndian.Uint64(key)
			if id < height {
				ctx.Logger().Info("TODO: let's remove valset request", "id", id)
			}
			// todo: also remove all confirmations < height
			return false
		})
		return nil
	default:
		return sdkerrors.Wrapf(types.ErrDuplicate, "event type: %s", att.ClaimType)
	}
	return nil
}
