package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EthereumEventProcessor processes `accepted` EthereumEvents
type EthereumEventProcessor struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
}

func (a EthereumEventProcessor) DetectMaliciousSupply(ctx sdk.Context, denom string, amount sdk.Int) (err error) {
	currentSupply := a.bankKeeper.GetSupply(ctx, denom)
	newSupply := new(big.Int).Add(currentSupply.Amount.BigInt(), amount.BigInt())
	if newSupply.BitLen() > 256 {
		return sdkerrors.Wrapf(types.ErrSupplyOverflow, "malicious supply of %s detected", denom)
	}
	return nil
}

// Handle is the entry point for EthereumEvent processing
func (a EthereumEventProcessor) Handle(ctx sdk.Context, eve types.EthereumEvent) (err error) {
	switch event := eve.(type) {
	case *types.SendToCosmosEvent:
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, event.TokenContract)
		addr, _ := sdk.AccAddressFromBech32(event.CosmosReceiver)
		coins := sdk.Coins{sdk.NewCoin(denom, event.Amount)}
		if !isCosmosOriginated {
			if err := a.DetectMaliciousSupply(ctx, denom, event.Amount); err != nil {
				return err
			}
			// If it is not cosmos originated, mint the coins (aka vouchers)
			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}
		}
		return a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	case *types.BatchExecutedEvent:
		a.keeper.batchTxExecuted(ctx, common.HexToAddress(event.TokenContract), event.BatchNonce)
		return
	case *types.ERC20DeployedEvent:
		// Check if it already exists
		if existingERC20, exists := a.keeper.getCosmosOriginatedERC20(ctx, event.CosmosDenom); exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 %s already exists for denom %s", existingERC20.Hex(), event.CosmosDenom))
		}

		// Check if denom exists
		// TODO: document that peggy chains require denom metadata set
		metadata, _ := a.keeper.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom)
		if metadata.Base == "" {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("denom not found %s", event.CosmosDenom))
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
		// TODO: issue event hook for consumer modules
	case *types.SignerSetTxExecutedEvent:
		// TODO here we should check the contents of the validator set against
		// the store, if they differ we should take some action to indicate to the
		// user that bridge highjacking has occurred
		a.keeper.setLastObservedSignerSetTx(ctx, types.SignerSetTx{
			Nonce:   event.SignerSetTxNonce,
			Signers: event.Members,
		})

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %T", event)
	}
	return nil
}
