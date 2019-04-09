package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMakeBridgeClaim:
			return handleMsgMakeBridgeClaim(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to make a bridge claim
func handleMsgMakeBridgeClaim(ctx sdk.Context, keeper Keeper, msg MsgMakeBridgeClaim) sdk.Result {
	//do it
	return sdk.Result{} // return
}
