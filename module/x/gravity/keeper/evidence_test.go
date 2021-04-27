package keeper

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestSubmitBadSignatureEvidence(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	inner_msg := types.OutgoingTxBatch{
		BatchTimeout: 420,
	}

	any, _ := codectypes.NewAnyWithValue(&inner_msg)

	msg := types.MsgSubmitBadSignatureEvidence{
		Subject:   any,
		Signature: "foo",
	}
	input.GravityKeeper.CheckBadSignatureEvidence(ctx, &msg)
}
