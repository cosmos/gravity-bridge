package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid        = sdkerrors.Register(ModuleName, 3, "invalid")
	ErrSupplyOverflow = sdkerrors.Register(ModuleName, 4, "malicious ERC20 with invalid supply sent over bridge")
)
