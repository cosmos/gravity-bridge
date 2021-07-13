package gravity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgSendToEthereum:
			res, err := msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCancelSendToEthereum:
			res, err := msgServer.CancelSendToEthereum(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRequestBatchTx:
			res, err := msgServer.RequestBatchTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSubmitEthereumTxConfirmation:
			res, err := msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSubmitEthereumEvent:
			res, err := msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDelegateKeys:
			res, err := msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Gravity Msg type: %v", msg.Type()))
		}
	}
}
