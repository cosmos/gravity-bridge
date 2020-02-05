//nolint:dupl
package ethbridge

import (
	"fmt"
	"strconv"

	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
	"github.com/cosmos/peggy/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "nftbridge" type messages.
func NewHandler(
	nftKeeper types.NFTKeeper, bridgeKeeper Keeper,
	cdc *codec.Codec) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateNFTBridgeClaim:
			return handleMsgCreateNFTBridgeClaim(ctx, cdc, bridgeKeeper, msg)
		case MsgBurnNFT:
			return handleMsgBurn(ctx, cdc, nftKeeper, bridgeKeeper, msg)
		case MsgLockNFT:
			return handleMsgLock(ctx, cdc, nftKeeper, bridgeKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized nftbridge message type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle a message to create a bridge claim
func handleMsgCreateNFTBridgeClaim(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgCreateNFTBridgeClaim,
) (*sdk.Result, error) {
	status, err := bridgeKeeper.ProcessClaim(ctx, types.NFTBridgeClaim(msg))
	if err != nil {
		return nil, err
	}
	if status.Text == oracle.SuccessStatusText {
		if err := bridgeKeeper.ProcessSuccessfulClaim(ctx, status.FinalClaim); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			ethbridge.EventTypeCreateClaim,
			sdk.NewAttribute(ethbridge.AttributeKeyEthereumSender, msg.EthereumSender.String()),
			sdk.NewAttribute(ethbridge.AttributeKeyCosmosReceiver, msg.CosmosReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyID, msg.ID),
			sdk.NewAttribute(ethbridge.AttributeKeyClaimType, msg.ClaimType.String()),
		),
		sdk.NewEvent(
			ethbridge.EventTypeProphecyStatus,
			sdk.NewAttribute(ethbridge.AttributeKeyStatus, status.Text.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBurn(
	ctx sdk.Context, cdc *codec.Codec, nftKeeper types.NFTKeeper,
	bridgeKeeper Keeper, msg MsgBurnNFT,
) (*sdk.Result, error) {

	if err := bridgeKeeper.ProcessBurn(ctx, msg.CosmosSender, msg.Denom, msg.ID); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			ethbridge.EventTypeBurn,
			sdk.NewAttribute(ethbridge.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(ethbridge.AttributeKeyTokenContract, msg.TokenContract.String()),
			sdk.NewAttribute(ethbridge.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(ethbridge.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyID, msg.ID),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}

func handleMsgLock(
	ctx sdk.Context, cdc *codec.Codec, nftKeeper types.NFTKeeper,
	bridgeKeeper Keeper, msg MsgLockNFT,
) (*sdk.Result, error) {

	if err := bridgeKeeper.ProcessLock(ctx, msg.CosmosSender, msg.Denom, msg.ID); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			ethbridge.EventTypeLock,
			sdk.NewAttribute(ethbridge.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(ethbridge.AttributeKeyTokenContract, msg.TokenContract.String()),
			sdk.NewAttribute(ethbridge.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(ethbridge.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyID, msg.ID),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}
