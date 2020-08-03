package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrMyCustomError = sdkerrors.Register(ModuleName, 1, "leaving this here as a reference for when we do our errors better")
)
