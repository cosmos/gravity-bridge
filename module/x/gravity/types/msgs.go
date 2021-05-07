package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgDelegateKeys{}
	_ sdk.Msg = &MsgSendToEthereum{}
	_ sdk.Msg = &MsgRequestBatchTx{}
	_ sdk.Msg = &MsgSubmitEthereumEvent{}
	_ sdk.Msg = &MsgSubmitEthereumSignature{}
	_ sdk.Msg = &MsgCancelSendToEthereum{}
)

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSetOrchestratorAddress(val sdk.ValAddress, oper sdk.AccAddress, eth string) *MsgDelegateKeys {
	return &MsgDelegateKeys{
		ValidatorAddress:    val.String(),
		OrchestratorAddress: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKeys) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKeys) Type() string { return "delegate_keys"}

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKeys) ValidateBasic() (err error) {
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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


// Route should return the name of the module
func (msg *MsgSubmitEthereumSignature) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEthereumSignature) Type() string { return "submit_ethereum_signature" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEthereumSignature) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event, err := UnpackSignature(msg.Signature)
	if err != nil {
		return err
	}

	return event.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEthereumSignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEthereumSignature) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToEthereum returns a new MsgSendToEthereum
func NewMsgSendToEthereum(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEthereum {
	return &MsgSendToEthereum{
		Sender:    sender.String(),
		EthRecipient:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
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
func (msg MsgSendToEthereum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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
func NewMsgRequestBatchTx(orchestrator sdk.AccAddress) *MsgRequestBatchTx {
	return &MsgRequestBatchTx{
		Signer: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatchTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatchTx) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatchTx) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatchTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatchTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}
