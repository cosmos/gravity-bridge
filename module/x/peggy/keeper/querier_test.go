package keeper

import (
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
			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"nonce": "1", "validator": "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",  "signature": null}}`),
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
