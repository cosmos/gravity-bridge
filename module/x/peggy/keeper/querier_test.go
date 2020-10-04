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
)

func TestQueryValsetConfirm(t *testing.T) {
	var (
		nonce                                = types.NewUInt64Nonce(1)
		myValidatorCosmosAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
	)
	k, ctx, _ := CreateTestEnv(t)
	k.SetValsetConfirm(ctx, types.MsgValsetConfirm{
		Nonce:     nonce,
		Validator: myValidatorCosmosAddr,
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
			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",  "signature": ""}}`),
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
{"nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a", "signature": "signature 1"},
{"nonce": "1", "validator": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du", "signature": "signature 2"},
{"nonce": "1", "validator": "cosmos1qgpqyqszqgpqyqszqgpqyqszqgpqyqszrh8mx2", "signature": "signature 3"}
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

func TestLastValsetRequestnonces(t *testing.T) {
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
	"powers": [
	  "715827882",
	  "715827882",
	  "715827882",
	  "715827882",
	  "715827882",
	  "715827882"
	],
	"eth_addresses": [
	  "0x0101010101010101010101010101010101010101",
	  "0x0202020202020202020202020202020202020202",
	  "0x0303030303030303030303030303030303030303",
	  "0x0404040404040404040404040404040404040404",
	  "0x0505050505050505050505050505050505050505",
	  "0x0606060606060606060606060606060606060606"
	]
  },
  {
	"nonce": "104",
	"powers": [
	  "858993459",
	  "858993459",
	  "858993459",
	  "858993459",
	  "858993459"
	],
	"eth_addresses": [
	  "0x0101010101010101010101010101010101010101",
	  "0x0202020202020202020202020202020202020202",
	  "0x0303030303030303030303030303030303030303",
	  "0x0404040404040404040404040404040404040404",
	  "0x0505050505050505050505050505050505050505"
	]
  },
  {
	"nonce": "103",
	"powers": [
	  "1073741823",
	  "1073741823",
	  "1073741823",
	  "1073741823"
	],
	"eth_addresses": [
	  "0x0101010101010101010101010101010101010101",
	  "0x0202020202020202020202020202020202020202",
	  "0x0303030303030303030303030303030303030303",
	  "0x0404040404040404040404040404040404040404"
	]
  },
  {
	"nonce": "102",
	"powers": [
	  "1431655765",
	  "1431655765",
	  "1431655765"
	],
	"eth_addresses": [
	  "0x0101010101010101010101010101010101010101",
	  "0x0202020202020202020202020202020202020202",
	  "0x0303030303030303030303030303030303030303"
	]
  },
  {
	"nonce": "101",
	"powers": [
	  "2147483647",
	  "2147483647"
	],
	"eth_addresses": [
	  "0x0101010101010101010101010101010101010101",
	  "0x0202020202020202020202020202020202020202"
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
	k.SetBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, nonce, otherValidatorCosmosAddr, []byte("signature"))
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
	"powers": [
	  "2147483647",
	  "2147483647"
	],
	"eth_addresses": [
	  "",
	  ""
	]
  }
}
`),
		},
		"empty when last req was signed": {
			srcAddr: sdk.AccAddress(otherValidatorCosmosAddr).String(),
			expResp: nil,
		},
		"not empty unknown address": {
			srcAddr: sdk.AccAddress(unknownValidatorCosmosAddr).String(),
			expResp: []byte(`
{
  "type": "peggy/Valset",
  "value": {
	"nonce": "201",
	"powers": [
	  "2147483647",
	  "2147483647"
	],
	"eth_addresses": [
	  "",
	  ""
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

func createEthAddress(i int) types.EthereumAddress {
	return types.EthereumAddress(gethCommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20)))

}
