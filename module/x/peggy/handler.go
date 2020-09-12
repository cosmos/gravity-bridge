package peggy

import (
	"encoding/hex"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/althea-net/peggy/module/x/peggy/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "Peggy" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgSetEthAddress:
			return handleMsgSetEthAddress(ctx, keeper, msg)
		case MsgValsetConfirm:
			return handleMsgValsetConfirm(ctx, keeper, msg)
		case MsgValsetRequest:
			return handleMsgValsetRequest(ctx, keeper, msg)
		case MsgSendToEth:
			return handleMsgSendToEth(ctx, keeper, msg)
		case MsgRequestBatch:
			return handleMsgRequestBatch(ctx, keeper, msg)
		case MsgConfirmBatch:
			return handleMsgConfirmBatch(ctx, keeper, msg)
		case MsgBatchInChain:
			return handleMsgBatchInChain(ctx, keeper, msg)
		case MsgEthDeposit:
			return handleMsgEthDeposit(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Peggy Msg type: %v", msg.Type()))
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

	sigBytes, hexErr := hex.DecodeString(msg.Signature)
	if hexErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Signature hex decoding error")
	}
	err := utils.ValidateEthSig(checkpoint, sigBytes, ethAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Failed to validate Checkpoint Sig")
	}

	// Save valset confirmation
	keeper.SetValsetConfirm(ctx, msg)
	return &sdk.Result{}, nil
}

func handleMsgSetEthAddress(ctx sdk.Context, keeper Keeper, msg MsgSetEthAddress) (*sdk.Result, error) {
	keeper.SetEthAddress(ctx, msg.Validator, msg.Address)
	return &sdk.Result{}, nil
}

func handleMsgSendToEth(ctx sdk.Context, keeper Keeper, msg MsgSendToEth) (*sdk.Result, error) {
	// TODO add this transcation to the Peggy Tx Pool
	return &sdk.Result{}, nil
}

func handleMsgRequestBatch(ctx sdk.Context, keeper Keeper, msg MsgRequestBatch) (*sdk.Result, error) {
	// TODO perform the batch creation process here, including pulling transactions out of
	// the Peggy Tx Pool and bundling them into transactions
	return &sdk.Result{}, nil
}

func handleMsgConfirmBatch(ctx sdk.Context, keeper Keeper, msg MsgConfirmBatch) (*sdk.Result, error) {
	// TODO add batch confirmation to the store, and if this confirmation means the batch counts as
	// `observed` (confirmations from 66% of the active voting power exist as of this block) then consider
	// the batch completed.
	return &sdk.Result{}, nil
}

func handleMsgBatchInChain(ctx sdk.Context, keeper Keeper, msg MsgBatchInChain) (*sdk.Result, error) {
	// TODO add batch confirmation to the store, and if this confirmation means the batch counts as
	// `observed` (confirmations from 66% of the active voting power exist as of this block) then consider
	// the batch completed.
	return &sdk.Result{}, nil
}

func handleMsgEthDeposit(ctx sdk.Context, keeper Keeper, msg MsgEthDeposit) (*sdk.Result, error) {
	// TODO issue tokens from the store of the appropriate denom once this deposit counts as `observed`
	return &sdk.Result{}, nil
}
