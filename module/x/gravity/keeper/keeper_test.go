package keeper

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestPrefixRange(t *testing.T) {
	cases := map[string]struct {
		src      []byte
		expStart []byte
		expEnd   []byte
		expPanic bool
	}{
		"normal":              {src: []byte{1, 3, 4}, expStart: []byte{1, 3, 4}, expEnd: []byte{1, 3, 5}},
		"normal short":        {src: []byte{79}, expStart: []byte{79}, expEnd: []byte{80}},
		"empty case":          {src: []byte{}},
		"roll-over example 1": {src: []byte{17, 28, 255}, expStart: []byte{17, 28, 255}, expEnd: []byte{17, 29, 0}},
		"roll-over example 2": {src: []byte{15, 42, 255, 255},
			expStart: []byte{15, 42, 255, 255}, expEnd: []byte{15, 43, 0, 0}},
		"pathological roll-over": {src: []byte{255, 255, 255, 255}, expStart: []byte{255, 255, 255, 255}},
		"nil prohibited":         {expPanic: true},
	}

	for testName, tc := range cases {
		tc := tc
		t.Run(testName, func(t *testing.T) {
			if tc.expPanic {
				require.Panics(t, func() {
					prefixRange(tc.src)
				})
				return
			}
			start, end := prefixRange(tc.src)
			assert.Equal(t, tc.expStart, start)
			assert.Equal(t, tc.expEnd, end)
		})
	}
}

func TestCurrentSignerSetTxNormalization(t *testing.T) {
	specs := map[string]struct {
		srcPowers []uint64
		expPowers []uint64
	}{
		"one": {
			srcPowers: []uint64{100},
			expPowers: []uint64{4294967295},
		},
		"two": {
			srcPowers: []uint64{100, 1},
			expPowers: []uint64{4252442866, 42524428},
		},
	}
	input := CreateTestEnv(t)
	ctx := input.Context
	for msg, spec := range specs {
		spec := spec
		t.Run(msg, func(t *testing.T) {
			operators := make([]MockStakingValidatorData, len(spec.srcPowers))
			for i, v := range spec.srcPowers {
				operators[i] = MockStakingValidatorData{
					// any unique addr
					Operator: bytes.Repeat([]byte{byte(i)}, sdk.AddrLen),
					Power:    int64(v),
				}
			}
			input.GravityKeeper.StakingKeeper = NewStakingKeeperWeightedMock(operators...)
			r := input.GravityKeeper.GetCurrentSignerSetTx(ctx)
			assert.Equal(t, spec.expPowers, types.EthereumSigners(r.Members).GetPowers())
		})
	}
}

func TestEthereumEventVoteRecordIterator(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	// add some ethereumEventVoteRecords to the store

	voteRecord1 := &types.EthereumEventVoteRecord{
		Accepted: true,
		Votes:    []string{},
	}
	send1 := &types.MsgSendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
		Orchestrator:   AccAddrs[0].String(),
	}
	voteRecord2 := &types.EthereumEventVoteRecord{
		Accepted: true,
		Votes:    []string{},
	}
	send2 := &types.MsgSendToCosmosEvent{
		EventNonce:     2,
		TokenContract:  TokenContractAddrs[0],
		Amount:         sdk.NewInt(100),
		EthereumSender: EthAddrs[0].String(),
		CosmosReceiver: AccAddrs[0].String(),
		Orchestrator:   AccAddrs[0].String(),
	}
	input.GravityKeeper.SetEthereumEventVoteRecord(ctx, send1.EventNonce, send1.ClaimHash(), voteRecord1)
	input.GravityKeeper.SetEthereumEventVoteRecord(ctx, send2.EventNonce, send2.ClaimHash(), voteRecord2)

	voteRecords := []types.EthereumEventVoteRecord{}
	input.GravityKeeper.IterateEthereumVoteRecords(ctx, func(_ []byte, voteRecord types.EthereumEventVoteRecord) bool {
		voteRecords = append(voteRecords, voteRecord)
		return false
	})

	require.Len(t, voteRecords, 2)
}

