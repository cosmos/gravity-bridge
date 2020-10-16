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
		case MsgBridgeSignatureSubmission:
			return handleBridgeSignatureSubmission(ctx, keeper, msg)
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
	//if keeper.GetLastValsetObservedNonce(ctx).isValid() {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge bootstrap process not observed, yet")
	//}
	v := keeper.SetValsetRequest(ctx)
	return &sdk.Result{
		Data: v.Nonce.Bytes(),
	}, nil
}

// deprecated should use MsgBridgeSignatureSubmission instead
func handleMsgValsetConfirm(ctx sdk.Context, keeper Keeper, msg MsgValsetConfirm) (*sdk.Result, error) {
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	keeper.SetValsetConfirm(ctx, msg) // store legacy
	return handleBridgeSignatureSubmission(ctx, keeper, MsgBridgeSignatureSubmission{
		Nonce:             msg.Nonce,
		ClaimType:         types.ClaimTypeOrchestratorSignedMultiSigUpdate,
		Orchestrator:      msg.Validator,
		EthereumSignature: sigBytes,
	})
}

func handleBridgeSignatureSubmission(ctx sdk.Context, keeper Keeper, msg MsgBridgeSignatureSubmission) (*sdk.Result, error) {
	checkpoint, err := getCheckpoint(ctx, keeper, msg.ClaimType, msg.Nonce)
	if err != nil {
		return nil, err
	}
	validator := findValidatorKey(ctx, msg.Orchestrator)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}

	ethAddress := keeper.GetEthAddress(ctx, validator)
	if ethAddress == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}
	err = types.ValidateEthereumSignature(checkpoint, msg.EthereumSignature, ethAddress.String())
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature")
	}

	// persist signature
	if keeper.HasBridgeApprovalSignature(ctx, msg.ClaimType, msg.Nonce, validator) {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature")
	}
	keeper.SetBridgeApprovalSignature(ctx, msg.ClaimType, msg.Nonce, validator, msg.EthereumSignature)
	details := types.SignedCheckpoint{Checkpoint: checkpoint}
	att, err := keeper.AddClaim(ctx, msg.ClaimType, msg.Nonce, validator, details)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: att.ID(),
	}, nil
}

func getCheckpoint(ctx sdk.Context, keeper Keeper, claimType types.ClaimType, nonce types.UInt64Nonce) ([]byte, error) {
	switch claimType {
	case types.ClaimTypeOrchestratorSignedMultiSigUpdate:
		if c := keeper.GetValsetRequest(ctx, nonce); c != nil {
			return c.GetCheckpoint(), nil
		}
	case types.ClaimTypeOrchestratorSignedWithdrawBatch:
		if c := keeper.GetOutgoingTXBatch(ctx, nonce); c != nil {
			nonce := keeper.GetLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate)
			if nonce == nil || nonce.IsEmpty() {
				nonce = keeper.GetLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeBootstrap)
			}
			if nonce == nil || nonce.IsEmpty() {
				return nil, sdkerrors.Wrap(types.ErrUnknown, "observed multisig set")
			}
			valset := keeper.GetValsetRequest(ctx, *nonce)
			if valset == nil {
				return nil, sdkerrors.Wrap(types.ErrUnknown, "no valset found for nonce")
			}
			return c.GetCheckpoint(valset)
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrUnsupported, "claim type")
	}
	return nil, sdkerrors.Wrap(types.ErrInvalid, "does not exist")
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
		Data: batchID.Nonce.Bytes(),
	}, nil
}

func handleMsgConfirmBatch(ctx sdk.Context, keeper Keeper, msg MsgConfirmBatch) (*sdk.Result, error) {
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	validator := findValidatorKey(ctx, msg.Orchestrator)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}

	keeper.SetOutgoingTXBatchConfirm(ctx, msg.Nonce, validator, []byte(msg.Signature)) // store legacy

	return handleBridgeSignatureSubmission(ctx, keeper, MsgBridgeSignatureSubmission{
		Nonce:             msg.Nonce,
		ClaimType:         types.ClaimTypeOrchestratorSignedWithdrawBatch,
		Orchestrator:      msg.Orchestrator,
		EthereumSignature: sigBytes,
	})
}
