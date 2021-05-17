package gravity

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/gravity-bridge/module/x/gravity/keeper"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestSignerSetTxCreationIfNotAvailable(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	gravity := input.GravityKeeper

	// EndBlocker should set a new validator set if not available
	EndBlocker(ctx, gravity)
	otx := gravity.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(1))
	require.NotNil(t, otx)
	_, ok := otx.(*types.SignerSetTx)
	require.True(t, ok)
	require.True(t, len(gravity.GetSignerSetTxs(ctx)) == 1)
}

func TestSignerSetTxCreationUponUnbonding(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	gravity := input.GravityKeeper
	gravity.NewSignerSetTx(ctx)

	input.Context = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// begin unbonding
	sh := staking.NewHandler(input.StakingKeeper)
	undelegateMsg := keeper.NewTestMsgUnDelegateValidator(keeper.ValAddrs[0], keeper.StakingAmount)
	sh(input.Context, undelegateMsg)

	// Run the staking endblocker to ensure valset is set in state
	staking.EndBlocker(input.Context, input.StakingKeeper)
	EndBlocker(input.Context, gravity)

	assert.Equal(t, 2, gravity.GetLatestSignerSetTxNonce(ctx))
}

func TestSignerSetTxSlashing_SignerSetTxCreated_Before_ValidatorBonded(t *testing.T) {
	//	Don't slash validators if valset is created before he is bonded.

	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.GravityKeeper
	params := input.GravityKeeper.GetParams(ctx)

	signerSet := pk.NewSignerSetTx(ctx)
	height := uint64(ctx.BlockHeight()) - (params.SignedSignerSetTxsWindow + 1)
	signerSet.Height = height
	pk.SetOutgoingTx(ctx, signerSet)

	EndBlocker(ctx, pk)

	// ensure that the  validator who is bonded after valset is created is not slashed
	val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	require.False(t, val.IsJailed())
}

func TestSignerSetTxSlashing_SignerSetTxCreated_After_ValidatorBonded(t *testing.T) {
	//	Slashing Conditions for Bonded Validator

	input, ctx := keeper.SetupFiveValChain(t)
	pk := input.GravityKeeper
	params := input.GravityKeeper.GetParams(ctx)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + int64(params.SignedSignerSetTxsWindow) + 2)
	signerSet := pk.NewSignerSetTx(ctx)
	height := uint64(ctx.BlockHeight()) - (params.SignedSignerSetTxsWindow + 1)
	signerSet.Height = height
	pk.SetOutgoingTx(ctx, signerSet)

	for i, val := range keeper.ValAddrs {
		if i == 0 {
			// don't sign with first validator
			continue
		}
		pk.SetEthereumSignature(ctx, &types.SignerSetTxSignature{signerSet.Nonce, keeper.AccAddrs[i].String(), []byte("dummysig")}, val)
	}

	EndBlocker(ctx, pk)

	// ensure that the  validator who is bonded before valset is created is slashed
	val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	require.True(t, val.IsJailed())

	// ensure that the  validator who attested the valset is not slashed.
	val = input.StakingKeeper.Validator(ctx, keeper.ValAddrs[1])
	require.False(t, val.IsJailed())

}

