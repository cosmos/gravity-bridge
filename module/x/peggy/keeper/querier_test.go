package keeper

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryValsetConfirm(t *testing.T) {
	var (
		nonce                 int64          = 1
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

func TestAllValsetConfirmsByNonce(t *testing.T) {
	var (
		nonce int64 = 1
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

func TestLastValsetRequestNonces(t *testing.T) {
	k, ctx, _ := CreateTestEnv(t)
	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			ethAddr := fmt.Sprintf("my eth addr %d", j+1)
			k.SetEthAddress(ctx, valAddr, ethAddr)
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
	"Nonce": "105",
	"Powers": [
	  "100",
	  "100",
	  "100",
	  "100",
	  "100",
	  "100"
	],
	"EthAddresses": [
	  "my eth addr 1",
	  "my eth addr 2",
	  "my eth addr 3",
	  "my eth addr 4",
	  "my eth addr 5",
	  "my eth addr 6"
	]
  },
  {
	"Nonce": "104",
	"Powers": [
	  "100",
	  "100",
	  "100",
	  "100",
	  "100"
	],
	"EthAddresses": [
	  "my eth addr 1",
	  "my eth addr 2",
	  "my eth addr 3",
	  "my eth addr 4",
	  "my eth addr 5"
	]
  },
  {
	"Nonce": "103",
	"Powers": [
	  "100",
	  "100",
	  "100",
	  "100"
	],
	"EthAddresses": [
	  "my eth addr 1",
	  "my eth addr 2",
	  "my eth addr 3",
	  "my eth addr 4"
	]
  },
  {
	"Nonce": "102",
	"Powers": [
	  "100",
	  "100",
	  "100"
	],
	"EthAddresses": [
	  "my eth addr 1",
	  "my eth addr 2",
	  "my eth addr 3"
	]
  },
  {
	"Nonce": "101",
	"Powers": [
	  "100",
	  "100"
	],
	"EthAddresses": [
	  "my eth addr 1",
	  "my eth addr 2"
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
		aValidatorCosmosAddr       = bytes.Repeat([]byte{1}, sdk.AddrLen)
		otherValidatorCosmosAddr   = bytes.Repeat([]byte{2}, sdk.AddrLen)
		unknownValidatorCosmosAddr = bytes.Repeat([]byte{3}, sdk.AddrLen)
	)
	k.StakingKeeper = NewStakingKeeperMock(aValidatorCosmosAddr, otherValidatorCosmosAddr)
	// seed with requests
	ctx = ctx.WithBlockHeight(200)
	k.SetValsetRequest(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	k.SetValsetConfirm(ctx, types.MsgValsetConfirm{Nonce: ctx.BlockHeight(), Validator: otherValidatorCosmosAddr})
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
	"Nonce": "201",
	"Powers": [
	  "100",
	  "100"
	],
	"EthAddresses": [
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
	"Nonce": "201",
	"Powers": [
	  "100",
	  "100"
	],
	"EthAddresses": [
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
