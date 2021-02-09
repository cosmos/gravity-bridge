package peggy

import (
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO:
// Have the validators put in a erc20<>denom relation with ERC20DeployedEvent
// Send some coins of that denom into the cosmos module
// Check that the coins are locked, not burned
// Have the validators put in a deposit event for that ERC20
// Check that the coins are unlocked and sent to the right account

func addDenomToERC20Relation(tv *testingVars) {
	denom := "uatom"
	erc20 := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"

	ethClaim := types.MsgERC20DeployedClaim{
		CosmosDenom:   denom,
		TokenContract: erc20,
		Name:          "Atom",
		Symbol:        "ATOM",
		Decimals:      6,
		EventNonce:    tv.myNonce,
		Orchestrator:  tv.myOrchestratorAddr.String(),
	}

	// when
	tv.ctx = tv.ctx.WithBlockTime(tv.myBlockTime)
	_, err := tv.h(tv.ctx, &ethClaim)
	EndBlocker(tv.ctx, tv.input.PeggyKeeper)
	require.NoError(tv.t, err)

	// and attestation persisted
	a := tv.input.PeggyKeeper.GetAttestation(tv.ctx, tv.myNonce, ethClaim.ClaimHash())
	require.NotNil(tv.t, a)

	// and erc20<>denom relation added to db
	isCosmosOriginated, gotERC20, err := tv.input.PeggyKeeper.DenomToERC20(tv.ctx, denom)
	require.NoError(tv.t, err)
	assert.True(tv.t, isCosmosOriginated)

	isCosmosOriginated, gotDenom := tv.input.PeggyKeeper.ERC20ToDenom(tv.ctx, erc20)
	assert.True(tv.t, isCosmosOriginated)

	assert.Equal(tv.t, denom, gotDenom)
	assert.Equal(tv.t, erc20, gotERC20)
}

type testingVars struct {
	myOrchestratorAddr sdk.AccAddress
	myCosmosAddr       sdk.AccAddress
	myValAddr          sdk.ValAddress
	myNonce            uint64
	anyETHAddr         string
	tokenETHAddr       string
	myBlockTime        time.Time
	input              keeper.TestInput
	ctx                sdk.Context
	h                  sdk.Handler
	t                  *testing.T
}

func initializeTestingVars(t *testing.T) testingVars {
	var tv testingVars

	tv.t = t

	tv.myOrchestratorAddr = make([]byte, sdk.AddrLen)
	tv.myCosmosAddr, _ = sdk.AccAddressFromBech32("cosmos16ahjkfqxpp6lvfy9fpfnfjg39xr96qett0alj5")
	tv.myValAddr = sdk.ValAddress(tv.myOrchestratorAddr) // revisit when proper mapping is impl in keeper
	tv.myNonce = uint64(1)
	tv.anyETHAddr = "0xf9613b532673Cc223aBa451dFA8539B87e1F666D"
	tv.tokenETHAddr = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
	tv.myBlockTime = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)

	tv.input = keeper.CreateTestEnv(t)
	tv.ctx = tv.input.Context
	tv.input.PeggyKeeper.StakingKeeper = keeper.NewStakingKeeperMock(tv.myValAddr)
	tv.input.PeggyKeeper.SetOrchestratorValidator(tv.ctx, tv.myValAddr, tv.myOrchestratorAddr)
	tv.h = NewHandler(tv.input.PeggyKeeper)

	return tv
}

func TestCosmosOriginated(t *testing.T) {
	tv := initializeTestingVars(t)
	addDenomToERC20Relation(&tv)
	//

	// // Test to reject denom with wrong info
	// // when
	// ctx = ctx.WithBlockTime(myBlockTime)
	// _, err = h(ctx, &ethClaim)
	// EndBlocker(ctx, input.PeggyKeeper)
	// // then
	// require.Error(t, err)
	// balance = input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	// assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance)

	// // Test to reject skipped nonce
	// ethClaim = types.MsgDepositClaim{
	// 	EventNonce:     uint64(3),
	// 	TokenContract:  tokenETHAddr,
	// 	Amount:         sdk.NewInt(12),
	// 	EthereumSender: anyETHAddr,
	// 	CosmosReceiver: myCosmosAddr.String(),
	// 	Orchestrator:   myOrchestratorAddr.String(),
	// }

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim)
	EndBlocker(ctx, input.PeggyKeeper)
	// then
	require.Error(t, err)
	balance = input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance)

	// Test to finally accept consecutive nonce
	ethClaim = types.MsgDepositClaim{
		EventNonce:     uint64(2),
		Amount:         sdk.NewInt(13),
		TokenContract:  tokenETHAddr,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   myOrchestratorAddr.String(),
	}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim)
	EndBlocker(ctx, input.PeggyKeeper)

	// then
	require.NoError(t, err)
	balance = input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 25)}, balance)
}
