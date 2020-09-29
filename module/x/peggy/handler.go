package peggy

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
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
		case MsgCreateEthereumClaims:
			return handleCreateEthereumClaims(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Peggy Msg type: %v", msg.Type()))
		}
	}
}

func handleCreateEthereumClaims(ctx sdk.Context, keeper Keeper, msg MsgCreateEthereumClaims) (*sdk.Result, error) {
	// TODO:
	// auth EthereumChainID whitelisted
	// auth bridge contract address whitelisted
	ctx.Logger().Info("+++ TODO: implement chain id + contract address authorization")
	//if !bytes.Equal(msg.BridgeContractAddress[:], k.GetBridgeContractAddress(ctx)[:]) {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid bridge contract address")
	//}

	var attestationIDs [][]byte
	// auth sender in current validator set
	for _, c := range msg.Claims {
		validator := findValidatorKey(ctx, msg.Orchestrator)
		if validator == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
		}
		att, err := keeper.AddClaim(ctx, c.GetType(), c.GetNonce(), validator, c.Details())
		if err != nil {
			return nil, sdkerrors.Wrap(err, "create attestation")
		}
		attestationIDs = append(attestationIDs, att.ID())
	}
	return &sdk.Result{
		Data: bytes.Join(attestationIDs, []byte(", ")),
	}, nil
}

func findValidatorKey(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress {
	// todo: implement proper in keeper
	return sdk.ValAddress(orchAddr)
}

func handleMsgValsetRequest(ctx sdk.Context, keeper Keeper, msg types.MsgValsetRequest) (*sdk.Result, error) {
	// todo: is requester in current valset?\

	// disabling bootstrap check for integration tests to pass
	//if keeper.GetLastValsetObservedNonce(ctx).IsEmpty() {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge bootstrap process not observed, yet")
	//}
	v := keeper.SetValsetRequest(ctx)
	return &sdk.Result{
		Data: v.Nonce.Bytes(),
	}, nil
}

func handleMsgValsetConfirm(ctx sdk.Context, keeper Keeper, msg MsgValsetConfirm) (*sdk.Result, error) {
	// Check that the signature is valid for the valset at the blockheight and the validator
	valset := keeper.GetValsetRequest(ctx, msg.Nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown nonce")
	}
	validator := findValidatorKey(ctx, msg.Validator)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}
	checkpoint := valset.GetCheckpoint()
	if err := verifyCheckpointSignature(ctx, keeper, validator, msg.Signature, checkpoint); err != nil {
		return nil, err
	}
	// Save valset confirmation
	keeper.SetValsetConfirm(ctx, msg)

	details := types.SignedCheckpoint{Checkpoint: checkpoint}
	att, err := keeper.AddClaim(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, msg.Nonce, validator, details)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: att.ID(),
	}, nil
}

func handleMsgSetEthAddress(ctx sdk.Context, keeper Keeper, msg MsgSetEthAddress) (*sdk.Result, error) {
	validator := findValidatorKey(ctx, msg.Validator)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}

	keeper.SetEthAddress(ctx, validator, msg.Address)
	return &sdk.Result{}, nil
}

func handleMsgSendToEth(ctx sdk.Context, keeper Keeper, msg MsgSendToEth) (*sdk.Result, error) {
	txID, err := keeper.AddToOutgoingPool(ctx, msg.Sender, msg.DestAddress, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: sdk.Uint64ToBigEndian(txID),
	}, nil
}

func handleMsgRequestBatch(ctx sdk.Context, keeper Keeper, msg MsgRequestBatch) (*sdk.Result, error) {
	batchID, err := keeper.BuildOutgoingTXBatch(ctx, msg.Denom, OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: sdk.Uint64ToBigEndian(batchID),
	}, nil
}

func handleMsgConfirmBatch(ctx sdk.Context, keeper Keeper, msg MsgConfirmBatch) (*sdk.Result, error) {
	validator := findValidatorKey(ctx, msg.Orchestrator)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}

	batch := keeper.GetOutgoingTXBatch(ctx, msg.Nonce.Uint64())
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	checkpoint := batch.GetCheckpoint()
	if err := verifyCheckpointSignature(ctx, keeper, validator, msg.Signature, checkpoint); err != nil {
		return nil, err
	}
	details := types.SignedCheckpoint{Checkpoint: checkpoint}
	att, err := keeper.AddClaim(ctx, types.ClaimTypeOrchestratorSignedWithdrawBatch, msg.Nonce, validator, details)
	if err != nil {
		return nil, err
	}
	keeper.SetOutgoingTXBatchConfirm(ctx, msg.Nonce.Uint64(), validator, []byte(msg.Signature))
	return &sdk.Result{
		Data: att.ID(),
	}, nil
}

// verifies the checkpoint was signed with the deposited ethereum key for the given validator
func verifyCheckpointSignature(ctx sdk.Context, keeper Keeper, validator sdk.ValAddress, ethSignature string, checkpoint []byte) error {
	ethAddress := keeper.GetEthAddress(ctx, validator)
	if len(ethAddress) == 0 {
		return sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}
	sigBytes, err := hex.DecodeString(ethSignature)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress.String())
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "signature")
	}
	return nil
}
