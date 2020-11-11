package keeper

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// func TestQueryValsetConfirm(t *testing.T) {
// 	var (
// 		nonce                                         = types.NewUInt64Nonce(1)
// 		myValidatorCosmosAddr   sdk.AccAddress        = make([]byte, sdk.AddrLen)
// 		myValidatorEthereumAddr types.EthereumAddress = createEthAddress(50)
// 	)
// 	k, ctx, _ := CreateTestEnv(t)
// 	k.SetValsetConfirm(ctx, types.MsgValsetConfirm{
// 		Nonce:     nonce,
// 		Validator: myValidatorCosmosAddr,
// 		Address:   myValidatorEthereumAddr,
// 	})

// 	specs := map[string]struct {
// 		srcNonce string
// 		srcAddr  string
// 		expErr   bool
// 		expResp  []byte
// 	}{
// 		"all good": {
// 			srcNonce: "1",
// 			srcAddr:  myValidatorCosmosAddr.String(),
// 			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",  "signature": ""}}`),
// 		},
// 		"unknown nonce": {
// 			srcNonce: "999999",
// 			srcAddr:  myValidatorCosmosAddr.String(),
// 		},
// 		"invalid address": {
// 			srcNonce: "1",
// 			srcAddr:  "not a valid addr",
// 			expErr:   true,
// 		},
// 		"invalid nonce": {
// 			srcNonce: "not a valid nonce",
// 			srcAddr:  myValidatorCosmosAddr.String(),
// 			expErr:   true,
// 		},
// 	}
// 	for msg, spec := range specs {
// 		t.Run(msg, func(t *testing.T) {
// 			got, err := queryValsetConfirm(ctx, []string{spec.srcNonce, spec.srcAddr}, k)
// 			if spec.expErr {
// 				require.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)
// 			if spec.expResp == nil {
// 				assert.Nil(t, got)
// 				return
// 			}
// 			assert.JSONEq(t, string(spec.expResp), string(got))
// 		})
// 	}
// }

func TestAllValsetConfirmsBynonce(t *testing.T) {
	var (
		nonce = types.NewUInt64Nonce(1)
	)
	k, ctx, _ := CreateTestEnv(t)

	// seed confirmations
	for i := 0; i < 3; i++ {
		addr := bytes.Repeat([]byte{byte(i)}, sdk.AddrLen)
		k.SetValsetApprovalSignature(ctx, nonce, addr, []byte(fmt.Sprintf("signature %d", i+1)))
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

func createEthAddress(i int) types.EthereumAddress {
	return types.EthereumAddress(gethCommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20)))
}

func createFakeEthSignature(n int) []byte {
	return bytes.Repeat([]byte{byte(n)}, 64)
}

func createValAddress() sdk.ValAddress {
	return sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
}
