package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/cosmos/peggy/x/oracle"
)

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendNFTFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, denom, id string) error
	SendNFTFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, denom, id string) error
	MintNFT(ctx sdk.Context, name string, denom, id string) error
	BurnNFT(ctx sdk.Context, name string, denom, id string) error
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(ctx sdk.Context, claim oracle.Claim) (oracle.Status, error)
	GetProphecy(ctx sdk.Context, id string) (oracle.Prophecy, bool)
}