func TestValsetSlashing_UnbondingValidator_UnbondWindow_NotExpired(t *testing.T) {
	//	Slashing Conditions for Unbonding Validator

	//  Create 5 validators
	input, ctx := keeper.SetupFiveValChain(t)
	// val := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	// fmt.Println("val1  tokens", val.GetTokens().ToDec())

	gravity := input.GravityKeeper
	params := input.GravityKeeper.GetParams(ctx)

	// Define slashing variables
	validatorStartHeight := ctx.BlockHeight()                                                             // 0
	signerSetTxHeight := validatorStartHeight + 1                                                         // 1
	valUnbondingHeight := signerSetTxHeight + 1                                                           // 2
	signerSetTxSlashedAt := signerSetTxHeight + int64(params.SignedSignerSetTxsWindow)                    // 11
	validatorUnbondingWindowExpiry := valUnbondingHeight + int64(params.UnbondSlashingSignerSetTxsWindow) // 17
	currentBlockHeight := signerSetTxSlashedAt + 1                                                        // 12

	assert.True(t, signerSetTxSlashedAt < currentBlockHeight)
	assert.True(t, signerSetTxHeight < validatorUnbondingWindowExpiry)

	// Create Valset request
	ctx = ctx.WithBlockHeight(signerSetTxHeight)
	vs := gravity.NewSignerSetTx(ctx)
	vs.Height = uint64(signerSetTxHeight)
	vs.Nonce = uint64(signerSetTxHeight)
	gravity.SetOutgoingTx(ctx, vs)

	// Start Unbonding validators
	// Validator-1  Unbond slash window is not expired. if not attested, slash
	// Validator-2  Unbond slash window is not expired. if attested, don't slash
	input.Context = ctx.WithBlockHeight(valUnbondingHeight)
	sh := staking.NewHandler(input.StakingKeeper)
	undelegateMsg1 := keeper.NewTestMsgUnDelegateValidator(keeper.ValAddrs[0], keeper.StakingAmount)
	sh(input.Context, undelegateMsg1)
	undelegateMsg2 := keeper.NewTestMsgUnDelegateValidator(keeper.ValAddrs[1], keeper.StakingAmount)
	sh(input.Context, undelegateMsg2)

	for i, val := range keeper.ValAddrs {
		if i == 0 {
			// don't sign with first validator
			continue
		}
		gravity.SetEthereumSignature(ctx, &types.SignerSetTxSignature{vs.Nonce, keeper.EthAddrs[i].Hex(), []byte("dummySig")}, val)
	}
	staking.EndBlocker(input.Context, input.StakingKeeper)

	ctx = ctx.WithBlockHeight(currentBlockHeight)
	EndBlocker(ctx, gravity)

	// Assertions
	val1 := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0])
	assert.True(t, val1.IsJailed())
	fmt.Println("val1  tokens", val1.GetTokens().ToDec())
	// check if tokens are slashed for val1.

	val2 := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[1])
	assert.True(t, val2.IsJailed())
	fmt.Println("val2  tokens", val2.GetTokens().ToDec())
	// check if tokens shouldn't be slashed for val2.
}

func TestBatchSlashing(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	gravity := input.GravityKeeper
	params := gravity.GetParams(ctx)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + int64(params.SignedBatchesWindow) + 2)

	// First store a batch
	batch := &types.BatchTx{
		Nonce:         1,
		Transactions:  []*types.SendToEthereum{},
		TokenContract: keeper.TokenContractAddrs[0],
		Height:        uint64(ctx.BlockHeight() - int64(params.SignedBatchesWindow+1)),
	}
	gravity.SetOutgoingTx(ctx, batch)

	for i, val := range keeper.ValAddrs {
		if i == 0 {
			// don't sign with first validator
			continue
		}
		if i == 1 {
			// don't sign with 2nd validator. set val bond height > batch block height
			validator := input.StakingKeeper.Validator(ctx, keeper.ValAddrs[i])
			valConsAddr, _ := validator.GetConsAddr()
			valSigningInfo := slashingtypes.ValidatorSigningInfo{StartHeight: int64(batch.Height + 1)}
			input.SlashingKeeper.SetValidatorSigningInfo(ctx, valConsAddr, valSigningInfo)
			continue
		}
		gravity.SetEthereumSignature(ctx, &types.BatchTxSignature{
			Nonce:          batch.Nonce,
			TokenContract:  keeper.TokenContractAddrs[0],
			EthereumSigner: keeper.EthAddrs[i].String(),
			Signature:      []byte("dummysig"),
		}, val)
	}

	EndBlocker(ctx, gravity)

	// ensure that the  validator is jailed and slashed
	require.True(t, input.StakingKeeper.Validator(ctx, keeper.ValAddrs[0]).IsJailed())

	// ensure that the 2nd  validator is not jailed and slashed
	require.False(t, input.StakingKeeper.Validator(ctx, keeper.ValAddrs[1]).IsJailed())

	// Ensure that the last slashed valset nonce is set properly
	assert.Equal(t, input.GravityKeeper.GetLastSlashedBatchBlock(ctx), batch.Height)

}

