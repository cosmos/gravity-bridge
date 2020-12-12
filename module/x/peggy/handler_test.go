package peggy

import (
	"math"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleValsetRequest(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr, _                   = sdk.AccAddressFromBech32("cosmos1990z7dqsvh8gthw9pa5sn4wuy2xrsd80mg5z6y")
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		myBlockHeight      int64          = 200
	)

	k, ctx, _ := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	h := NewHandler(k)
	msg := &types.MsgValsetRequest{Requester: myCosmosAddr.String()}
	ctx = ctx.WithBlockTime(myBlockTime).WithBlockHeight(myBlockHeight)
	res, err := h(ctx, msg)
	// then
	require.NoError(t, err)
	nonce := types.UInt64FromBytes(res.Data)
	require.False(t, nonce == 0)
	require.Equal(t, uint64(myBlockHeight), nonce)
	// and persisted
	valset := k.GetValsetRequest(ctx, nonce)
	require.NotNil(t, valset)
	assert.Equal(t, nonce, valset.Nonce)
	require.Len(t, valset.Members, 1)
	assert.Equal(t, []uint64{math.MaxUint32}, types.BridgeValidators(valset.Members).GetPowers())
	assert.Equal(t, "", valset.Members[0].EthereumAddress)
}

func TestHandleMsgSendToEth(t *testing.T) {
	var (
		userCosmosAddr, _            = sdk.AccAddressFromBech32("cosmos1990z7dqsvh8gthw9pa5sn4wuy2xrsd80mg5z6y")
		blockTime                    = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		blockHeight        int64     = 200
		denom                        = "peggy/0xB5E9944950C97acab395a324716D186632789712"
		startingCoinAmount sdk.Int   = sdk.NewIntFromUint64(150)
		sendAmount         sdk.Int   = sdk.NewIntFromUint64(50)
		feeAmount          sdk.Int   = sdk.NewIntFromUint64(5)
		startingCoins      sdk.Coins = sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}
		sendingCoin        sdk.Coin  = sdk.NewCoin(denom, sendAmount)
		feeCoin            sdk.Coin  = sdk.NewCoin(denom, feeAmount)
		ethDestination               = "0x3c9289da00b02dC623d0D8D907619890301D26d4"
	)

	// we start by depositing some funds into the users balance to send
	k, ctx, keepers := keeper.CreateTestEnv(t)
	h := NewHandler(k)
	keepers.BankKeeper.MintCoins(ctx, types.ModuleName, startingCoins)
	keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userCosmosAddr, startingCoins)
	balance1 := keepers.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}, balance1)

	// send some coins
	msg := &types.MsgSendToEth{
		Sender:    userCosmosAddr.String(),
		EthDest:   ethDestination,
		Amount:    sendingCoin,
		BridgeFee: feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err := h(ctx, msg)
	require.NoError(t, err)
	balance2 := keepers.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount.Sub(sendAmount).Sub(feeAmount))}, balance2)

	// do the same thing again and make sure it works twice
	msg1 := &types.MsgSendToEth{
		Sender:    userCosmosAddr.String(),
		EthDest:   ethDestination,
		Amount:    sendingCoin,
		BridgeFee: feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err1 := h(ctx, msg1)
	require.NoError(t, err1)
	balance3 := keepers.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	finalAmount3 := startingCoinAmount.Sub(sendAmount).Sub(sendAmount).Sub(feeAmount).Sub(feeAmount)
	assert.Equal(t, sdk.Coins{sdk.NewCoin(denom, finalAmount3)}, balance3)

	// now we should be out of coins and error
	msg2 := &types.MsgSendToEth{
		Sender:    userCosmosAddr.String(),
		EthDest:   ethDestination,
		Amount:    sendingCoin,
		BridgeFee: feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err2 := h(ctx, msg2)
	require.Error(t, err2)
	balance4 := keepers.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewCoin(denom, finalAmount3)}, balance4)
}

