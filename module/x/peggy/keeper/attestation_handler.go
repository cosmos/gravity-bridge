package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type AttestationHandler struct {
	keeper Keeper
}

func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation) error {
	switch att.ClaimType {
	case types.ClaimTypeEthereumBridgeDeposit:
		// todo: mint new vouchers
	case types.ClaimTypeEthereumBridgeWithdrawalBatch:
		// todo: mark batch as successful
	case types.ClaimTypeEthereumBridgeMultiSigUpdate:
		// todo: update nonce for "MultiSig Set"
	default:
		return sdkerrors.Wrapf(types.ErrDuplicate, "event type: %X", att.ClaimType) // todo: claim type to string
	}
	return nil
}
