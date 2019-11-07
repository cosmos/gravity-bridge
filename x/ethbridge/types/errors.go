package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType local code type
type CodeType = sdk.CodeType

// Exported code type numbers
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidEthNonce     CodeType = 1
	CodeInvalidEthAddress   CodeType = 2
	CodeErrJSONMarshalling  CodeType = 3
	CodeInvalidEthSymbol    CodeType = 4
	CodeErrInvalidClaimType CodeType = 5
	CodeErrInvalidChainID   CodeType = 6
)

// ErrInvalidEthNonce implements sdk.Error.
func ErrInvalidEthNonce(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEthNonce, "invalid ethereum nonce provided, must be >= 0")
}

// ErrInvalidEthAddress implements sdk.Error.
func ErrInvalidEthAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEthAddress, "invalid ethereum address provided, must be a valid hex-encoded Ethereum address")
}

// ErrInvalidChainID implements sdk.Error.
func ErrInvalidChainID(codespace sdk.CodespaceType, chainID string) sdk.Error {
	return sdk.NewError(codespace, CodeErrInvalidChainID, fmt.Sprintf("invalid ethereum chain id '%s'", chainID))
}

// ErrJSONMarshalling implements sdk.Error.
func ErrJSONMarshalling(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrJSONMarshalling, "error marshalling JSON for this claim")
}

// ErrInvalidEthSymbol implements sdk.Error.
func ErrInvalidEthSymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEthSymbol, "invalid symbol provided, symbol \"eth\" must have null address set as token contract address")
}

// ErrInvalidClaimType implements sdk.Error.
func ErrInvalidClaimType() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeErrInvalidClaimType, "invalid claim type provided")
}
