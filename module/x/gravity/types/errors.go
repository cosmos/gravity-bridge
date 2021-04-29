package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ethereum bridge errors
var (
	ErrInvalidGravityDenom         = sdkerrors.Register(ModuleName, 2, "invalid denomination for bridge transfer")
	ErrContractNotFound            = sdkerrors.Register(ModuleName, 3, "contract address not found for ERC20 token")
	ErrContractExists              = sdkerrors.Register(ModuleName, 4, "contract address already exists")
	ErrEventNotFound               = sdkerrors.Register(ModuleName, 5, "ethereum event not found")
	ErrEventUnsupported            = sdkerrors.Register(ModuleName, 6, "ethereum event type unsupported")
	ErrEventInvalid                = sdkerrors.Register(ModuleName, 7, "invalid ethereum event")
	ErrTxNotFound                  = sdkerrors.Register(ModuleName, 8, "outgoing transaction not found")
	ErrValidatorNotBonded          = sdkerrors.Register(ModuleName, 9, "validator is not bonded")
	ErrSignerSetNotFound           = sdkerrors.Register(ModuleName, 10, "ethereum signer set not found")
	ErrValidatorEthAddressNotFound = sdkerrors.Register(ModuleName, 11, "ethereum address not found for validator")
	ErrConfirmUnsupported          = sdkerrors.Register(ModuleName, 12, "confirm tx type unsupported")
	ErrSignatureInvalid            = sdkerrors.Register(ModuleName, 13, "invalid signature")
	ErrSignatureDuplicate          = sdkerrors.Register(ModuleName, 14, "confirm already signed")
)
