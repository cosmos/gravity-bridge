//nolint:dupl
package ethbridge

import (
	"fmt"
	"strconv"

	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper,
	cdc *codec.Codec) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateEthBridgeClaim:
			return handleMsgCreateEthBridgeClaim(ctx, cdc, bridgeKeeper, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgLock:
			return handleMsgLock(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create a bridge claim
func handleMsgCreateEthBridgeClaim(ctx sdk.Context, cdc *codec.Codec,
	bridgeKeeper Keeper,
	msg MsgCreateEthBridgeClaim) sdk.Result {

	status, sdkErr := bridgeKeeper.ProcessClaim(ctx, types.EthBridgeClaim(msg))
	if sdkErr != nil {
		return sdkErr.Result()
	}

	if status.Text == oracle.SuccessStatusText {
		sdkErr = bridgeKeeper.ProcessSuccessfulClaim(ctx, status.FinalClaim)
		if sdkErr != nil {
			return sdkErr.Result()
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyEthereumSender, msg.EthereumSender.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.CosmosReceiver.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyClaimType, msg.ClaimType.String()),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurn(ctx sdk.Context, cdc *codec.Codec,
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper, msg MsgBurn) sdk.Result {
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		return sdk.ErrInvalidAddress(msg.CosmosSender.String()).Result()
	}

	err := bridgeKeeper.ProcessBurn(ctx, msg.CosmosSender, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(types.AttributeKeyTokenContract, msg.TokenContract.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}

}

func handleMsgLock(ctx sdk.Context, cdc *codec.Codec,
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper, msg MsgLock) sdk.Result {
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		return sdk.ErrInvalidAddress(msg.CosmosSender.String()).Result()
	}

	err := bridgeKeeper.ProcessLock(ctx, msg.CosmosSender, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(types.AttributeKeyTokenContract, msg.TokenContract.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}

}
