package peggy

import (
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/althea-net/peggy/module/x/peggy/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgSetEthAddress:
			return handleMsgSetEthAddress(ctx, keeper, msg)
		case MsgValsetConfirm:
			return handleMsgValsetConfirm(ctx, keeper, msg)
		case MsgValsetRequest:
			return handleMsgValsetRequest(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}

func handleMsgValsetRequest(ctx sdk.Context, keeper Keeper, msg types.MsgValsetRequest) (*sdk.Result, error) {
	keeper.SetValsetRequest(ctx)
	return &sdk.Result{}, nil
}

func handleMsgValsetConfirm(ctx sdk.Context, keeper Keeper, msg MsgValsetConfirm) (*sdk.Result, error) {
	// Check that the signature is valid for the valset at the blockheight and the validator
	valset := keeper.GetValsetRequest(ctx, msg.Nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown nonce")
	}

	checkpoint := valset.GetCheckpoint()
	ethAddress := keeper.GetEthAddress(ctx, msg.Validator)
	if len(ethAddress) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty eth address")
	}

	err := utils.ValidateEthSig(checkpoint, msg.Signature, ethAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Save valset confirmation
	keeper.SetValsetConfirm(ctx, msg)
	return &sdk.Result{}, nil
}

func handleMsgSetEthAddress(ctx sdk.Context, keeper Keeper, msg MsgSetEthAddress) (*sdk.Result, error) {
	keeper.SetEthAddress(ctx, msg.Validator, msg.Address)
	return &sdk.Result{}, nil
}
