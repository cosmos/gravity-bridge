package keeper

import (
	"fmt"

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
// TODO-JT add handler for ERC20DeployedEvent
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, denom := a.keeper.ERC20ToDenom(ctx, claim.TokenContract)

		if isCosmosOriginated {
			// If it is cosmos originated, unlock the coins
			coins := sdk.Coins{sdk.NewCoin(denom, claim.Amount)}

			addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		} else {
			// If it is not cosmos originated, mint the coins (aka vouchers)
			coins := sdk.Coins{sdk.NewCoin(denom, claim.Amount)}

			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}

			addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		}
	case *types.MsgWithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)
	case *types.MsgERC20DeployedClaim:
		// Check if attributes of ERC20 match Cosmos denom
		metadata, exists := a.keeper.bankKeeper.GetDenomMetaData(ctx, claim.CosmosDenom)

		if !exists {
			return sdkerrors.Wrap(types.ErrUnknown, fmt.Sprintf("denom not found %s", claim.CosmosDenom))
		}

		if claim.Name != metadata.Description {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 name %s does not match denom description %s", claim.Name, metadata.Description))
		}

		if claim.Symbol != metadata.Display {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 symbol %s does not match denom display %s", claim.Symbol, metadata.Display))
		}

		if claim.Decimals != uint64(metadata.DenomUnits[0].Exponent) {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 decimals %d does not match denom denomunits exponent %d", claim.Decimals, uint64(metadata.DenomUnits[0].Exponent)))
		}

		// Add to denom-erc20 mapping
		a.keeper.setCosmosOriginatedDenomToERC20(ctx, claim.CosmosDenom, claim.TokenContract)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
