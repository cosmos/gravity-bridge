package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

//Exported code type numbers
const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeProphecyNotFound   CodeType = 1
	CodeMinimumPowerTooLow CodeType = 2
	CodeNoClaims           CodeType = 3
	CodeInvalidIdentifier  CodeType = 4
)

func ErrProphecyNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProphecyNotFound, "prophecy with given id not found")
}

func ErrMinimumPowerTooLow(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinimumPowerTooLow, "minimum number for validator staking power must be greater than 1")
}

func ErrNoClaims(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoClaims, "cannot create prophecy without initial claim")
}

func ErrInvalidIdentifier(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIdentifier, "invalid identifier provided, must be a nonempty string")
}
