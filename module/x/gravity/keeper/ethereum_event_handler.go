package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EthereumEventProcessor processes `accepted` EthereumEvents
type EthereumEventProcessor struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
}

// Handle is the entry point for EthereumEvent processing
func (a EthereumEventProcessor) Handle(ctx sdk.Context, eve types.EthereumEvent) (err error) {
	switch event := eve.(type) {
	case *types.SendToCosmosEvent:
		// Check if coin is Cosmos-originated asset and get denom
		if isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, event.TokenContract); isCosmosOriginated {
			// If it is cosmos originated, unlock the coins
			addr, _ := sdk.AccAddressFromBech32(event.CosmosReceiver)
			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{sdk.NewCoin(denom, event.Amount)}); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		} else {
			// If it is not cosmos originated, mint the coins (aka vouchers)
			coins := sdk.Coins{sdk.NewCoin(denom, event.Amount)}
			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}

			addr, _ := sdk.AccAddressFromBech32(event.CosmosReceiver)
			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		}
	case *types.BatchExecutedEvent:
		return a.keeper.BatchTxExecuted(ctx, event.TokenContract, event.GetNonce())
	case *types.ERC20DeployedEvent:
		// Check if it already exists
		if existingERC20, exists := a.keeper.GetCosmosOriginatedERC20(ctx, event.CosmosDenom); exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 %s already exists for denom %s", existingERC20, event.CosmosDenom))
		}

		// Check if denom exists
		// TODO: document that peggy chains require denom metadata set
		metadata := a.keeper.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom)
		if metadata.Base == "" {
			return sdkerrors.Wrap(types.ErrUnknown, fmt.Sprintf("denom not found %s", event.CosmosDenom))
		}

		// Check if attributes of ERC20 match Cosmos denom
		if event.Erc20Name != metadata.Display {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 name %s does not match denom display %s", event.Erc20Name, metadata.Description))
		}

		if event.Erc20Symbol != metadata.Display {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 symbol %s does not match denom display %s", event.Erc20Symbol, metadata.Display))
		}

		// ERC20 tokens use a very simple mechanism to tell you where to display the decimal point.
		// The "decimals" field simply tells you how many decimal places there will be.
		// Cosmos denoms have a system that is much more full featured, with enterprise-ready token denominations.
		// There is a DenomUnits array that tells you what the name of each denomination of the
		// token is.
		// To correlate this with an ERC20 "decimals" field, we have to search through the DenomUnits array
		// to find the DenomUnit which matches up to the main token "display" value. Then we take the
		// "exponent" from this DenomUnit.
		// If the correct DenomUnit is not found, it will default to 0. This will result in there being no decimal places
		// in the token's ERC20 on Ethereum. So, for example, if this happened with Atom, 1 Atom would appear on Ethereum
		// as 1 million Atoms, having 6 extra places before the decimal point.
		// This will only happen with a Denom Metadata which is for all intents and purposes invalid, but I am not sure
		// this is checked for at any other point.
		decimals := uint32(0)
		for _, denomUnit := range metadata.DenomUnits {
			if denomUnit.Denom == metadata.Display {
				decimals = denomUnit.Exponent
				break
			}
		}

		if decimals != uint32(event.Erc20Decimals) {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 decimals %d does not match denom decimals %d", event.Erc20Decimals, decimals))
		}

		// Add to denom-erc20 mapping
		a.keeper.setCosmosOriginatedDenomToERC20(ctx, event.CosmosDenom, event.TokenContract)

	case *types.ContractCallExecutedEvent:
		// todo: issue event hook for consumer modules
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %T", event)
	}
	return nil
}
