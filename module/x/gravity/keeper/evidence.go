package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func (k Keeper) CheckBadSignatureEvidence(
	ctx sdk.Context,
	msg *types.MsgSubmitBadSignatureEvidence) error {
	var subject types.EthereumSigned

	k.cdc.UnpackAny(msg.Subject, &subject)

	switch subject := subject.(type) {
	case *types.OutgoingTxBatch:
	case *types.Valset:
	case *types.OutgoingLogicCall:
		// Get checkpoint of the supposed bad signature (fake valset, batch, or logic call submitted to eth)
		gravityID := k.GetGravityID(ctx)
		checkpoint := subject.GetCheckpoint(gravityID)

		// Try to find the checkpoint in the archives. If it exists, we don't slash because
		// this is not a bad signature
		if k.GetPastEthSignatureCheckpoint(ctx, checkpoint) {
			return sdkerrors.Wrap(types.ErrInvalid, "Checkpoint exists, cannot slash")
		}

		// Decode Eth signature to bytes
		sigBytes, err := hex.DecodeString(msg.Signature)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
		}

		// Get eth address of the offending validator using the checkpoint and the signature
		ethAddress, err := types.EthAddressFromSignature(checkpoint, sigBytes)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s", ethAddress, gravityID, hex.EncodeToString(checkpoint), msg.Signature))
		}

		// Find the offending validator by eth address
		val, found := k.GetValidatorByEthAddress(ctx, ethAddress)
		if !found {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("Did not find validator for eth address %s", ethAddress))
		}

		// Slash the offending validator
		cons, _ := val.GetConsAddr()
		params := k.GetParams(ctx)
		k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBadEthSignature)
		if !val.IsJailed() {
			k.StakingKeeper.Jail(ctx, cons)
		}

	default:
		panic("foo") // Replace this with a returned error
	}

	return nil
}

// SetPastEthSignatureCheckpoint puts the checkpoint of a valseet, batch, or logic call into a set
// in order to prove later that it existed at one point.
func (k Keeper) SetPastEthSignatureCheckpoint(ctx sdk.Context, checkpoint []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPastEthSignatureCheckpointKey(checkpoint), []byte{0x1})
}

// GetPastEthSignatureCheckpoint tells you whether a given checkpoint has ever existed
func (k Keeper) GetPastEthSignatureCheckpoint(ctx sdk.Context, checkpoint []byte) (found bool) {
	store := ctx.KVStore(k.storeKey)
	if bytes.Equal(store.Get(types.GetPastEthSignatureCheckpointKey(checkpoint)), []byte{0x1}) {
		return true
	} else {
		return false
	}
}
