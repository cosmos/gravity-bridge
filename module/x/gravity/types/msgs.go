package types

import (
	"fmt"

	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ sdk.Msg = &MsgDelegateKey{}
	_ sdk.Msg = &MsgSubmitEvent{}
	_ sdk.Msg = &MsgSubmitConfirm{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgTransfer{}
	_ sdk.Msg = &MsgCancelTransfer{}
)

// NewMsgDelegateKey returns a new msgSetOrchestratorAddress
func NewMsgDelegateKey(validatorAddr sdk.ValAddress, operatorAddr sdk.AccAddress, ethereumAddr common.Address) *MsgDelegateKey {
	return &MsgDelegateKey{
		ValidatorAddress:    validatorAddr.String(),
		OrchestratorAddress: operatorAddr.String(),
		EthAddress:          ethereumAddr.String(),
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKey) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKey) Type() string { return "delegate_key" }

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKey) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress)
	}
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgDelegateKey) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// GetSignBytes encodes the message for signing
func (msg *MsgDelegateKey) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}

// NewMsgTransfer returns a new MsgTransfer
func NewMsgTransfer(sender sdk.AccAddress, ethRecipientAddr common.Address, amount, fee sdk.Coin) *MsgTransfer {
	return &MsgTransfer{
		Sender:       sender.String(),
		EthRecipient: ethRecipientAddr.String(),
		Amount:       amount,
		BridgeFee:    fee,
	}
}

// Route should return the name of the module
func (msg MsgTransfer) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTransfer) Type() string { return "transfer" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgTransfer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
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
	if err := ValidateEthAddress(msg.EthRecipient); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTransfer) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(denom string, orchestrator sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		OrchestratorAddress: orchestrator.String(),
		Denom:               denom,
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}

	// NOTE: this only supports gravity denoms (plain or with the 'gravity/' namespace)
	// consider only passing ethereum addresses?
	return ValidateGravityDenom(msg.Denom)
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgCancelTransfer returns a new MsgCancelTransfer
func NewMsgCancelTransfer(txID tmbytes.HexBytes, sender sdk.AccAddress) *MsgCancelTransfer {
	return &MsgCancelTransfer{
		TxID:   txID,
		Sender: sender.String(),
	}
}

// Route should return the name of the module
func (msg *MsgCancelTransfer) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCancelTransfer) Type() string { return "cancel_transfer" }

// ValidateBasic performs stateless checks
func (msg *MsgCancelTransfer) ValidateBasic() (err error) {
	if len(msg.TxID) == 0 {
		return fmt.Errorf("tx id cannot be empty")
	}

	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelTransfer) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg *MsgCancelTransfer) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// NewMsgSubmitConfirm returns a new MsgSubmitConfirm
func NewMsgSubmitConfirm(confirm *types.Any, orchestratorAddr sdk.AccAddress) *MsgSubmitConfirm {
	return &MsgSubmitConfirm{
		Confirm: confirm,
		Signer:  orchestratorAddr.String(),
	}
}

// Route should return the name of the module
func (msg *MsgSubmitConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitConfirm) Type() string { return "submit_confirm" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitConfirm) ValidateBasic() (err error) {
	_, err = sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}

	confirm, err := UnpackConfirm(msg.Confirm)
	if err != nil {
		return err
	}

	return confirm.Validate()
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgSubmitConfirm) GetConfirm() Confirm {
	confirm, _ := UnpackConfirm(m.Confirm)
	return confirm
}

func (m *MsgSubmitConfirm) SetConfirm(confirm Confirm) error {
	any, err := PackConfirm(confirm)
	if err != nil {
		return err
	}

	m.Confirm = any
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitConfirm) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var confirm Confirm
	return unpacker.UnpackAny(m.Confirm, &confirm)
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitConfirm) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}

// NewMsgSubmitEvent returns a new MsgSubmitEvent
func NewMsgSubmitEvent(event *types.Any, signer sdk.AccAddress) *MsgSubmitEvent {
	return &MsgSubmitEvent{
		Event:  event,
		Signer: signer.String(),
	}
}

// Route should return the name of the module
func (msg *MsgSubmitEvent) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEvent) Type() string { return "submit_event" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEvent) ValidateBasic() (err error) {
	_, err = sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	event, err := UnpackEvent(msg.Event)
	if err != nil {
		return err
	}

	return event.Validate()
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgSubmitEvent) GetEvent() EthereumEvent {
	event, _ := UnpackEvent(m.Event)
	return event
}

func (m *MsgSubmitEvent) SetEvent(event EthereumEvent) error {
	any, err := PackEvent(event)
	if err != nil {
		return err
	}
	m.Event = any
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitEvent) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var event EthereumEvent
	return unpacker.UnpackAny(m.Event, &event)
}

// GetSignBytes encodes the message for signing
func (msg MsgSubmitEvent) GetSignBytes() []byte {
	panic("gravity messages do not support amino")
}
