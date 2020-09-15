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
	k, ctx := CreateTestEnv(t)
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
	k, ctx := CreateTestEnv(t)

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
