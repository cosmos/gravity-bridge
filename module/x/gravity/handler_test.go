package gravity

import (
	"bytes"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

func TestHandleMsgSendToEthereum(t *testing.T) {
	var (
		userCosmosAddr, _               = sdk.AccAddressFromBech32("cosmos1990z7dqsvh8gthw9pa5sn4wuy2xrsd80mg5z6y")
		blockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		blockHeight           int64     = 200
		denom                           = "gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
		startingCoinAmount, _           = sdk.NewIntFromString("150000000000000000000") // 150 ETH worth, required to reach above u64 limit (which is about 18 ETH)
		sendAmount, _                   = sdk.NewIntFromString("50000000000000000000")  // 50 ETH
		feeAmount, _                    = sdk.NewIntFromString("5000000000000000000")   // 5 ETH
		startingCoins         sdk.Coins = sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}
		sendingCoin           sdk.Coin  = sdk.NewCoin(denom, sendAmount)
		feeCoin               sdk.Coin  = sdk.NewCoin(denom, feeAmount)
		ethDestination                  = "0x3c9289da00b02dC623d0D8D907619890301D26d4"
	)

	// we start by depositing some funds into the users balance to send
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	h := NewHandler(input.GravityKeeper)
	input.BankKeeper.MintCoins(ctx, types.ModuleName, startingCoins)
	input.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userCosmosAddr, startingCoins) // 150
	balance1 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount)}, balance1) // 150

	// send some coins
	msg := &types.MsgSendToEthereum{
		Sender:            userCosmosAddr.String(),
		EthereumRecipient: ethDestination,
		Amount:            sendingCoin,
		BridgeFee:         feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err := h(ctx, msg) // send 55
	require.NoError(t, err)
	balance2 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin(denom, startingCoinAmount.Sub(sendAmount).Sub(feeAmount))}, balance2)

	// do the same thing again and make sure it works twice
	msg1 := &types.MsgSendToEthereum{
		Sender:            userCosmosAddr.String(),
		EthereumRecipient: ethDestination,
		Amount:            sendingCoin,
		BridgeFee:         feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err1 := h(ctx, msg1) // send 55
	require.NoError(t, err1)
	balance3 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	finalAmount3 := startingCoinAmount.Sub(sendAmount).Sub(sendAmount).Sub(feeAmount).Sub(feeAmount)
	require.Equal(t, sdk.Coins{sdk.NewCoin(denom, finalAmount3)}, balance3)

	// now we should be out of coins and error
	msg2 := &types.MsgSendToEthereum{
		Sender:            userCosmosAddr.String(),
		EthereumRecipient: ethDestination,
		Amount:            sendingCoin,
		BridgeFee:         feeCoin}
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err2 := h(ctx, msg2) // send 55
	require.Error(t, err2)
	balance4 := input.BankKeeper.GetAllBalances(ctx, userCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin(denom, finalAmount3)}, balance4)
}

