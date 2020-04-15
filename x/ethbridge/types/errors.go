package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidEthNonce   = sdkerrors.Register(ModuleName, 1, "invalid ethereum nonce provided, must be >= 0")
	ErrInvalidEthAddress = sdkerrors.Register(ModuleName, 2,
		"invalid ethereum address provided, must be a valid hex-encoded Ethereum address")
	ErrJSONMarshalling  = sdkerrors.Register(ModuleName, 3, "error marshalling JSON for this claim")
	ErrInvalidEthSymbol = sdkerrors.Register(ModuleName, 4,
		"invalid symbol provided, symbol 'eth' must have null address set as token contract address")
	ErrInvalidClaimType       = sdkerrors.Register(ModuleName, 5, "invalid claim type provided")
	ErrInvalidEthereumChainID = sdkerrors.Register(ModuleName, 6, "invalid ethereum chain id")
	ErrInvalidAmount          = sdkerrors.Register(ModuleName, 7, "amount must be a valid integer > 0")
	ErrInvalidSymbol          = sdkerrors.Register(ModuleName, 8, "symbol must be 1 character or more")
	ErrInvalidBurnSymbol      = sdkerrors.Register(ModuleName, 9,
		fmt.Sprintf("symbol of token to burn must be in the form %v{ethereumSymbol}", PeggedCoinPrefix))
)
