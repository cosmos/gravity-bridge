package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestQueryValsetConfirm(t *testing.T) {
	var (
		nonce                                         = types.NewUInt64Nonce(1)
		myValidatorCosmosAddr   sdk.AccAddress        = make([]byte, sdk.AddrLen)
		myValidatorEthereumAddr types.EthereumAddress = createEthAddress(50)
	)
	k, ctx, _ := CreateTestEnv(t)
	k.SetValsetConfirm(ctx, types.MsgValsetConfirm{
		Nonce:     nonce,
		Validator: myValidatorCosmosAddr,
		Address:   myValidatorEthereumAddr,
	})

	specs := map[string]struct {
		srcNonce string
		srcAddr  string
		expErr   bool
		expResp  []byte
	}{
		"all good": {
			srcNonce: "1",
			srcAddr:  myValidatorCosmosAddr.String(),
			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",  "signature": ""}}`),
		},
		"unknown nonce": {
			srcNonce: "999999",
			srcAddr:  myValidatorCosmosAddr.String(),
		},
		"invalid address": {
			srcNonce: "1",
			srcAddr:  "not a valid addr",
			expErr:   true,
		},
		"invalid nonce": {
			srcNonce: "not a valid nonce",
			srcAddr:  myValidatorCosmosAddr.String(),
			expErr:   true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := queryValsetConfirm(ctx, []string{spec.srcNonce, spec.srcAddr}, k)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if spec.expResp == nil {
				assert.Nil(t, got)
				return
			}
			assert.JSONEq(t, string(spec.expResp), string(got))
		})
	}
}

func TestAllValsetConfirmsBynonce(t *testing.T) {
	var (
		nonce = types.NewUInt64Nonce(1)
	)
	k, ctx, _ := CreateTestEnv(t)

	// seed confirmations
	for i := 0; i < 3; i++ {
		addr := bytes.Repeat([]byte{byte(i)}, sdk.AddrLen)
		k.SetValsetConfirm(ctx, types.MsgValsetConfirm{
			Nonce:     nonce,
			Validator: addr,
			Address:   createEthAddress(50),
			Signature: fmt.Sprintf("signature %d", i+1),
		})
	}

	specs := map[string]struct {
		srcNonce string
		expErr   bool
		expResp  []byte
	}{
		"all good": {
			srcNonce: "1",
			expResp: []byte(`[
{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a", "signature": "signature 1"},
{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du", "signature": "signature 2"},
{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1qgpqyqszqgpqyqszqgpqyqszqgpqyqszrh8mx2", "signature": "signature 3"}
]`),
		},
		"unknown nonce": {
			srcNonce: "999999",
			expResp:  nil,
		},
		"invalid nonce": {
			srcNonce: "not a valid nonce",
			expErr:   true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := allValsetConfirmsByNonce(ctx, spec.srcNonce, k)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if spec.expResp == nil {
				assert.Nil(t, got)
				return
			}
			assert.JSONEq(t, string(spec.expResp), string(got))
		})
	}
}

