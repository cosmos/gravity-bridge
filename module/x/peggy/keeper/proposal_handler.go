package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func HandlePeggyBootstrapProposal(ctx sdk.Context, k Keeper, p *types.PeggyBootstrapProposal) error {
	// TODO
	return nil
}

func HandlePeggyUpgradeProposal(ctx sdk.Context, k Keeper, p *types.PeggyUpgradeProposal) error {
	// TODO
	return nil
}
