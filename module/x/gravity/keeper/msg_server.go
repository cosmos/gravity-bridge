package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetDelegateKeys(c context.Context, msg *types.MsgDelegateKeys) (*types.MsgDelegateKeysResponse, error) {
	// ensure that this passes validation
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	val, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	orch, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)

	// ensure that the validator exists
	if k.Keeper.StakingKeeper.Validator(ctx, val) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	}

	// TODO consider impact of maliciously setting duplicate delegate
	// addresses since no signatures from the private keys of these addresses
	// are required for this message it could be sent in a hostile way.

	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, val, orch)
	// set the ethereum address
	k.SetEthAddress(ctx, val, common.HexToAddress(msg.EthereumAddress))
	// TODO: set the third index from eth addr -> orch addr so we can query all of them

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeySetOperatorAddr, orch.String()),
		),
	)

	return &types.MsgDelegateKeysResponse{}, nil

}

// SubmitEthereumSignature handles MsgSubmitEthereumSignature
func (k msgServer) SubmitEthereumSignature(c context.Context, msg *types.MsgSubmitEthereumSignature) (*types.MsgSubmitEthereumSignatureResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signature, err := types.UnpackSignature(msg.Signature)
	if err != nil {
		return nil, err
	}

	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	otx := k.GetOutgoingTx(ctx, signature.GetStoreIndex())
	if otx == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find outgoing tx")
	}

	gravityID := k.GetGravityID(ctx)
	checkpoint := otx.GetCheckpoint([]byte(gravityID))

	ethAddress := k.GetEthAddress(ctx, val)
	if ethAddress == "" || ethAddress != signature.GetSigner().String() {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	}

	if err = types.ValidateEthereumSignature(checkpoint, signature.GetSignature(), ethAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s", ethAddress, gravityID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// TODO: should validators be able to overwrite their signatures?
	if k.GetEthereumSignature(ctx, signature.GetStoreIndex(), val) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}

	key := k.SetEthereumSignature(ctx, signature, val)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyEthereumSignatureKey, string(key)),
		),
	)

	return &types.MsgSubmitEthereumSignatureResponse{}, nil
}

// func (k Keeper) ValidateEthereumSignature

// SubmitEthereumEvent handles MsgSubmitEthereumEvent
func (k msgServer) SubmitEthereumEvent(c context.Context, msg *types.MsgSubmitEthereumEvent) (*types.MsgSubmitEthereumEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	event, err := types.UnpackEvent(msg.Event)
	if err != nil {
		return nil, err
	}

	// return an error if the validator isn't in the active set
	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	// Add the claim to the store
	_, err = k.RecordEventVote(ctx, event, val)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create event vote record")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, fmt.Sprintf("%T", event)),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordID, string(types.GetEthereumEventVoteRecordKey(event.GetNonce(), event.Hash()))),
		),
	)

	return &types.MsgSubmitEthereumEventResponse{}, nil
}

// SendToEthereum handles MsgSendToEthereum
func (k msgServer) SendToEthereum(c context.Context, msg *types.MsgSendToEthereum) (*types.MsgSendToEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	txID, err := k.AddToOutgoingPool(ctx, sender, msg.EthereumRecipient, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)

	return &types.MsgSendToEthereumResponse{}, nil
}

// RequestBatchTx handles MsgRequestBatchTx
func (k msgServer) RequestBatchTx(c context.Context, msg *types.MsgRequestBatchTx) (*types.MsgRequestBatchTxResponse, error) {
	// TODO: limit this to only orchestrators and validators?
	ctx := sdk.UnwrapSDKContext(c)

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID, err := k.BuildBatchTx(ctx, tokenContract, BatchTxSize)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyContract, tokenContract),
			sdk.NewAttribute(types.AttributeKeyBatchNonce, fmt.Sprint(batchID.Nonce)),
		),
	)

	return &types.MsgRequestBatchTxResponse{}, nil
}

func (k msgServer) CancelSendToEthereum(c context.Context, msg *types.MsgCancelSendToEthereum) (*types.MsgCancelSendToEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.RemoveFromOutgoingPoolAndRefund(ctx, msg.Id, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(msg.Id)),
		),
	)

	return &types.MsgCancelSendToEthereumResponse{}, nil
}

func (k Keeper) getSignerValidator(ctx sdk.Context, signerString string) (sdk.ValAddress, error) {
	signer, _ := sdk.AccAddressFromBech32(signerString)

	var validatorI stakingtypes.ValidatorI
	validator := k.GetOrchestratorValidator(ctx, signer)
	if validator == nil {
		validatorI = k.StakingKeeper.Validator(ctx, sdk.ValAddress(signer))
	} else {
		validatorI = k.StakingKeeper.Validator(ctx, validator)
	}

	if validatorI == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "not orchestrator or validator")
	} else if !validatorI.IsBonded() {
		return nil, sdkerrors.Wrap(types.ErrUnbonded, fmt.Sprintf("validator: %s", validatorI.GetOperator()))
	}

	return validatorI.GetOperator(), nil
}
