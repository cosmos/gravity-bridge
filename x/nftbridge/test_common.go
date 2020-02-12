package nftbridge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/modules/incubator/nft"
	oracle "github.com/cosmos/peggy/x/oracle"
	keeperLib "github.com/cosmos/peggy/x/oracle/keeper"
)

// CreateTestHandler creates a test handler
func CreateTestHandler(t *testing.T, consensusNeeded float64, validatorAmounts []int64) (sdk.Context, oracle.Keeper, bank.Keeper, nft.Keeper, auth.AccountKeeper, []sdk.ValAddress, sdk.Handler) {
	ctx, oracleKeeper, bankKeeper, _, nftKeeper, accountKeeper, validatorAddresses := oracle.CreateTestKeepers(t, consensusNeeded, validatorAmounts, ModuleName)
	// bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	// nftKeeper.SetModuleAccount(ctx, bridgeAccount)

	cdc := keeperLib.MakeTestCodec()
	bridgeKeeper := NewKeeper(cdc, nftKeeper, oracleKeeper)
	handler := NewHandler(nftKeeper, bridgeKeeper, cdc)

	return ctx, oracleKeeper, bankKeeper, nftKeeper, accountKeeper, validatorAddresses, handler
}
