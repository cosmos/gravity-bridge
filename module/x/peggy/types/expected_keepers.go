package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []staking.Validator
	GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64
	GetLastTotalPower(ctx sdk.Context) (power sdk.Int)
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}
