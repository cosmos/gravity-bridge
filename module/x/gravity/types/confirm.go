package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

type Confirm interface {
	proto.Message

	GetType() string
	// TODO: delete
	GetOrchestratorAddress() string
	GetNonce() uint64
	GetSignature() string
	Validate() error

	// TODO: consider deleting
	GetTokenContract() string
	GetInvalidationId() string
	GetInvalidationNonce() uint64
}

var (
	_ Confirm = &ConfirmBatch{}
	_ Confirm = &ConfirmLogicCall{}
	_ Confirm = &ConfirmSignerSet{}
)

// GetType should return the action
func (msg ConfirmBatch) GetType() string { return "batch" }

// Validate performs stateless checks
func (msg ConfirmBatch) Validate() error {
	if _, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "could not decode hex string %s", msg.Signature)
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (msg ConfirmBatch) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationId is a noop to implement confirm interface
func (msg ConfirmBatch) GetInvalidationId() string { return "" }

// GetType should return the action
func (msg ConfirmLogicCall) GetType() string { return "logic_Call" }

// Validate performs stateless checks
func (msg ConfirmLogicCall) Validate() error {
	if _, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	_, err = hex.DecodeString(msg.InvalidationId)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.InvalidationId)
	}
	return nil
}

func (msg ConfirmLogicCall) GetNonce() uint64 {
	return 0
}

func (msg ConfirmLogicCall) GetTokenContract() string {
	return ""
}

// NewConfirmSignerSet returns a new ConfirmSignerSet
func NewConfirmSignerSet(nonce uint64, ethSigner string, validator sdk.AccAddress, signature string) *ConfirmSignerSet {
	return &ConfirmSignerSet{
		Nonce:               nonce,
		OrchestratorAddress: validator.String(),
		EthSigner:           ethSigner,
		Signature:           signature,
	}
}

// GetType should return the action
func (c *ConfirmSignerSet) GetType() string { return "valset" }

// Validate performs stateless checks
func (c *ConfirmSignerSet) Validate() (err error) {
	if _, err = sdk.AccAddressFromBech32(c.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, c.OrchestratorAddress)
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}

	// TODO: validate signatre
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (c ConfirmSignerSet) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationId is a noop to implement confirm interface
func (c ConfirmSignerSet) GetInvalidationId() string { return "" }

func (c *ConfirmSignerSet) GetTokenContract() string {
	return ""
}
