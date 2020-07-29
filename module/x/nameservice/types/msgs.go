package types

import (
	"regexp"

	"github.com/althea-net/peggy/module/x/nameservice/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

// ValsetConfirm
// -------------
type MsgValsetConfirm struct {
	Nonce     int64          `json:"nonce"`
	Validator sdk.AccAddress `json:"validator"`
	Signature []byte         `json:"signature"`
}

func NewMsgValsetConfirm(nonce int64, validator sdk.AccAddress, signature []byte) MsgValsetConfirm {
	return MsgValsetConfirm{
		Nonce:     nonce,
		Validator: validator,
		Signature: signature,
	}
}

// Route should return the name of the module
func (msg MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetConfirm) Type() string { return "valset_confirm" }

// Stateless checks
func (msg MsgValsetConfirm) ValidateBasic() error {
	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}

// ValsetRequest
// -------------
type MsgValsetRequest struct {
}

func NewMsgValsetRequest() MsgValsetRequest {
	return MsgValsetRequest{}
}

// Route should return the name of the module
func (msg MsgValsetRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetRequest) Type() string { return "valset_request" }

func (msg MsgValsetRequest) ValidateBasic() error { return nil }

// GetSignBytes encodes the message for signing
func (msg MsgValsetRequest) GetSignBytes() []byte {
	return []byte{} // Does this work?
}

// GetSigners defines whose signature is required
func (msg MsgValsetRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{} // Does this work?
}

// SetEthAddress
// -------------
type MsgSetEthAddress struct {
	Address   string         `json:"address"`
	Validator sdk.AccAddress `json:"validator"`
	Signature []byte         `json:"signature"`
}

func NewMsgSetEthAddress(address string, validator sdk.AccAddress, signature []byte) MsgSetEthAddress {
	return MsgSetEthAddress{
		Address:   address,
		Validator: validator,
		Signature: signature,
	}
}

// Route should return the name of the module
func (msg MsgSetEthAddress) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetEthAddress) Type() string { return "set_eth_address" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid, and whether the Eth address has signed the validator address
// (proving control of the Eth address)
func (msg MsgSetEthAddress) ValidateBasic() error {
	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
	}
	valdiateEthAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if !valdiateEthAddr.MatchString(msg.Address) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "This is not a valid Ethereum address") // TODO: what error type to use here?
	}

	err := utils.ValidateEthSig(crypto.Keccak256(msg.Validator), msg.Signature, msg.Address)

	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetEthAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetEthAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string         `json:"name"`
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgSetName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetName) Type() string { return "set_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Value cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyName) Type() string { return "buy_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// MsgDeleteName defines a DeleteName message
type MsgDeleteName struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgDeleteName is a constructor function for MsgDeleteName
func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		Name:  name,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgDeleteName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeleteName) Type() string { return "delete_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
