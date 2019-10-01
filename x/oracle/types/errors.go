package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

// Exported code type numbers
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeProphecyNotFound              CodeType = 1
	CodeMinimumConsensusNeededInvalid CodeType = 2
	CodeNoClaims                      CodeType = 3
	CodeInvalidIdentifier             CodeType = 4
	CodeProphecyFinalized             CodeType = 5
	CodeDuplicateMessage              CodeType = 6
	CodeInvalidClaim                  CodeType = 7
	CodeInvalidValidator              CodeType = 8
	CodeInternalDB                    CodeType = 9
)

func ErrProphecyNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProphecyNotFound, "prophecy with given id not found")
}

func ErrMinimumConsensusNeededInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinimumConsensusNeededInvalid, "minimum consensus proportion of validator staking power must be > 0 and <= 1")
}

func ErrNoClaims(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoClaims, "cannot create prophecy without initial claim")
}

func ErrInvalidIdentifier(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIdentifier, "invalid identifier provided, must be a nonempty string")
}

func ErrProphecyFinalized(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProphecyFinalized, "prophecy already finalized")
}

func ErrDuplicateMessage(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateMessage, "already processed message from validator for this id")
}

func ErrInvalidClaim(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidClaim, "claim cannot be empty string")
}

func ErrInvalidValidator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "claim must be made by actively bonded validator")
}

func ErrInternalDB(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInternalDB, fmt.Sprintf("internal error serializing/deserializing prophecy: %s", err.Error()))
}