func TestHandleCreateEthereumClaimsSingleValidator(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr, _                   = sdk.AccAddressFromBech32("cosmos16ahjkfqxpp6lvfy9fpfnfjg39xr96qett0alj5")
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myNonce                           = uint64(1)
		anyETHAddr                        = "0xf9613b532673Cc223aBa451dFA8539B87e1F666D"
		tokenETHAddr                      = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	k, ctx, keepers := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	h := NewHandler(k)

	myErc20 := types.ERC20Token{
		Amount:   sdk.NewInt(12),
		Contract: tokenETHAddr,
	}

	ethClaim := types.MsgDepositClaim{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   myOrchestratorAddr.String(),
	}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err := h(ctx, &ethClaim)
	require.NoError(t, err)
	// and claim persisted
	claimFound := k.HasClaim(ctx, types.CLAIM_TYPE_DEPOSIT, myNonce, myValAddr, &ethClaim)
	assert.True(t, claimFound)
	// and attestation persisted
	a := k.GetAttestation(ctx, myNonce, &ethClaim)
	require.NotNil(t, a)
	// and vouchers added to the account
	balance := keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance)

	// Test to reject duplicate deposit
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim)
	// then
	require.Error(t, err)
	balance = keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance)

	// Test to reject skipped nonce
	ethClaim = types.MsgDepositClaim{
		EventNonce:     uint64(3),
		TokenContract:  tokenETHAddr,
		Amount:         sdk.NewInt(12),
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   myOrchestratorAddr.String(),
	}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim)
	// then
	require.Error(t, err)
	balance = keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance)

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
	// then
	require.NoError(t, err)
	balance = keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 25)}, balance)
}

func TestHandleCreateEthereumClaimsMultiValidator(t *testing.T) {
	var (
		orchestratorAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		orchestratorAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		orchestratorAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		myCosmosAddr, _      = sdk.AccAddressFromBech32("cosmos16ahjkfqxpp6lvfy9fpfnfjg39xr96qett0alj5")
		valAddr1             = sdk.ValAddress(orchestratorAddr1) // revisit when proper mapping is impl in keeper
		valAddr2             = sdk.ValAddress(orchestratorAddr2) // revisit when proper mapping is impl in keeper
		valAddr3             = sdk.ValAddress(orchestratorAddr3) // revisit when proper mapping is impl in keeper
		myNonce              = uint64(1)
		anyETHAddr           = "0xf9613b532673Cc223aBa451dFA8539B87e1F666D"
		tokenETHAddr         = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
		myBlockTime          = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	k, ctx, keepers := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
	h := NewHandler(k)

	myErc20 := types.ERC20Token{
		Amount:   sdk.NewInt(12),
		Contract: tokenETHAddr,
	}

	ethClaim1 := types.MsgDepositClaim{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   orchestratorAddr1.String(),
	}
	ethClaim2 := types.MsgDepositClaim{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   orchestratorAddr2.String(),
	}
	ethClaim3 := types.MsgDepositClaim{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
		Orchestrator:   orchestratorAddr3.String(),
	}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err := h(ctx, &ethClaim1)
	require.NoError(t, err)
	// and claim persisted
	claimFound1 := k.HasClaim(ctx, types.CLAIM_TYPE_DEPOSIT, myNonce, valAddr1, &ethClaim1)
	assert.True(t, claimFound1)
	// and attestation persisted
	a1 := k.GetAttestation(ctx, myNonce, &ethClaim1)
	require.NotNil(t, a1)
	// and vouchers not yet added to the account
	balance1 := keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.NotEqual(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance1)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim2)
	require.NoError(t, err)

	// and claim persisted
	claimFound2 := k.HasClaim(ctx, types.CLAIM_TYPE_DEPOSIT, myNonce, valAddr1, &ethClaim2)
	assert.True(t, claimFound2)
	// and attestation persisted
	a2 := k.GetAttestation(ctx, myNonce, &ethClaim1)
	require.NotNil(t, a2)
	// and vouchers now added to the account
	balance2 := keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance2)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, &ethClaim3)
	require.NoError(t, err)

	// and claim persisted
	claimFound3 := k.HasClaim(ctx, types.CLAIM_TYPE_DEPOSIT, myNonce, valAddr1, &ethClaim2)
	assert.True(t, claimFound3)
	// and attestation persisted
	a3 := k.GetAttestation(ctx, myNonce, &ethClaim1)
	require.NotNil(t, a3)
	// and no additional added to the account
	balance3 := keepers.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy/0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance3)
}
