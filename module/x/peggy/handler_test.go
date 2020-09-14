package peggy

import (
	"bytes"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCreateEthereumClaims(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr       sdk.AccAddress = bytes.Repeat([]byte{1}, 12)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myNonce                           = bytes.Repeat([]byte{2}, 12)
		anyETHAddr                        = types.NewEthereumAddress("any-address")
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	k, ctx := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	h := NewHandler(k)

	msg := MsgCreateEthereumClaims{
		EthereumChainID:       "0",
		BridgeContractAddress: types.NewEthereumAddress(""),
		Orchestrator:          myOrchestratorAddr,
		Claims: []EthereumClaim{
			EthereumBridgeDepositClaim{
				Nonce:          myNonce,
				ERC20Token:     types.ERC20Token{},
				EthereumSender: anyETHAddr,
				CosmosReceiver: myCosmosAddr,
			},
		},
	}
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err := h(ctx, msg)
	// then
	require.NoError(t, err)
	// and no
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce, myValAddr)
	assert.True(t, claimFound)
	// and attestation persisted
	a := k.GetAttestation(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce)
	require.NotNil(t, a)
	exp := types.Attestation{
		ClaimType:     types.ClaimTypeEthereumBridgeDeposit,
		Nonce:         myNonce,
		Certainty:     types.CertaintyObserved,
		Status:        types.ProcessStatusProcessed,
		ProcessResult: types.ProcessResultSuccess,
		Tally: types.AttestationTally{
			TotalVotesPower:    100,
			TotalVotesCount:    1,
			RequiredVotesPower: 66,
			RequiredVotesCount: 0,
		},
		SubmitTime:          myBlockTime,
		ConfirmationEndTime: time.Date(2020, 9, 14+1, 15, 20, 10, 0, time.UTC),
	}
	assert.Equal(t, exp, *a)
}
