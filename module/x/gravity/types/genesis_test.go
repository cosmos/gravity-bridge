package types

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateTestGenesisState() GenesisState {
	bridgeID := make([]byte, BridgeIDLen)
	rand.Read(bridgeID)

	return GenesisState{
		BridgeID: bridgeID,
		Params: DefaultParams(),
		LastObservedNonce: 10,
		SignerSets: []EthSignerSet{},
		BatchTxs: []BatchTx{},
		LogicCallTxs: []LogicCallTx{},
		TransferTxs: []TransferTx{},
		SignerSetConfirms: []ConfirmSignerSet{},
		BatchConfirms: []ConfirmBatch{},
		LogicCallConfirms: []ConfirmLogicCall{},
		Attestations: []Attestation{},
		DelegateKeys: []MsgDelegateKey{},
		Erc20ToDenoms: []ERC20ToDenom{},
	}
}

func TestGenesisState_ValidateBasic(t *testing.T) {
	bridgeID := make([]byte, BridgeIDLen)
	rand.Read(bridgeID)

	testCases := []struct {
		name string
		state GenesisState
		expError bool
	}{
		{"valid input", generateTestGenesisState(), false},
	}

	for _, tc := range testCases {
		err := tc.state.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}