func TestSignerSetTxEmission(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	gravity := input.GravityKeeper

	// Store a validator set with a power change as the most recent validator set
	sstx := gravity.NewSignerSetTx(ctx)
	// TODO: decrement height
	delta := float64(types.EthereumSigners(sstx.Signers).TotalPower()) * 0.05
	sstx.Signers[0].Power = uint64(float64(sstx.Signers[0].Power) - delta/2)
	sstx.Signers[1].Power = uint64(float64(sstx.Signers[1].Power) + delta/2)
	gravity.SetOutgoingTx(ctx, sstx)

	// EndBlocker should set a new validator set
	EndBlocker(ctx, gravity)
	require.NotNil(t, gravity.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(2)))
	require.True(t, len(gravity.GetSignerSetTxs(ctx)) == 2)
}

func TestSignerSetTxSetting(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	input.GravityKeeper.NewSignerSetTx(ctx)
	require.True(t, len(input.GravityKeeper.GetSignerSetTxs(ctx)) == 1)
}

/// Test batch timeout
func TestBatchTxTimeout(t *testing.T) {
	input, ctx := keeper.SetupFiveValChain(t)
	gravity := input.GravityKeeper
	params := gravity.GetParams(ctx)
	var (
		now                 = time.Now().UTC()
		mySender, _         = sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
		myReceiver          = common.HexToAddress("0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7")
		myTokenContractAddr = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5") // Pickle
		allVouchers         = sdk.NewCoins(types.NewERC20Token(99999, myTokenContractAddr.Hex()).GravityCoin())
	)

	require.Greater(t, params.AverageBlockTime, uint64(0))
	require.Greater(t, params.AverageEthereumBlockTime, uint64(0))

	// mint some vouchers first
	require.NoError(t, input.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	require.NoError(t, input.BankKeeper.SetBalances(ctx, mySender, allVouchers))

	// add some TX to the pool
	input.AddSendToEthTxsToPool(t, ctx, myTokenContractAddr, mySender, myReceiver, 2, 3, 2, 1, 5, 6)

	// when
	ctx = ctx.WithBlockTime(now).WithBlockHeight(250)

	// check that we can make a batch without first setting an ethereum block height
	b1, err1 := gravity.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err1)
	require.Equal(t, b1.Timeout, uint64(0))

	gravity.SetLastObservedEthereumBlockHeight(ctx, 500)

	b2, err2 := gravity.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err2)
	// this is exactly block 500 plus twelve hours
	require.Equal(t, b2.Timeout, uint64(504))

	// make sure the batches got stored in the first place
	gotFirstBatch := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b1.TokenContract), b1.Nonce))
	require.NotNil(t, gotFirstBatch)
	gotSecondBatch := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b2.TokenContract), b2.Nonce))
	require.NotNil(t, gotSecondBatch)

	// when, way into the future
	ctx = ctx.WithBlockTime(now).WithBlockHeight(9)

	b3, err2 := gravity.BuildBatchTx(ctx, myTokenContractAddr, 2)
	require.NoError(t, err2)

	EndBlocker(ctx, gravity)

	// this had a timeout of zero should be deleted.
	gotFirstBatch = input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b1.TokenContract), b1.Nonce))
	require.Nil(t, gotFirstBatch)
	// make sure the end blocker does not delete these, as the block height has not officially
	// been updated by a relay event
	gotSecondBatch = input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b2.TokenContract), b2.Nonce))
	require.NotNil(t, gotSecondBatch)
	gotThirdBatch := input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b3.TokenContract), b3.Nonce))
	require.NotNil(t, gotThirdBatch)

	gravity.SetLastObservedEthereumBlockHeight(ctx, 5000)
	EndBlocker(ctx, gravity)

	// make sure the end blocker does delete these, as we've got a new Ethereum block height
	gotFirstBatch = input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b1.TokenContract), b1.Nonce))
	require.Nil(t, gotFirstBatch)
	gotSecondBatch = input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b2.TokenContract), b2.Nonce))
	require.Nil(t, gotSecondBatch)
	gotThirdBatch = input.GravityKeeper.GetOutgoingTx(ctx, types.MakeBatchTxKey(common.HexToAddress(b3.TokenContract), b3.Nonce))
	require.NotNil(t, gotThirdBatch)
}