func TestLastValsetRequestNonces(t *testing.T) {
	k, ctx, _ := CreateTestEnv(t)
	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			k.SetEthAddress(ctx, valAddr, createEthAddress(j+1))
			validators = append(validators, valAddr)
		}
		k.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		k.SetValsetRequest(ctx)
	}

	specs := map[string]struct {
		expResp []byte
	}{
		"limit at 5": {
			expResp: []byte(`[
{
  "nonce": "105",
  "members": [
    {
      "power": "715827882",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0505050505050505050505050505050505050505"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0606060606060606060606060606060606060606"
    }
  ]
},
{
  "nonce": "104",
  "members": [
    {
      "power": "858993459",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0505050505050505050505050505050505050505"
    }
  ]
},
{
  "nonce": "103",
  "members": [
    {
      "power": "1073741823",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    }
  ]
},
{
  "nonce": "102",
  "members": [
    {
      "power": "1431655765",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "1431655765",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "1431655765",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    }
  ]
},
{
  "nonce": "101",
  "members": [
    {
      "power": "2147483647",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "2147483647",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    }
  ]
}
]`),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := lastValsetRequests(ctx, k)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

func TestLastPendingValsetRequest(t *testing.T) {
	k, ctx, _ := CreateTestEnv(t)
	var (
		aValidatorCosmosAddr                      = bytes.Repeat([]byte{1}, sdk.AddrLen)
		otherValidatorCosmosAddr   sdk.ValAddress = bytes.Repeat([]byte{2}, sdk.AddrLen)
		unknownValidatorCosmosAddr                = bytes.Repeat([]byte{3}, sdk.AddrLen)
	)
	k.StakingKeeper = NewStakingKeeperMock(aValidatorCosmosAddr, otherValidatorCosmosAddr)
	// seed with requests
	ctx = ctx.WithBlockHeight(200)
	k.SetValsetRequest(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	nonce := types.NewUInt64Nonce(uint64(ctx.BlockHeight()))
	k.SetBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, nonce, otherValidatorCosmosAddr, "signature")
	k.SetValsetRequest(ctx)

	specs := map[string]struct {
		srcAddr string
		expErr  bool
		expResp []byte
	}{
		"last req when unsigned": {
			srcAddr: sdk.AccAddress(aValidatorCosmosAddr).String(),
			expResp: []byte(`
{
  "type": "peggy/Valset",
  "value": {
	"nonce": "201",
    "members": [
     {
       "power": "2147483647",
       "ethereum_address": ""
     },
     {
       "power": "2147483647",
       "ethereum_address": ""
     }
   ]
  }
}
`),
		},
		"empty when last req was signed": {
			srcAddr: sdk.AccAddress(otherValidatorCosmosAddr).String(),
			expResp: nil,
		},
		"not empty unknown source address": {
			srcAddr: sdk.AccAddress(unknownValidatorCosmosAddr).String(),
			expResp: []byte(`
{
  "type": "peggy/Valset",
  "value": {
	"nonce": "201",
    "members": [
     {
       "power": "2147483647",
       "ethereum_address": ""
     },
     {
       "power": "2147483647",
       "ethereum_address": ""
     }
   ]
  }
}
`),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := lastPendingValsetRequest(ctx, spec.srcAddr, k)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if spec.expResp == nil {
				assert.Nil(t, got)
				return
			}
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

func TestLastApprovedMultiSigUpdate(t *testing.T) {
	const maxVal = 3
	validators := make(map[string]types.BridgeValidator, maxVal)
	validatorAddrs := make([]sdk.ValAddress, maxVal+1)
	sigs := make(map[string]string, maxVal)

	// some validator/ orchestrator test data
	for i := 0; i < maxVal; i++ {
		n := i + 1
		validatorAddrs[n] = createValAddress()
		valAddr := validatorAddrs[n].String()
		validators[valAddr] = types.BridgeValidator{
			Power:           uint64(n),
			EthereumAddress: createEthAddress(n),
		}
		sigs[valAddr] = createFakeEthSignature(n)
	}
	unsortedOrchestrators := make(types.BridgeValidators, 0, maxVal)
	for _, v := range validators {
		unsortedOrchestrators = append(unsortedOrchestrators, v)
	}
	nonce := types.NewUInt64Nonce(1)

	specs := map[string]struct {
		srcBridgedVal types.BridgeValidators
		srcSigners    map[types.EthereumAddress]sdk.ValAddress
		expResp       *MultiSigUpdateResponse
	}{

		"validators and sigs ordered by power": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
				{Power: 2, EthereumAddress: createEthAddress(2)},
				{Power: 3, EthereumAddress: createEthAddress(3)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(1): validatorAddrs[1],
				createEthAddress(2): validatorAddrs[2],
				createEthAddress(3): validatorAddrs[3],
			},
			expResp: &MultiSigUpdateResponse{
				Valset: types.Valset{
					Nonce: nonce,
					Members: types.BridgeValidators{
						{Power: 3, EthereumAddress: createEthAddress(3)},
						{Power: 2, EthereumAddress: createEthAddress(2)},
						{Power: 1, EthereumAddress: createEthAddress(1)},
					},
				},
				Signatures: []SignatureWithAddress{
					{
						createFakeEthSignature(3),
						createEthAddress(3),
					},
					{
						createFakeEthSignature(2),
						createEthAddress(2),
					},
					{
						createFakeEthSignature(1),
						createEthAddress(1),
					},
				},
			},
		},
		"validators with power 0 excluded": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 0, EthereumAddress: createEthAddress(1)},
				{Power: 2, EthereumAddress: createEthAddress(2)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(2): validatorAddrs[2],
			},
			expResp: &MultiSigUpdateResponse{
				Valset: types.Valset{
					Nonce: nonce,
					Members: types.BridgeValidators{
						{Power: 2, EthereumAddress: createEthAddress(2)},
					},
				},
				Signatures: []SignatureWithAddress{
					{
						createFakeEthSignature(2),
						createEthAddress(2),
					},
				},
			},
		},
		"validators not in current mutliSig set excluded": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(1):   validatorAddrs[1],
				createEthAddress(999): validatorAddrs[0],
			},
			expResp: &MultiSigUpdateResponse{
				Valset: types.Valset{
					Nonce: nonce,
					Members: types.BridgeValidators{
						{Power: 1, EthereumAddress: createEthAddress(1)},
					},
				},
				Signatures: []SignatureWithAddress{
					{
						createFakeEthSignature(1),
						createEthAddress(1),
					},
				},
			},
		},
		"without signatures": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
			},
			expResp: &MultiSigUpdateResponse{
				Valset: types.Valset{Nonce: nonce, Members: types.BridgeValidators{{Power: 1, EthereumAddress: createEthAddress(1)}}},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			k, ctx, _ := CreateTestEnv(t)
			// persist multiSig set as latest observed
			k.storeValset(ctx, types.NewValset(nonce, spec.srcBridgedVal))
			k.setLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate, nonce)
			// store approved attestation
			k.setLastAttestedNonce(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, nonce)
			// with all the orchestrator's signatures stored
			for ethAddr, valAddr := range spec.srcSigners {
				k.SetEthAddress(ctx, valAddr, ethAddr)
				if sig, ok := sigs[valAddr.String()]; ok {
					k.SetBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, nonce, valAddr, sig)
				}
			}

			// when
			bz, err := lastApprovedMultiSigUpdate(ctx, k)

			// then
			require.NoError(t, err)
			if spec.expResp == nil {
				require.Nil(t, bz)
				return
			}
			require.NotNil(t, bz)

			var got MultiSigUpdateResponse
			k.cdc.MustUnmarshalJSON(bz, &got)

			assert.Equal(t, *spec.expResp, got)
		})
	}
}

