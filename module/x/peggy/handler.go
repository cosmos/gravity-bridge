package peggy

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "Peggy" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgSetEthAddress:
			return handleMsgSetEthAddress(ctx, keeper, msg)
		case *types.MsgValsetConfirm:
			return handleMsgConfirmValset(ctx, keeper, msg)
		case *types.MsgValsetRequest:
			return handleMsgValsetRequest(ctx, keeper, msg)
		case *types.MsgSendToEth:
			return handleMsgSendToEth(ctx, keeper, msg)
		case *types.MsgRequestBatch:
			return handleMsgRequestBatch(ctx, keeper, msg)
		case *types.MsgConfirmBatch:
			return handleMsgConfirmBatch(ctx, keeper, msg)
		case *types.MsgDepositClaim:
			return handleDepositClaim(ctx, keeper, msg)
		case *types.MsgWithdrawClaim:
			return handleWithdrawClaim(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Peggy Msg type: %v", msg.Type()))
		}
	}
}

func handleDepositClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgDepositClaim) (*sdk.Result, error) {
	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data: bytes.Join(attestationIDs, []byte(", ")),
	}, nil
}

func handleWithdrawClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgWithdrawClaim) (*sdk.Result, error) {
	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data: bytes.Join(attestationIDs, []byte(", ")),
	}, nil
}

func findValidatorKey(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress {
	// todo: implement proper in keeper
	// TODO: do we want ValAddress or do we want the AccAddress for the validator?
	// this is a v important question for encoding
	return sdk.ValAddress(orchAddr)
}

func handleMsgValsetRequest(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgValsetRequest) (*sdk.Result, error) {
	// todo: is requester in current valset?\

	// disabling bootstrap check for integration tests to pass
	//if keeper.GetLastValsetObservedNonce(ctx).isValid() {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge bootstrap process not observed, yet")
	//}
	v := keeper.SetValsetRequest(ctx)
	return &sdk.Result{
		Data: types.UInt64Bytes(v.Nonce),
	}, nil
}

// This function takes in a signature submitted by a validator's Eth Signer
func handleMsgConfirmBatch(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgConfirmBatch) (*sdk.Result, error) {

	batch := keeper.GetOutgoingTXBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	peggyID := keeper.GetPeggyID(ctx)
	checkpoint, err := batch.GetCheckpoint(peggyID)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)
	validator := findValidatorKey(ctx, valaddr)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}

	ethAddress := keeper.GetEthAddress(ctx, sdk.AccAddress(validator))
	if ethAddress == "" {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}
	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// check if we already have this confirm
	if keeper.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, valaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := keeper.SetBatchConfirm(ctx, msg)
	return &sdk.Result{
		Data: key,
	}, nil
}

// This function takes in a signature submitted by a validator's Eth Signer
func handleMsgConfirmValset(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgValsetConfirm) (*sdk.Result, error) {

	valset := keeper.GetValsetRequest(ctx, msg.Nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
	}

	peggyID := keeper.GetPeggyID(ctx)
	checkpoint := valset.GetCheckpoint(peggyID)

	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)
	validator := findValidatorKey(ctx, valaddr)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	}

	ethAddress := keeper.GetEthAddress(ctx, sdk.AccAddress(validator))
	if ethAddress == "" {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}
	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// persist signature
	if keeper.GetValsetConfirm(ctx, msg.Nonce, valaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := keeper.SetValsetConfirm(ctx, *msg)
	return &sdk.Result{
		Data: key,
	}, nil
}

func handleMsgSetEthAddress(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgSetEthAddress) (*sdk.Result, error) {
	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)
	validator := findValidatorKey(ctx, valaddr)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}

	keeper.SetEthAddress(ctx, sdk.AccAddress(validator), msg.Address)
	return &sdk.Result{}, nil
}

func handleMsgSendToEth(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgSendToEth) (*sdk.Result, error) {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	txID, err := keeper.PushToOutgoingPool(ctx, sender, msg.EthDest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: sdk.Uint64ToBigEndian(txID),
	}, nil
}

func handleMsgRequestBatch(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRequestBatch) (*sdk.Result, error) {
	// ensure that peggy denom is valid
	ec, err := types.ERC20FromPeggyCoin(sdk.NewInt64Coin(msg.Denom, 0))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid denom: %s", err)
	}

	batchID, err := k.BuildOutgoingTXBatch(ctx, ec.Contract, keeper.OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data: types.UInt64Bytes(batchID.BatchNonce),
	}, nil
}
