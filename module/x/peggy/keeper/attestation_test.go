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
		myNonce                           = types.NewUInt64Nonce(123)
		anyETHAddr                        = types.NewEthereumAddress("any-address")
		tokenETHAddr                      = types.NewEthereumAddress("any-erc20-token-addr")
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)

	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	depositDetails := types.BridgeDeposit{
		Nonce: myNonce,
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
	// and last observed status updated
	gotAttestation = k.GetLastObservedAttestation(ctx, types.ClaimTypeEthereumBridgeDeposit)
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
	myNonce := types.NewUInt64Nonce(myBatchID)
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
	gotAttestation = k.GetLastObservedAttestation(ctx, types.ClaimTypeEthereumBridgeWithdrawalBatch)
	assert.Equal(t, exp, *gotAttestation)
}

func TestObserveBridgeMultiSigUpdate(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myBlockHeight      uint64         = 100
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)

	ctx = ctx.WithBlockTime(myBlockTime).WithBlockHeight(int64(myBlockHeight))
	k.SetValsetRequest(ctx)

	// when
	myNonce := types.NewUInt64Nonce(myBlockHeight)
	gotAttestation, err := k.AddClaim(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate, myNonce, myValAddr, nil)
	// then
	require.NoError(t, err)

	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate, myNonce, myValAddr, nil)
	assert.True(t, claimFound)

	// and expected state
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeEthereumBridgeMultiSigUpdate,
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
	gotAttestation = k.GetLastObservedAttestation(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate)
	assert.Equal(t, exp, *gotAttestation)
	gotNonce := k.GetLastValsetObservedNonce(ctx)
	require.NotNil(t, gotNonce)
	assert.Equal(t, myNonce, *gotNonce)
}

func TestObserveBridgeBootstrap(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myEthAddr                         = types.NewEthereumAddress("0xa")
		myBlockHeight      uint64         = 100
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)
	ctx = ctx.WithBlockTime(myBlockTime).WithBlockHeight(int64(myBlockHeight))

	// when
	myNonce := types.NewUInt64Nonce(myBlockHeight)
	details := types.BridgeBootstrap{
		AllowedValidatorSet: []types.EthereumAddress{myEthAddr},
		ValidatorPowers:     []uint64{10},
		PeggyID:             []byte("my random string"),
		StartThreshold:      67,
	}
	gotAttestation, err := k.AddClaim(ctx, types.ClaimTypeEthereumBootstrap, myNonce, myValAddr, details)
	// then
	require.NoError(t, err)

	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBootstrap, myNonce, myValAddr, details)
	assert.True(t, claimFound)

	// and expected state
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeEthereumBootstrap,
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
		Details:             details,
		SubmitTime:          myBlockTime,
		ConfirmationEndTime: time.Date(2020, 9, 14+1, 15, 20, 10, 0, time.UTC),
	}
	assert.Equal(t, exp, *gotAttestation)
	// and last observed status updated
	gotAttestation = k.GetLastObservedAttestation(ctx, types.ClaimTypeEthereumBootstrap)
	assert.Equal(t, exp, *gotAttestation)
	gotNonce := k.GetLastValsetObservedNonce(ctx)
	require.NotNil(t, gotNonce)
}

func TestApproveBridgeMultiSigUpdate(t *testing.T) {
	var (
		myOrchestratorCosmosAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myOrchestratorETHAddr                   = types.NewEthereumAddress("0x8858eeb3dfffa017d4bce9801d340d36cf895ccf")
		myValAddr                               = sdk.ValAddress(myOrchestratorCosmosAddr)
		myBlockHeight            uint64         = 100
		myBlockTime                             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	k, ctx, _ := CreateTestEnv(t)
	k.StakingKeeper = NewStakingKeeperMock(myValAddr)

	ctx = ctx.WithBlockTime(myBlockTime).WithBlockHeight(int64(myBlockHeight))
	k.SetEthAddress(ctx, myValAddr, myOrchestratorETHAddr)
	k.SetValsetRequest(ctx)

	// when
	myNonce := types.NewUInt64Nonce(myBlockHeight)
	checkpoint := k.GetValsetRequest(ctx, int64(myBlockHeight)).GetCheckpoint()
	details := types.SignedCheckpoint{
		Checkpoint: checkpoint,
	}
	gotAttestation, err := k.AddClaim(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, myNonce, myValAddr, details)
	// then
	require.NoError(t, err)

	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, myNonce, myValAddr, details)
	assert.True(t, claimFound)

	// and expected state
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeOrchestratorSignedMultiSigUpdate,
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
		Details:             types.SignedCheckpoint{Checkpoint: checkpoint},
		SubmitTime:          myBlockTime,
		ConfirmationEndTime: time.Date(2020, 9, 14+1, 15, 20, 10, 0, time.UTC),
	}
	assert.Equal(t, exp, *gotAttestation)
	// and last observed status updated
	gotAttestation = k.GetLastObservedAttestation(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate)
	assert.Equal(t, exp, *gotAttestation)

	gotNonce := k.GetLastValsetApprovedNonce(ctx)
	require.NotNil(t, gotNonce)
	assert.Equal(t, myNonce, *gotNonce)
}
