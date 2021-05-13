package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EthereumEventVoteRecordHandler processes `observed` EthereumEventVoteRecords
type EthereumEventVoteRecordHandler struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
}

// Handle is the entry point for EthereumEventVoteRecord processing.
func (a EthereumEventVoteRecordHandler) Handle(ctx sdk.Context, voteRecord types.EthereumEventVoteRecord, event types.EthereumEvent) error {
	switch event := event.(type) {
	case *types.MsgSendToCosmosEvent:
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, event.TokenContract)

		if isCosmosOriginated {
			// If it is cosmos originated, unlock the coins
			coins := sdk.Coins{sdk.NewCoin(denom, event.Amount)}

			addr, err := sdk.AccAddressFromBech32(event.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid receiver address")
			}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		} else {
			// If it is not cosmos originated, mint the coins (aka vouchers)
			coins := sdk.Coins{sdk.NewCoin(denom, event.Amount)}

			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}

			addr, err := sdk.AccAddressFromBech32(event.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid receiver address")
			}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		}
	case *types.MsgBatchExecutedEvent:
		return a.keeper.BatchTxExecuted(ctx, event.TokenContract, event.BatchNonce)
	case *types.MsgERC20DeployedEvent:
		// Check if it already exists
		existingERC20, exists := a.keeper.GetCosmosOriginatedERC20(ctx, event.CosmosDenom)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 %s already exists for denom %s", existingERC20, event.CosmosDenom))
		}

		// Check if denom exists
		metadata := a.keeper.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom)
		if metadata.Base == "" {
			return sdkerrors.Wrap(types.ErrUnknown, fmt.Sprintf("denom not found %s", event.CosmosDenom))
		}

		// Check if attributes of ERC20 match Cosmos denom
		if event.Name != metadata.Display {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 name %s does not match denom display %s", event.Name, metadata.Description))
		}

		if event.Symbol != metadata.Display {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 symbol %s does not match denom display %s", event.Symbol, metadata.Display))
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

		if decimals != uint32(event.Decimals) {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 decimals %d does not match denom decimals %d", event.Decimals, decimals))
		}

		// Add to denom-erc20 mapping
		a.keeper.setCosmosOriginatedDenomToERC20(ctx, event.CosmosDenom, event.TokenContract)
	case *types.MsgSignerSetUpdatedEvent:
		// TODO here we should check the contents of the validator set against
		// the store, if they differ we should take some action to indicate to the
		// user that bridge highjacking has occurred
		a.keeper.SetLastObservedSignerSetTx(ctx, types.SignerSetTx{
			Nonce:   event.SignerSetNonce,
			Members: event.Members,
		})

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", event.GetType())
	}
	return nil
}
