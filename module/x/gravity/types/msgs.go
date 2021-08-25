package types

import (
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgDelegateKeys{}
	_ sdk.Msg = &MsgSendToEthereum{}
	_ sdk.Msg = &MsgCancelSendToEthereum{}
	_ sdk.Msg = &MsgRequestBatchTx{}
	_ sdk.Msg = &MsgSubmitEthereumEvent{}
	_ sdk.Msg = &MsgSubmitEthereumTxConfirmation{}

	_ cdctypes.UnpackInterfacesMessage = &MsgSubmitEthereumEvent{}
	_ cdctypes.UnpackInterfacesMessage = &MsgSubmitEthereumTxConfirmation{}
	_ cdctypes.UnpackInterfacesMessage = &EthereumEventVoteRecord{}
)

// NewMsgDelegateKeys returns a reference to a new MsgDelegateKeys.
func NewMsgDelegateKeys(val sdk.ValAddress, orchAddr sdk.AccAddress, ethAddr string, ethSig []byte) *MsgDelegateKeys {
	return &MsgDelegateKeys{
		ValidatorAddress:    val.String(),
		OrchestratorAddress: orchAddr.String(),
		EthereumAddress:     ethAddr,
		EthSignature:        ethSig,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKeys) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKeys) Type() string { return "delegate_keys" }

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKeys) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress)
	}
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if !common.IsHexAddress(msg.EthereumAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "ethereum address")
	}
	if len(msg.EthSignature) == 0 {
		return ErrEmptyEthSig
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgDelegateKeys) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgDelegateKeys) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// Route should return the name of the module
func (msg *MsgSubmitEthereumEvent) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEthereumEvent) Type() string { return "submit_ethereum_event" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEthereumEvent) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event, err := UnpackEvent(msg.Event)
	if err != nil {
		return err
	}
	return event.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEthereumEvent) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEthereumEvent) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgSubmitEthereumEvent) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var event EthereumEvent
	return unpacker.UnpackAny(msg.Event, &event)
}

// Route should return the name of the module
func (msg *MsgSubmitEthereumTxConfirmation) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEthereumTxConfirmation) Type() string { return "submit_ethereum_signature" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEthereumTxConfirmation) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event, err := UnpackConfirmation(msg.Confirmation)

	if err != nil {
		return err
	}

	return event.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEthereumTxConfirmation) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEthereumTxConfirmation) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgSubmitEthereumTxConfirmation) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var sig EthereumTxConfirmation
	return unpacker.UnpackAny(msg.Confirmation, &sig)
}

// NewMsgSendToEthereum returns a new MsgSendToEthereum
func NewMsgSendToEthereum(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEthereum {
	return &MsgSendToEthereum{
		Sender:            sender.String(),
		EthereumRecipient: destAddress,
		Amount:            send,
		BridgeFee:         bridgeFee,
	}
}

// Route should return the name of the module
func (msg MsgSendToEthereum) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToEthereum) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEthereum) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
	// this check is VERY IMPORTANT
	if msg.Amount.Denom != msg.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins,
			fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if !common.IsHexAddress(msg.EthereumRecipient) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "ethereum address")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToEthereum) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEthereum) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatchTx returns a new msgRequestBatch
func NewMsgRequestBatchTx(denom string, signer sdk.AccAddress) *MsgRequestBatchTx {
	return &MsgRequestBatchTx{
		Denom:  denom,
		Signer: signer.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatchTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatchTx) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatchTx) ValidateBasic() error {
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return sdkerrors.Wrap(err, "denom is invalid")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatchTx) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatchTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgCancelSendToEthereum returns a new MsgCancelSendToEthereum
func NewMsgCancelSendToEthereum(id uint64, orchestrator sdk.AccAddress) *MsgCancelSendToEthereum {
	return &MsgCancelSendToEthereum{
		Id:     id,
		Sender: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgCancelSendToEthereum) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCancelSendToEthereum) Type() string { return "cancel_send_to_ethereum" }

// ValidateBasic performs stateless checks
func (msg MsgCancelSendToEthereum) ValidateBasic() error {
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalid, "Id cannot be 0")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCancelSendToEthereum) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgCancelSendToEthereum) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}
