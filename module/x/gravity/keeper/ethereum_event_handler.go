package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EthereumEventProcessor processes `accepted` EthereumEvents
type EthereumEventProcessor struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
}

func (a EthereumEventProcessor) DetectMaliciousSupply(ctx sdk.Context, denom string, amount sdk.Int) (err error) {
	currentSupply := a.keeper.bankKeeper.GetSupply(ctx, denom)
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

			// if it is not cosmos originated, mint the coins (aka vouchers)
			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}
		}

		if err := a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
			return err
		}
		a.keeper.AfterSendToCosmosEvent(ctx, *event)
		return nil

	case *types.BatchExecutedEvent:
		a.keeper.batchTxExecuted(ctx, common.HexToAddress(event.TokenContract), event.BatchNonce)
		a.keeper.AfterBatchExecutedEvent(ctx, *event)
		return nil

	case *types.ERC20DeployedEvent:
		if err := a.verifyERC20DeployedEvent(ctx, event); err != nil {
			return err
		}

		// add to denom-erc20 mapping
		a.keeper.setCosmosOriginatedDenomToERC20(ctx, event.CosmosDenom, event.TokenContract)
		a.keeper.AfterERC20DeployedEvent(ctx, *event)
		return nil

	case *types.ContractCallExecutedEvent:
		a.keeper.AfterContractCallExecutedEvent(ctx, *event)
		return nil

	case *types.SignerSetTxExecutedEvent:
		// TODO here we should check the contents of the validator set against
		// the store, if they differ we should take some action to indicate to the
		// user that bridge highjacking has occurred
		a.keeper.setLastObservedSignerSetTx(ctx, types.SignerSetTx{
			Nonce:   event.SignerSetTxNonce,
			Signers: event.Members,
		})
		a.keeper.AfterSignerSetExecutedEvent(ctx, *event)
		return nil

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %T", event)
	}
}

func (a EthereumEventProcessor) verifyERC20DeployedEvent(ctx sdk.Context, event *types.ERC20DeployedEvent) error {
	if existingERC20, exists := a.keeper.getCosmosOriginatedERC20(ctx, event.CosmosDenom); exists {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"ERC20 token %s already exists for denom %s", existingERC20.Hex(), event.CosmosDenom,
		)
	}

	// We expect that all Cosmos-based tokens have metadata defined. In the case
	// a token does not have metadata defined, e.g. an IBC token, we successfully
	// handle the token under the following conditions:
	//
	// 1. The ERC20 name is equal to the token's denomination. Otherwise, this
	// 		means that ERC20 tokens would have an untenable UX.
	// 2. The ERC20 token has zero decimals as this is what we default to since
	// 		we cannot know or infer the real decimal value for the Cosmos token.
	// 3. The ERC20 symbol is empty.
	//
	// NOTE: This path is not encouraged and all supported assets should have
	// metadata defined. If metadata cannot be defined, consider adding the token's
	// metadata on the fly.
	if md, ok := a.keeper.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom); ok && md.Base != "" {
		return verifyERC20Token(md, event)
	}

	if supply := a.keeper.bankKeeper.GetSupply(ctx, event.CosmosDenom); supply.IsZero() {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"no supply exists for token %s without metadata", event.CosmosDenom,
		)
	}

	if event.Erc20Name != event.CosmosDenom {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"invalid ERC20 name for token without metadata; got: %s, expected: %s", event.Erc20Name, event.CosmosDenom,
		)
	}

	if event.Erc20Symbol != "" {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"expected empty ERC20 symbol for token without metadata; got: %s", event.Erc20Symbol,
		)
	}

	if event.Erc20Decimals != 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"expected zero ERC20 decimals for token without metadata; got: %d", event.Erc20Decimals,
		)
	}

	return nil
}

func verifyERC20Token(metadata banktypes.Metadata, event *types.ERC20DeployedEvent) error {
	if event.Erc20Name != metadata.Display {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"ERC20 name %s does not match the denom display %s", event.Erc20Name, metadata.Display,
		)
	}

	if event.Erc20Symbol != metadata.Display {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"ERC20 symbol %s does not match denom display %s", event.Erc20Symbol, metadata.Display,
		)
	}

	// ERC20 tokens use a very simple mechanism to tell you where to display the
	// decimal point. The "decimals" field simply tells you how many decimal places
	// there will be.
	//
	// Cosmos denoms have a system that is much more full featured, with
	// enterprise-ready token denominations. There is a DenomUnits array that
	// tells you what the name of each denomination of the token is.
	//
	// To correlate this with an ERC20 "decimals" field, we have to search through
	// the DenomUnits array to find the DenomUnit which matches up to the main
	// token "display" value. Then we take the "exponent" from this DenomUnit.
	//
	// If the correct DenomUnit is not found, it will default to 0. This will
	// result in there being no decimal places in the token's ERC20 on Ethereum.
	// For example, if this happened with ATOM, 1 ATOM would appear on Ethereum
	// as 1 million ATOM, having 6 extra places before the decimal point.
	var decimals uint32
	for _, denomUnit := range metadata.DenomUnits {
		if denomUnit.Denom == metadata.Display {
			decimals = denomUnit.Exponent
			break
		}
	}

	if uint64(decimals) != event.Erc20Decimals {
		return sdkerrors.Wrapf(
			types.ErrInvalidERC20Event,
			"ERC20 decimals %d does not match denom decimals %d", event.Erc20Decimals, decimals,
		)
	}

	return nil
}
