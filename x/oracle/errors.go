package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

//Exported code type numbers
const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeInvalidNonce  CodeType = 1
	CodeNotFound      CodeType = 2
	CodeMinimumTooLow CodeType = 3
)

// ErrInvalidNonce if prophecy or claim is missing nonce
func ErrInvalidNonce(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidNonce, "invalid nonce provided, must be an integer >= 0")
}

func ErrNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotFound, "prophecy or claim with given nonce not found")
}

func ErrMinimumTooLow(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinimumTooLow, "minimum number of validators must be greater than 1")
}
