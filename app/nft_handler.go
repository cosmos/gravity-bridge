package app

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
)

// OverrideNFTModule overrides the NFT module for custom handlers
type OverrideNFTModule struct {
	nft.AppModule
	k nft.Keeper
}

// NewHandler module handler for the OerrideNFTModule
func (am OverrideNFTModule) NewHandler() sdk.Handler {
	return CustomNFTHandler(am.k)
}

// NewOverrideNFTModule generates a new NFT Module
func NewOverrideNFTModule(appModule nft.AppModule, keeper nft.Keeper) OverrideNFTModule {
	return OverrideNFTModule{
		AppModule: appModule,
		k:         keeper,
	}
}

// CustomNFTHandler routes the messages to the handlers
func CustomNFTHandler(k nft.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case nft.MsgTransferNFT:
			return handleMsgTransferNFTCustom(ctx, msg, k)
		// case nft.MsgEditNFTMetadata:
		// 	return nft.HandleMsgEditNFTMetadata(ctx, msg, k)
		// case nft.MsgMintNFT:
		// 	return HandleMsgMintNFTCustom(ctx, msg, k)
		// case nft.MsgBurnNFT:
		// 	return nft.HandleMsgBurnNFT(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,
				fmt.Sprintf("unrecognized %s message type: %T", nft.ModuleName, msg))
		}
	}
}

// handleMsgTransferNFTCustom handles MsgTransferNFT
func handleMsgTransferNFTCustom(ctx sdk.Context, msg nft.MsgTransferNFT, k nft.Keeper,
) (*sdk.Result, error) {

	foundNFT, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if foundNFT.GetOwner().Equals(msg.Sender) {
		return nft.HandleMsgTransferNFT(ctx, msg, k)
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
		fmt.Sprintf("can't transfer NFT you don't own: %T", msg))
}