func TestMsgSubmitEthreumEventSendToCosmosSingleValidator(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr, _                   = sdk.AccAddressFromBech32("cosmos16ahjkfqxpp6lvfy9fpfnfjg39xr96qett0alj5")
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myNonce                           = uint64(1)
		anyETHAddr                        = "0xf9613b532673Cc223aBa451dFA8539B87e1F666D"
		tokenETHAddr                      = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		amountA, _                        = sdk.NewIntFromString("50000000000000000000")  // 50 ETH
		amountB, _                        = sdk.NewIntFromString("100000000000000000000") // 100 ETH
	)
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	gk := input.GravityKeeper
	bk := input.BankKeeper
	gk.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	gk.SetOrchestratorValidatorAddress(ctx, myValAddr, myOrchestratorAddr)
	h := NewHandler(gk)

	myErc20 := types.ERC20Token{
		Amount:   amountA,
		Contract: tokenETHAddr,
	}

	sendToCosmosEvent := &types.SendToCosmosEvent{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}

	eva, err := types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgSubmitEvent := &types.MsgSubmitEthereumEvent{eva, myOrchestratorAddr.String()}
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, msgSubmitEvent)
	EndBlocker(ctx, gk)
	require.NoError(t, err)

	// and attestation persisted
	a := gk.GetEthereumEventVoteRecord(ctx, myNonce, sendToCosmosEvent.Hash())
	require.NotNil(t, a)
	// and vouchers added to the account

	balance := bk.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", amountA)}, balance)

	// Test to reject duplicate deposit
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, msgSubmitEvent)
	EndBlocker(ctx, gk)
	// then
	require.Error(t, err)
	balance = bk.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", amountA)}, balance)

	// Test to reject skipped nonce

	sendToCosmosEvent = &types.SendToCosmosEvent{
		EventNonce:     uint64(3),
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}

	eva, err = types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgSubmitEvent = &types.MsgSubmitEthereumEvent{eva, myOrchestratorAddr.String()}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, msgSubmitEvent)
	require.Error(t, err)

	EndBlocker(ctx, gk)
	// then
	balance = bk.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", amountA)}, balance)

	// Test to finally accept consecutive nonce
	sendToCosmosEvent = &types.SendToCosmosEvent{
		EventNonce:     uint64(2),
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}

	eva, err = types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgSubmitEvent = &types.MsgSubmitEthereumEvent{eva, myOrchestratorAddr.String()}
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, msgSubmitEvent)
	EndBlocker(ctx, gk)

	// then
	require.NoError(t, err)
	balance = bk.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewCoin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", amountB)}, balance)
}

func TestMsgSubmitEthreumEventSendToCosmosMultiValidator(t *testing.T) {
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
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	input.GravityKeeper.StakingKeeper = keeper.NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
	input.GravityKeeper.SetOrchestratorValidatorAddress(ctx, valAddr1, orchestratorAddr1)
	input.GravityKeeper.SetOrchestratorValidatorAddress(ctx, valAddr2, orchestratorAddr2)
	input.GravityKeeper.SetOrchestratorValidatorAddress(ctx, valAddr3, orchestratorAddr3)
	h := NewHandler(input.GravityKeeper)

	myErc20 := types.ERC20Token{
		Amount:   sdk.NewInt(12),
		Contract: tokenETHAddr,
	}

	ethClaim1 := &types.SendToCosmosEvent{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}
	ethClaim1a, err := types.PackEvent(ethClaim1)
	require.NoError(t, err)
	ethClaim1Msg := &types.MsgSubmitEthereumEvent{ethClaim1a, orchestratorAddr1.String()}
	ethClaim2 := &types.SendToCosmosEvent{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}
	ethClaim2a, err := types.PackEvent(ethClaim2)
	require.NoError(t, err)
	ethClaim2Msg := &types.MsgSubmitEthereumEvent{ethClaim2a, orchestratorAddr2.String()}
	ethClaim3 := &types.SendToCosmosEvent{
		EventNonce:     myNonce,
		TokenContract:  myErc20.Contract,
		Amount:         myErc20.Amount,
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr.String(),
	}
	ethClaim3a, err := types.PackEvent(ethClaim3)
	require.NoError(t, err)
	ethClaim3Msg := &types.MsgSubmitEthereumEvent{ethClaim3a, orchestratorAddr3.String()}

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, ethClaim1Msg)
	EndBlocker(ctx, input.GravityKeeper)
	require.NoError(t, err)
	// and attestation persisted
	a1 := input.GravityKeeper.GetEthereumEventVoteRecord(ctx, myNonce, ethClaim1.Hash())
	require.NotNil(t, a1)
	// and vouchers not yet added to the account
	balance1 := input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	require.NotEqual(t, sdk.Coins{sdk.NewInt64Coin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance1)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, ethClaim2Msg)
	EndBlocker(ctx, input.GravityKeeper)
	require.NoError(t, err)

	// and attestation persisted
	a2 := input.GravityKeeper.GetEthereumEventVoteRecord(ctx, myNonce, ethClaim2.Hash())
	require.NotNil(t, a2)
	// and vouchers now added to the account
	balance2 := input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance2)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err = h(ctx, ethClaim3Msg)
	EndBlocker(ctx, input.GravityKeeper)
	require.NoError(t, err)

	// and attestation persisted
	a3 := input.GravityKeeper.GetEthereumEventVoteRecord(ctx, myNonce, ethClaim3.Hash())
	require.NotNil(t, a3)
	// and no additional added to the account
	balance3 := input.BankKeeper.GetAllBalances(ctx, myCosmosAddr)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("gravity0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e", 12)}, balance3)
}

