package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

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

var _ AttestationHandler = DefaultAttestationHandler{}

// DefaultAttestationHandler is the default handler for processing observed
// event attestations received from Ethereum.
type DefaultAttestationHandler struct {
	keeper Keeper
}

// NewAttestationHandler creates a default attestation handler instance
func NewAttestationHandler(k Keeper) AttestationHandler {
	return &DefaultAttestationHandler{
		keeper: k,
	}
}

// OnAttestation processes ethereum event upon attestation and performs a custom
// logic.
//
// TODO: clean up
func (h DefaultAttestationHandler) OnAttestation(ctx sdk.Context, attestation types.Attestation) error {
	event, found := h.keeper.GetEthereumEvent(ctx, attestation.EventID)
	if !found {
		return sdkerrors.Wrap(types.ErrEventNotFound, attestation.EventID.String())
	}

	switch event := event.(type) {
	case *types.DepositEvent:
		return h.keeper.OnReceiveDeposit(ctx, event)
	case *types.WithdrawEvent:
		tokenContract := common.HexToAddress(event.TokenContract)
		return h.keeper.OnBatchTxExecuted(ctx, tokenContract, event.TxID)
	case *types.CosmosERC20DeployedEvent:
		return h.keeper.RegisterERC20(ctx, event)
	default:
		return sdkerrors.Wrapf(types.ErrEventUnsupported, "event type %s: %T", event.GetType(), event)
	}
}

// OnReceiveDeposit
func (k Keeper) OnReceiveDeposit(ctx sdk.Context, event *types.DepositEvent) error {
	tokenContract := common.HexToAddress(event.TokenContract)
	// Check if coin is Cosmos-originated asset and get denom
	denom, isCosmosCoin := k.GetCoinDenomFromERC20Contract(ctx, tokenContract)

	coins := sdk.Coins{}
	if isCosmosCoin {
		//  if the coin is a native cosmos coin, unlock the coins and transfer to recipient
		coins = sdk.Coins{sdk.NewCoin(denom, event.Amount)}
	} else {
		// if the coin is an ERC20 token, mint the ERC20 token vouchers and transfer to recipient
		voucherDenom := types.GravityDenom(event.TokenContract)
		coins := sdk.Coins{sdk.NewCoin(voucherDenom, event.Amount)}

		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return sdkerrors.Wrapf(err, "mint erc20 vouchers: %s", coins)
		}
	}

	addr, err := sdk.AccAddressFromBech32(event.CosmosReceiver)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid receiver address")
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

// RegisterERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, event *types.CosmosERC20DeployedEvent) error {
	// Check if it already exists
	contractAddr, found := k.GetERC20ContractFromCoinDenom(ctx, event.CosmosDenom)
	if found {
		return sdkerrors.Wrap(
			types.ErrContractExists,
			fmt.Sprintf("erc20 contract %s already registered for coin denom %s", contractAddr.String(), event.CosmosDenom))
	}

	// Check if denom exists
	metadata := k.bankKeeper.GetDenomMetaData(ctx, event.CosmosDenom)

	// NOTE: this will fail on all IBC vouchers or any Cosmos coin that hasn't
	// a denom metadata value defined
	// TODO: discuss if we should create/set a new metadata if it's not currently
	// set to store for the given cosmos denom
	if err := validateCoinMetadata(*event, metadata); err != nil {
		return err
	}

	tokenContract := common.HexToAddress(event.TokenContract)
	k.setERC20DenomMap(ctx, event.CosmosDenom, tokenContract)

	k.Logger(ctx).Debug("erc20 token registered", "contract-address", event.TokenContract, "cosmos-denom", event.CosmosDenom)
	return nil
}

// validateCoinMetadata performs a stateless validation on the metadata fields and compares its values
// with the deployed ERC20 contract values.
func validateCoinMetadata(event types.CosmosERC20DeployedEvent, metadata banktypes.Metadata) error {
	if err := metadata.Validate(); err != nil {
		return err
	}

	// Check if attributes of ERC20 match Cosmos denom
	if event.Name != metadata.Display {
		return sdkerrors.Wrapf(
			types.ErrEventInvalid,
			"ERC20 name %s does not match denom display %s", event.Name, metadata.Description)
	}

	if event.Symbol != metadata.Display {
		return sdkerrors.Wrapf(
			types.ErrEventInvalid,
			"ERC20 symbol %s does not match denom display %s", event.Symbol, metadata.Display)
	}

	// NOTE: denomination units can't be empty and are sorted in ASC order
	decimals := metadata.DenomUnits[len(metadata.DenomUnits)-1].Exponent

	if decimals != uint32(event.Decimals) {
		return sdkerrors.Wrapf(
			types.ErrEventInvalid,
			"ERC20 decimals %d does not match denom decimals %d", event.Decimals, decimals)
	}

	return nil
}
