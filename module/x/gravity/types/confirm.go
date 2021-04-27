package types

import (
	"encoding/hex"
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

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
	GetInvalidationID() tmbytes.HexBytes
	GetInvalidationNonce() uint64
}

var (
	_ Confirm = &ConfirmBatch{}
	_ Confirm = &ConfirmLogicCall{}
	_ Confirm = &ConfirmSignerSet{}
)

// GetType should return the action
func (c ConfirmBatch) GetType() string { return "batch" }

// Validate performs stateless checks
func (c ConfirmBatch) Validate() error {
	if _, err := sdk.AccAddressFromBech32(c.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, c.OrchestratorAddress)
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddress(c.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	_, err := hex.DecodeString(c.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "could not decode hex string %s", c.Signature)
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (c ConfirmBatch) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationID is a noop to implement confirm interface
func (c ConfirmBatch) GetInvalidationID() tmbytes.HexBytes { return nil }

// GetType should return the action
func (c ConfirmLogicCall) GetType() string { return "logic_Call" }

// Validate performs stateless checks
func (c ConfirmLogicCall) Validate() error {
	if _, err := sdk.AccAddressFromBech32(c.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, c.OrchestratorAddress)
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	_, err := hex.DecodeString(c.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", c.Signature)
	}
	if len(c.InvalidationID) == 0 {
		return fmt.Errorf("invalidation id is empty")
	}
	return nil
}

func (c ConfirmLogicCall) GetNonce() uint64 {
	return 0
}

func (c ConfirmLogicCall) GetTokenContract() string {
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

// GetInvalidationID is a noop to implement confirm interface
func (c ConfirmSignerSet) GetInvalidationID() tmbytes.HexBytes { return nil }

func (c ConfirmSignerSet) GetTokenContract() string {
	return ""
}
