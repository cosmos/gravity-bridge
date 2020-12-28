package peggy

import (
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
)

var (
	stakingAmt = sdk.TokensFromConsensusPower(10)
)

func setupVal5(t *testing.T) (keeper.TestKeepers, sdk.Handler) {
	peggyKeeper, ctx, testKeeper := keeper.CreateTestEnv(t)
	params := peggyKeeper.GetParams(ctx)
	params.SignedBlocksWindow = 10
	params.SlashFractionValset = sdk.NewDecWithPrec(1, 2)
	params.PeggyId = "lkasjdfklajsldkfjd"
	params.ContractSourceHash = "lkasjdfklajsldkfjd"
	params.StartThreshold = 0
	params.EthereumAddress = "0x8858eeb3dfffa017d4bce9801d340d36cf895ccf"
	params.BridgeChainId = 11
	peggyKeeper.SetParams(ctx, params)
	h := NewHandler(peggyKeeper)

	sh := staking.NewHandler(testKeeper.StakingKeeper)

	bd := testKeeper.StakingKeeper.GetParams(ctx).BondDenom
	for i := range []int{0, 1, 2, 3, 4} {
		acc := testKeeper.AccountKeeper.NewAccount(ctx, authtypes.NewBaseAccount(keeper.Addrs[i], keeper.AccPubKeys[i], uint64(i), 0))
		testKeeper.BankKeeper.SetBalances(ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, stakingAmt.Add(sdk.NewInt(100)))))
		testKeeper.AccountKeeper.SetAccount(ctx, acc)
	}

	// Validator created
	_, err := sh(ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.PubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh(ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.PubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh(ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.PubKeys[2], stakingAmt))
	require.NoError(t, err)
	_, err = sh(ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[3], keeper.PubKeys[3], stakingAmt))
	require.NoError(t, err)
	_, err = sh(ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[4], keeper.PubKeys[4], stakingAmt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, testKeeper.StakingKeeper)

	return testKeeper, h
}
