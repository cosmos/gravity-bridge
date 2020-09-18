package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper Keeper
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation) error {
	switch att.ClaimType {
	case types.ClaimTypeEthereumBridgeDeposit:
		// todo: mint new vouchers
	case types.ClaimTypeEthereumBridgeWithdrawalBatch:
		// todo: mark batch as successful
	case types.ClaimTypeEthereumBridgeMultiSigUpdate:
		// todo: update nonce for "MultiSig Set"
	default:
		return sdkerrors.Wrapf(types.ErrDuplicate, "event type: %s", att.ClaimType)
	}
	return nil
}