func TestMsgSetDelegateAddresses(t *testing.T) {
	var (
		ethAddress                    = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		cosmosAddress  sdk.AccAddress = bytes.Repeat([]byte{0x1}, sdk.AddrLen)
		ethAddress2                   = "0x26126048c706fB45a5a6De8432F428e794d0b952"
		cosmosAddress2 sdk.AccAddress = bytes.Repeat([]byte{0x2}, sdk.AddrLen)
		valAddress     sdk.ValAddress = bytes.Repeat([]byte{0x3}, sdk.AddrLen)
		blockTime                     = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		blockTime2                    = time.Date(2020, 9, 15, 15, 20, 10, 0, time.UTC)
		blockHeight    int64          = 200
		blockHeight2   int64          = 210
	)
	input := keeper.CreateTestEnv(t)
	input.GravityKeeper.StakingKeeper = keeper.NewStakingKeeperMock(valAddress)
	ctx := input.Context
	wctx := sdk.WrapSDKContext(ctx)
	k := input.GravityKeeper
	h := NewHandler(input.GravityKeeper)
	ctx = ctx.WithBlockTime(blockTime)

	msg := types.NewMsgDelegateKeys(valAddress, cosmosAddress, ethAddress)
	ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(blockHeight)
	_, err := h(ctx, msg)
	require.NoError(t, err)

	require.Equal(t, ethAddress, k.GetValidatorEthereumAddress(ctx, valAddress).Hex())
	require.Equal(t, valAddress, k.GetOrchestratorValidatorAddress(ctx, cosmosAddress))
	require.Equal(t, cosmosAddress, k.GetEthereumOrchestratorAddress(ctx, common.HexToAddress(ethAddress)))

	_, err = k.DelegateKeysByOrchestrator(wctx, &types.DelegateKeysByOrchestratorRequest{OrchestratorAddress: cosmosAddress.String()})
	require.NoError(t, err)

	_, err = k.DelegateKeysByEthereumSigner(wctx, &types.DelegateKeysByEthereumSignerRequest{EthereumSigner: ethAddress})
	require.NoError(t, err)

	_, err = k.DelegateKeysByValidator(wctx, &types.DelegateKeysByValidatorRequest{ValidatorAddress: valAddress.String()})
	require.NoError(t, err)

	// delegate new orch and eth addrs for same validator
	msg = types.NewMsgDelegateKeys(valAddress, cosmosAddress2, ethAddress2)
	ctx = ctx.WithBlockTime(blockTime2).WithBlockHeight(blockHeight2)
	_, err = h(ctx, msg)
	require.NoError(t, err)

	require.Equal(t, ethAddress2, k.GetValidatorEthereumAddress(ctx, valAddress).Hex())
	require.Equal(t, valAddress, k.GetOrchestratorValidatorAddress(ctx, cosmosAddress2))
	require.Equal(t, cosmosAddress2, k.GetEthereumOrchestratorAddress(ctx, common.HexToAddress(ethAddress2)))

	_, err = k.DelegateKeysByOrchestrator(wctx, &types.DelegateKeysByOrchestratorRequest{OrchestratorAddress: cosmosAddress2.String()})
	require.NoError(t, err)

	_, err = k.DelegateKeysByEthereumSigner(wctx, &types.DelegateKeysByEthereumSignerRequest{EthereumSigner: ethAddress2})
	require.NoError(t, err)

	_, err = k.DelegateKeysByValidator(wctx, &types.DelegateKeysByValidatorRequest{ValidatorAddress: valAddress.String()})
	require.NoError(t, err)
}
