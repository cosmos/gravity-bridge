package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDuplicate    = sdkerrors.Register(ModuleName, 2, "duplicate")
	ErrInvalidState = sdkerrors.Register(ModuleName, 3, "invalid state")
	ErrTimeout      = sdkerrors.Register(ModuleName, 4, "timeout")
	ErrUnknown      = sdkerrors.Register(ModuleName, 5, "unkown")
)
