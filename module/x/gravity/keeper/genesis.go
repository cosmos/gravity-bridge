package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO:
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {

}

// TODO:
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		Params:        k.GetParams(ctx),
		Erc20ToDenoms: k.GetERC20Denoms(ctx),
	}
}
