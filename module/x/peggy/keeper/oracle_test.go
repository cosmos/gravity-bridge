package keeper

import (
	"bytes"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObserveDeposit(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr       sdk.AccAddress = bytes.Repeat([]byte{1}, 12)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myNonce                           = bytes.Repeat([]byte{2}, 12)
		anyETHAddr                        = types.NewEthereumAddress("any-address")
		tokenETHAddr                      = types.NewEthereumAddress("any-erc20-token-addr")
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	depositDetails := types.BridgeDeposit{
		ERC20Token: types.ERC20Token{
			Amount:               12,
			Symbol:               "ALX",
			TokenContractAddress: tokenETHAddr,
		},
		EthereumSender: anyETHAddr,
		CosmosReceiver: myCosmosAddr,
	}
	gotAttestation, err := k.AddClaim(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce, myValAddr, depositDetails)
	// then
	require.NoError(t, err)

	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce, myValAddr, depositDetails)
	assert.True(t, claimFound)

	// and expected state
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeEthereumBridgeDeposit,
		Nonce:         myNonce,
		Certainty:     types.CertaintyObserved,
		Status:        types.ProcessStatusProcessed,
		ProcessResult: types.ProcessResultSuccess,
		Tally: types.AttestationTally{
			TotalVotesPower:    sdk.NewUint(100),
			TotalVotesCount:    1,
			RequiredVotesPower: sdk.NewUint(66),
			RequiredVotesCount: 0,
		},
		SubmitTime:          myBlockTime,
		ConfirmationEndTime: time.Date(2020, 9, 14+1, 15, 20, 10, 0, time.UTC),
		Details:             depositDetails,
	}
	assert.Equal(t, exp, *gotAttestation)
}

func TestObserveWithdrawBatch(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myBatchID          uint64         = 1
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)
	k.storeBatch(ctx, myBatchID, types.OutgoingTxBatch{
		BatchStatus: types.BatchStatusSubmitted,
	})

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	myNonce := types.NonceFromUint64(myBatchID)
	gotAttestation, err := k.AddClaim(ctx, types.ClaimTypeEthereumBridgeWithdrawalBatch, myNonce, myValAddr, nil)
	// then
	require.NoError(t, err)

	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBridgeWithdrawalBatch, myNonce, myValAddr, nil)
	assert.True(t, claimFound)

	// and expected state
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeEthereumBridgeWithdrawalBatch,
		Nonce:         myNonce,
		Certainty:     types.CertaintyObserved,
		Status:        types.ProcessStatusProcessed,
		ProcessResult: types.ProcessResultSuccess,
		Tally: types.AttestationTally{
			TotalVotesPower:    sdk.NewUint(100),
			TotalVotesCount:    1,
			RequiredVotesPower: sdk.NewUint(66),
			RequiredVotesCount: 0,
		},
		SubmitTime:          myBlockTime,
		ConfirmationEndTime: time.Date(2020, 9, 14+1, 15, 20, 10, 0, time.UTC),
	}
	assert.Equal(t, exp, *gotAttestation)
	// and last observed status updated
	assert.Equal(t, myBatchID, k.GetLastObservedBatchID(ctx))
}
