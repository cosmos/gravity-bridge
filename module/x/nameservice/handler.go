package nameservice

import (
	"fmt"

	"github.com/althea-net/peggy/module/x/nameservice/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// case MsgSetName:
		// 	return handleMsgSetName(ctx, keeper, msg)
		// case MsgBuyName:
		// 	return handleMsgBuyName(ctx, keeper, msg)
		// case MsgDeleteName:
		// 	return handleMsgDeleteName(ctx, keeper, msg)
		case MsgSetEthAddress:
			return handleMsgSetEthAddress(ctx, keeper, msg)
		case MsgValsetConfirm:
			return handleMsgValsetConfirm(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}

func handleMsgValsetConfirm(ctx sdk.Context, keeper Keeper, msg MsgValsetConfirm) (*sdk.Result, error) {
	// Check that the signature is valid for the valset at the blockheight and the validator
	valset := keeper.GetValsetRequest(ctx, msg.Nonce)

	checkpoint := valset.GetCheckpoint()
	ethAddress := keeper.GetEthAddress(ctx, msg.Validator)

	err := utils.ValidateEthSig(checkpoint, msg.Signature, ethAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	// Save valset confirmation
	keeper.SetValsetConfirm(ctx, msg)
	return &sdk.Result{}, nil
}

func handleMsgSetEthAddress(ctx sdk.Context, keeper Keeper, msg MsgSetEthAddress) (*sdk.Result, error) {
	keeper.SetEthAddress(ctx, msg.Validator, msg.Address)
	return &sdk.Result{}, nil
}

// // Handle a message to set name
// func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) (*sdk.Result, error) {
// 	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
// 	}
// 	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
// 	return &sdk.Result{}, nil                // return
// }

// // Handle a message to buy name
// func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) (*sdk.Result, error) {
// 	// Checks if the the bid price is greater than the price paid by the current owner
// 	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
// 	}
// 	if keeper.HasOwner(ctx, msg.Name) {
// 		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
// 	keeper.SetPrice(ctx, msg.Name, msg.Bid)
// 	return &sdk.Result{}, nil
// }

// // Handle a message to delete name
// func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg MsgDeleteName) (*sdk.Result, error) {
// 	if !keeper.IsNamePresent(ctx, msg.Name) {
// 		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
// 	}
// 	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
// 	}

// 	keeper.DeleteWhois(ctx, msg.Name)
// 	return &sdk.Result{}, nil
// }
