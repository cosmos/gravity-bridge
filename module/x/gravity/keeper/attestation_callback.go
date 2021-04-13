package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// AttestationHandler defines an interface that processes incoming attestations
// from Ethereum. While the default handler only mints ERC20 tokens, additional
// custom functionality can be implemented by passing an external handler to the
// bridge keeper.
//
// Examples of custom functionality could be, but not limited to:
//
// - Transfering newly minted ERC20 tokens (represented as an sdk.Coins) to a
// given recipient, either local or via IBC to a counterparty chain
//
// - Pooling the tokens into an escrow account for interest accruing DeFi solutions
//
// - Deposit into an AMM pair
type AttestationHandler interface {
	OnAttestation(ctx sdk.Context, attestation types.Attestation) error
}

// DefaultAttestationHandler is the default handler for processing observed
// event attestations received from Ethereum.
type DefaultAttestationHandler struct {
	keeper Keeper
}

var _ AttestationHandler = DefaultAttestationHandler{}

// OnAttestation processes ethereum event upon attestation and performs a custom
// logic.
//
// TODO: add handler for ERC20DeployedEvent
// TODO: clean up
func (handler DefaultAttestationHandler) OnAttestation(ctx sdk.Context, attestation types.Attestation) error {
	event, found := handler.keeper.GetEthEvent(ctx, attestation.EventID)
	if !found {
		// TODO: err msg
		return fmt.Errorf("not found")
	}

	switch event := event.(type) {
	case *types.DepositEvent:
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, event.TokenContract)

		if isCosmosOriginated {
			// If it is cosmos originated, unlock the coins
			coins := sdk.Coins{sdk.NewCoin(denom, event.Amount)}

			addr, err := sdk.AccAddressFromBech32(event.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
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
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		}
	case *types.WithdrawEvent:
		a.keeper.OutgoingTxBatchExecuted(ctx, event.TokenContract, event.BatchNonce)
	case *types.ERC20DeployedEvent:
		return RegisterERC20(a.keeper, ctx, event)
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", event.GetType())
	}

	return nil
}

// RegisterERC20
func RegisterERC20(k Keeper, ctx sdk.Context, event types.CosmosERC20DeployedEvent) error {
	// Check if it already exists
	contractAddr, found := k.GetERC20ContractFromCoinDenom(ctx, event.CosmosDenom)
	if found {
		return sdkerrors.Wrap(
			// TODO: fix
			types.ErrContractNotFound,
			fmt.Sprintf("erc20 contract %s already registered for coin denom %s", contractAddr.String(), event.CosmosDenom))
	}

	// Check if denom exists
	metadata := k.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom)

	if err := validateCoinMetadata(event, metadata); err != nil {
		return err
	}

	// Add to denom-erc20 mapping
	return nil
}

func validateCoinMetadata(event types.CosmosERC20DeployedEvent, metadata banktypes.Metadata) error {
	if metadata.Base == "" {
		return sdkerrors.Wrapf(types.ErrUnknown, "denom not found %s", event.CosmosDenom)
	}

	// Check if attributes of ERC20 match Cosmos denom
	if event.Name != metadata.Display {
		return sdkerrors.Wrapf(
			// TODO: fix
			types.ErrContractNotFound,
			"ERC20 name %s does not match denom display %s", event.Name, metadata.Description)
	}

	if event.Symbol != metadata.Display {
		return sdkerrors.Wrapf(
			// TODO: fix
			types.ErrContractNotFound,
			"ERC20 symbol %s does not match denom display %s", event.Symbol, metadata.Display)
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
		return sdkerrors.Wrapf(
			// TODO: fix
			types.ErrContractNotFound,
			"ERC20 decimals %d does not match denom decimals %d", event.Decimals, decimals)
	}

	return nil
}
