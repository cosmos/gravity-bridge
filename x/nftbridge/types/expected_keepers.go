package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft/exported"
)

// NFTKeeper defines the expected nft keeper
type NFTKeeper interface {
	// SendNFTFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, denom, id string) error
	// SendNFTFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, denom, id string) error
	// MintNFT(ctx sdk.Context, name string, denom, id string) error
	// BurnNFT(ctx sdk.Context, name string, denom, id string) error
	// SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	GetNFT(ctx sdk.Context, denom, id string) (nft exported.NFT, err error)
	UpdateNFT(ctx sdk.Context, denom string, nft exported.NFT) (err error)
	MintNFT(ctx sdk.Context, denom string, nft exported.NFT) (err error)
	DeleteNFT(ctx sdk.Context, denom, id string) (err error)
	IsNFT(ctx sdk.Context, denom, id string) (exists bool)
}