func TestQueryInflightBatches(t *testing.T) {
	const maxVal = 3
	validators := make(map[string]types.BridgeValidator, maxVal)
	validatorAddrs := make([]sdk.ValAddress, maxVal+1)
	sigs := make(map[string]string, maxVal)

	// some validator/ orchestrator test data
	for i := 0; i < maxVal; i++ {
		n := i + 1
		validatorAddrs[n] = createValAddress()
		valAddr := validatorAddrs[n].String()
		validators[valAddr] = types.BridgeValidator{
			Power:           uint64(n),
			EthereumAddress: createEthAddress(n),
		}
		sigs[valAddr] = createFakeEthSignature(n)
	}
	unsortedOrchestrators := make(types.BridgeValidators, 0, maxVal)
	for _, v := range validators {
		unsortedOrchestrators = append(unsortedOrchestrators, v)
	}
	var (
		nonce = types.NewUInt64Nonce(1)
		batch = types.OutgoingTxBatch{Nonce: nonce}
	)

	specs := map[string]struct {
		srcBridgedVal     types.BridgeValidators
		srcSigners        map[types.EthereumAddress]sdk.ValAddress
		lastObservedBatch types.UInt64Nonce
		expResp           []ApprovedOutgoingTxBatchResponse
	}{
		"validators and sigs ordered by power": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
				{Power: 2, EthereumAddress: createEthAddress(2)},
				{Power: 3, EthereumAddress: createEthAddress(3)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(1): validatorAddrs[1],
				createEthAddress(2): validatorAddrs[2],
				createEthAddress(3): validatorAddrs[3],
			},
			lastObservedBatch: types.NewUInt64Nonce(0),

			expResp: []ApprovedOutgoingTxBatchResponse{
				{
					Batch: batch,
					Signatures: []SignatureWithAddress{
						{
							createFakeEthSignature(3),
							createEthAddress(3),
						},
						{
							createFakeEthSignature(2),
							createEthAddress(2),
						},
						{
							createFakeEthSignature(1),
							createEthAddress(1),
						},
					},
				},
			},
		},
		"validators with power 0 excluded": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 0, EthereumAddress: createEthAddress(1)},
				{Power: 2, EthereumAddress: createEthAddress(2)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(2): validatorAddrs[2],
			},
			expResp: []ApprovedOutgoingTxBatchResponse{
				{
					Batch: batch,
					Signatures: []SignatureWithAddress{
						{
							createFakeEthSignature(2),
							createEthAddress(2),
						},
					},
				},
			},
		},
		"validators not in current mutliSig set excluded": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
			},
			srcSigners: map[types.EthereumAddress]sdk.ValAddress{
				createEthAddress(1):   validatorAddrs[1],
				createEthAddress(999): validatorAddrs[0],
			},
			expResp: []ApprovedOutgoingTxBatchResponse{
				{
					Batch: batch,
					Signatures: []SignatureWithAddress{
						{
							createFakeEthSignature(1),
							createEthAddress(1),
						},
					},
				},
			},
		},
		"without signatures": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
			},
			expResp: []ApprovedOutgoingTxBatchResponse{
				{Batch: batch},
			},
		},
		"nothing in flight": {
			srcBridgedVal: []types.BridgeValidator{
				{Power: 1, EthereumAddress: createEthAddress(1)},
			},
			lastObservedBatch: nonce,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			k, ctx, _ := CreateTestEnv(t)
			// persist multiSig set
			k.storeValset(ctx, types.NewValset(nonce, spec.srcBridgedVal))
			k.setLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate, nonce)
			// persist batch
			k.storeBatch(ctx, batch)
			// store an approved attestation
			if !spec.lastObservedBatch.IsEmpty() {
				k.setLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeWithdrawalBatch, spec.lastObservedBatch)
			}
			// with all the orchestrator's signatures stored
			for ethAddr, valAddr := range spec.srcSigners {
				k.SetEthAddress(ctx, valAddr, ethAddr)
				if sig, ok := sigs[valAddr.String()]; ok {
					k.SetBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedWithdrawBatch, nonce, valAddr, sig)
				}
			}

			// when
			bz, err := queryInflightBatches(ctx, k)

			// then
			require.NoError(t, err)
			if spec.expResp == nil {
				require.Nil(t, bz)
				return
			}
			require.NotNil(t, bz)

			var got []ApprovedOutgoingTxBatchResponse
			k.cdc.MustUnmarshalJSON(bz, &got)

			assert.Equal(t, spec.expResp, got)
		})
	}
}

func createEthAddress(i int) types.EthereumAddress {
	return types.EthereumAddress(gethCommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20)))
}

func createFakeEthSignature(n int) string {
	return hex.EncodeToString(bytes.Repeat([]byte{byte(n)}, 64))
}

func createValAddress() sdk.ValAddress {
	return sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
}
