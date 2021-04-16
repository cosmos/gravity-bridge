package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Confirm interface {
	Type() ConfirmType
	GetOrchestratorAddress() string
	GetNonce() uint64
	GetSignature() string

	GetTokenContract() string
	GetInvalidationID() string
	GetInvalidationNonce() uint64
}

var (
	_ Confirm = &ConfirmBatch{}
	_ Confirm = &ConfirmLogicCall{}
	_ Confirm = &ValsetConfirm{}
)

// Type should return the action
func (msg ConfirmBatch) Type() ConfirmType { return ConfirmType_CONFIRM_TYPE_BATCH }

// ValidateBasic performs stateless checks
func (msg ConfirmBatch) ValidateBasic() error {
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
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (msg ConfirmBatch) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationId is a noop to implement confirm interface
func (msg ConfirmBatch) GetInvalidationID() string { return "" }

// Type should return the action
func (msg ConfirmLogicCall) Type() ConfirmType { return ConfirmType_CONFIRM_TYPE_LOGIC }

// ValidateBasic performs stateless checks
func (msg ConfirmLogicCall) ValidateBasic() error {
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

// GetInvalidationId is a noop to implement confirm interface
func (msg ConfirmLogicCall) GetInvalidationID() string { return "" }

// NewValsetConfirm returns a new ValsetConfirm
func NewValsetConfirm(nonce uint64, ethAddress string, validator sdk.AccAddress, signature string) *ValsetConfirm {
	return &ValsetConfirm{
		Nonce:               nonce,
		OrchestratorAddress: validator.String(),
		EthAddress:          ethAddress,
		Signature:           signature,
	}
}

// Type should return the action
func (msg *ValsetConfirm) Type() ConfirmType { return ConfirmType_CONFIRM_TYPE_VALSET }

// ValidateBasic performs stateless checks
func (msg *ValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (msg ValsetConfirm) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationId is a noop to implement confirm interface
func (msg ValsetConfirm) GetInvalidationID() string { return "" }

func (msg *ValsetConfirm) GetTokenContract() string {
	return ""
}
