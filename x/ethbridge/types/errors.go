package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType local code type
type CodeType = sdk.CodeType

// Exported code type numbers
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidEthNonce    CodeType = 1
	CodeInvalidEthAddress  CodeType = 2
	CodeErrJSONMarshalling CodeType = 3
)

// ErrInvalidEthNonce implements sdk.Error.
func ErrInvalidEthNonce(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEthNonce, "invalid ethereum nonce provided, must be >= 0")
}

// ErrInvalidEthAddress implements sdk.Error.
func ErrInvalidEthAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEthAddress, "invalid ethereum address provided, must be a valid hex-encoded Ethereum address")
}

// ErrJSONMarshalling implements sdk.Error.
func ErrJSONMarshalling(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrJSONMarshalling, "error marshalling JSON for this claim")
}
