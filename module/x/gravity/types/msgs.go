package types

import (
	"fmt"

	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

var (
	_ sdk.Msg = &MsgDelegateKey{}
	_ sdk.Msg = &MsgSubmitClaim{}
	_ sdk.Msg = &MsgSubmitConfirm{}
	_ sdk.Msg = &MsgSendToEth{}
	_ sdk.Msg = &MsgRequestBatch{}
)

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSetDelegateKeys(val sdk.ValAddress, oper sdk.AccAddress, eth string) *MsgDelegateKey {
	return &MsgDelegateKey{
		Validator:    val.String(),
		Orchestrator: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKey) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKey) Type() string { return "set_operator_address" }

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKey) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgDelegateKey) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// GetSignBytes encodes the message for signing
func (msg *MsgDelegateKey) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}

// NewMsgSendToEth returns a new msgSendToEth
func NewMsgSendToEth(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEth {
	return &MsgSendToEth{
		Sender:    sender.String(),
		EthDest:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
	}
}

// Route should return the name of the module
func (msg *MsgSendToEth) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSendToEth) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg *MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
	if msg.Amount.Denom != msg.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := ValidateEthAddress(msg.EthDest); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgSendToEth) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg *MsgSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(orchestrator sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		Sender: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgCancelSendToEth(val sdk.ValAddress, id uint64) *MsgCancelSendToEth {
	return &MsgCancelSendToEth{
		TransactionId: id,
	}
}

// Route should return the name of the module
func (msg *MsgCancelSendToEth) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCancelSendToEth) Type() string { return "cancel_send_to_eth" }

// ValidateBasic performs stateless checks
func (msg *MsgCancelSendToEth) ValidateBasic() (err error) {
	_, err = sdk.ValAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelSendToEth) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}

// GetSigners defines whose signature is required
func (msg *MsgCancelSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSubmitConfirm(confirm *types.Any, signer string) *MsgSubmitConfirm {
	return &MsgSubmitConfirm{
		Confirm: confirm,
		Signer:  signer,
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
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgSubmitConfirm) GetConfirm() Confirm {
	confirm, ok := msg.Confirm.GetCachedValue().(Confirm)
	if !ok {
		return nil
	}
	return confirm
}

func (msg *MsgSubmitConfirm) SetConfirm(confirm Confirm) error {
	message, ok := confirm.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(message)
	if err != nil {
		return err
	}
	msg.Confirm = any
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgSubmitConfirm) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var claim EthereumClaim
	return unpacker.UnpackAny(msg.Confirm, &claim)
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitConfirm) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSubmitClaim(claim *types.Any, signer string) *MsgSubmitClaim {
	return &MsgSubmitClaim{
		Claim:  claim,
		Signer: signer,
	}
}

// Route should return the name of the module
func (msg *MsgSubmitClaim) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitClaim) Type() string { return "submit_claim" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitClaim) ValidateBasic() (err error) {
	_, err = sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgSubmitClaim) GetClaim() EthereumClaim {
	content, ok := m.Claim.GetCachedValue().(EthereumClaim)
	if !ok {
		return nil
	}
	return content
}

func (m *MsgSubmitClaim) SetClaim(claim EthereumClaim) error {
	msg, ok := claim.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Claim = any
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitClaim) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var claim EthereumClaim
	return unpacker.UnpackAny(m.Claim, &claim)
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitClaim) GetSignBytes() []byte {
	panic("Gravity messages do not support amino")
}
