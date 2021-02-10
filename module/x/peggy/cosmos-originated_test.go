package peggy

import (
	"testing"

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

func TestCosmosOriginated(t *testing.T) {
	tv := initializeTestingVars(t)
	addDenomToERC20Relation(tv)
	lockCoinsInModule(tv)
}

type testingVars struct {
	myOrchestratorAddr sdk.AccAddress
	// myCosmosAddr       sdk.AccAddress
	myValAddr sdk.ValAddress
	// myNonce            uint64
	// anyETHAddr         string
	// erc20              string
	// denom              string
	// myBlockTime        time.Time
	input keeper.TestInput
	ctx   sdk.Context
	h     sdk.Handler
	t     *testing.T
}

func initializeTestingVars(t *testing.T) *testingVars {
	var tv testingVars

	tv.t = t

	tv.myOrchestratorAddr = make([]byte, sdk.AddrLen)
	tv.myValAddr = sdk.ValAddress(tv.myOrchestratorAddr) // revisit when proper mapping is impl in keeper

	tv.input = keeper.CreateTestEnv(t)
	tv.ctx = tv.input.Context
	tv.input.PeggyKeeper.StakingKeeper = keeper.NewStakingKeeperMock(tv.myValAddr)
	tv.input.PeggyKeeper.SetOrchestratorValidator(tv.ctx, tv.myValAddr, tv.myOrchestratorAddr)
	tv.h = NewHandler(tv.input.PeggyKeeper)

	return &tv
}

func addDenomToERC20Relation(tv *testingVars) {
	var (
		myNonce = uint64(1)
		denom   = "uatom"
		erc20   = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
	)

	ethClaim := types.MsgERC20DeployedClaim{
		CosmosDenom:   denom,
		TokenContract: erc20,
		Name:          "Atom",
		Symbol:        "ATOM",
		Decimals:      6,
		EventNonce:    myNonce,
		Orchestrator:  tv.myOrchestratorAddr.String(),
	}

	_, err := tv.h(tv.ctx, &ethClaim)
	require.NoError(tv.t, err)

	EndBlocker(tv.ctx, tv.input.PeggyKeeper)

	// check if attestation persisted
	a := tv.input.PeggyKeeper.GetAttestation(tv.ctx, myNonce, ethClaim.ClaimHash())
	require.NotNil(tv.t, a)

	// check if erc20<>denom relation added to db
	isCosmosOriginated, gotERC20, err := tv.input.PeggyKeeper.DenomToERC20(tv.ctx, denom)
	require.NoError(tv.t, err)
	assert.True(tv.t, isCosmosOriginated)

	isCosmosOriginated, gotDenom := tv.input.PeggyKeeper.ERC20ToDenom(tv.ctx, erc20)
	assert.True(tv.t, isCosmosOriginated)

	assert.Equal(tv.t, denom, gotDenom)
	assert.Equal(tv.t, erc20, gotERC20)
}

func lockCoinsInModule(tv *testingVars) {
	var (
		userCosmosAddr, _            = sdk.AccAddressFromBech32("cosmos1990z7dqsvh8gthw9pa5sn4wuy2xrsd80mg5z6y")
		denom                        = "uatom"
		startingCoinAmount sdk.Int   = sdk.NewIntFromUint64(150)
		sendAmount         sdk.Int   = sdk.NewIntFromUint64(50)
		feeAmount          sdk.Int   = sdk.NewIntFromUint64(5)
		startingCoins      sdk.Coins = sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}
		sendingCoin        sdk.Coin  = sdk.NewCoin(denom, sendAmount)
		feeCoin            sdk.Coin  = sdk.NewCoin(denom, feeAmount)
		ethDestination               = "0x3c9289da00b02dC623d0D8D907619890301D26d4"
	)

	// we start by depositing some funds into the users balance to send
	input := keeper.CreateTestEnv(tv.t)
	ctx := input.Context
	h := NewHandler(input.PeggyKeeper)
	input.BankKeeper.MintCoins(ctx, types.ModuleName, startingCoins)
	input.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userCosmosAddr, startingCoins)
	balance1 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	assert.Equal(tv.t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}, balance1)

	// send some coins
	msg := &types.MsgSendToEth{
		Sender:    userCosmosAddr.String(),
		EthDest:   ethDestination,
		Amount:    sendingCoin,
		BridgeFee: feeCoin,
	}

	_, err := h(ctx, msg)
	require.NoError(tv.t, err)

	// Check that user balance has gone down
	balance2 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	assert.Equal(tv.t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount.Sub(sendAmount).Sub(feeAmount))}, balance2)

	// Check that peggy balance has gone up
	peggyAddr := input.AccountKeeper.GetModuleAddress(types.ModuleName)
	assert.Equal(tv.t,
		sdk.Coins{sdk.NewCoin(denom, sendAmount)},
		input.BankKeeper.GetAllBalances(ctx, peggyAddr),
	)
}