func TestDelegateKeys(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	k := input.GravityKeeper
	var (
		ethAddrs = []string{"0x3146D2d6Eed46Afa423969f5dDC3152DfC359b09",
			"0x610277F0208D342C576b991daFdCb36E36515e76", "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4",
			"0xb2A7F3E84F8FdcA1da46c810AEa110dd96BAE6bF"}

		valAddrs = []string{"cosmosvaloper1jpz0ahls2chajf78nkqczdwwuqcu97w6z3plt4",
			"cosmosvaloper15n79nty2fj37ant3p2gj4wju4ls6eu6tjwmdt0", "cosmosvaloper16dnkc6ac6ruuyr6l372fc3p77jgjpet6fka0cq",
			"cosmosvaloper1vrptwhl3ht2txmzy28j9msqkcvmn8gjz507pgu"}

		orchAddrs = []string{"cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7", "cosmos1rhfs24tlw4na04v35tzmjncy785kkw9j27d5kx",
			"cosmos10upq3tmt04zf55f6hw67m0uyrda3mp722q70rw", "cosmos1nt2uwjh5peg9vz2wfh2m3jjwqnu9kpjlhgpmen"}
	)

	for i := range ethAddrs {
		// set some addresses
		val, err1 := sdk.ValAddressFromBech32(valAddrs[i])
		orch, err2 := sdk.AccAddressFromBech32(orchAddrs[i])
		require.NoError(t, err1)
		require.NoError(t, err2)
		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		k.SetEthAddressForValidator(ctx, val, ethAddrs[i])
	}

	addresses := k.GetDelegateKeys(ctx)
	for i := range addresses {
		res := addresses[i]
		assert.Equal(t, valAddrs[i], res.Validator)
		assert.Equal(t, orchAddrs[i], res.Orchestrator)
		assert.Equal(t, ethAddrs[i], res.EthAddress)
	}

}

func TestLastSlashedSignerSetTxNonce(t *testing.T) {
	input := CreateTestEnv(t)
	k := input.GravityKeeper
	ctx := input.Context

	vs := k.GetCurrentSignerSetTx(ctx)

	i := 1
	for ; i < 10; i++ {
		vs.Height = uint64(i)
		vs.Nonce = uint64(i)
		k.StoreSignerSetTxUnsafe(ctx, vs)
	}

	latestSignerSetTxNonce := k.GetLatestSignerSetTxNonce(ctx)
	assert.Equal(t, latestSignerSetTxNonce, uint64(i-1))

	//  lastSlashedSignerSetTxNonce should be zero initially.
	lastSlashedSignerSetTxNonce := k.GetLastSlashedSignerSetTxNonce(ctx)
	assert.Equal(t, lastSlashedSignerSetTxNonce, uint64(0))
	unslashedSignerSetTxs := k.GetUnSlashedSignerSetTxs(ctx, uint64(12))
	assert.Equal(t, len(unslashedSignerSetTxs), 9)

	// check if last Slashed SignerSetTx nonce is set properly or not
	k.SetLastSlashedSignerSetTxNonce(ctx, uint64(3))
	lastSlashedSignerSetTxNonce = k.GetLastSlashedSignerSetTxNonce(ctx)
	assert.Equal(t, lastSlashedSignerSetTxNonce, uint64(3))

	// when maxHeight < lastSlashedSignerSetTxNonce, len(unslashedSignerSetTxs) should be zero
	unslashedSignerSetTxs = k.GetUnSlashedSignerSetTxs(ctx, uint64(2))
	assert.Equal(t, len(unslashedSignerSetTxs), 0)

	// when maxHeight == lastSlashedSignerSetTxNonce, len(unslashedSignerSetTxs) should be zero
	unslashedSignerSetTxs = k.GetUnSlashedSignerSetTxs(ctx, uint64(3))
	assert.Equal(t, len(unslashedSignerSetTxs), 0)

	// when maxHeight > lastSlashedSignerSetTxNonce && maxHeight <= latestSignerSetTxNonce
	unslashedSignerSetTxs = k.GetUnSlashedSignerSetTxs(ctx, uint64(6))
	assert.Equal(t, len(unslashedSignerSetTxs), 2)

	// when maxHeight > latestSignerSetTxNonce
	unslashedSignerSetTxs = k.GetUnSlashedSignerSetTxs(ctx, uint64(15))
	assert.Equal(t, len(unslashedSignerSetTxs), 6)
	fmt.Println("unslashedSignerSetTxsRange", unslashedSignerSetTxs)
}
