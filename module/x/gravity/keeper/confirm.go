package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (k Keeper) ConfirmEvent(ctx sdk.Context, confirm types.Confirm, orchestratorAddr sdk.AccAddress, ethAddress common.Address) error {
	bridgeID := k.GetBridgeID(ctx)

	var (
		checkpoint []byte
		err        error
	)

	switch confirm := confirm.(type) {
	case *types.ConfirmBatch:
		checkpoint, err = k.ConfirmBatchTx(ctx, confirm, bridgeID)
	case *types.ConfirmLogicCall:
		checkpoint, err = k.ConfirmLogicCallTx(ctx, confirm, bridgeID)
	case *types.ConfirmSignerSet:
		checkpoint, err = k.ConfirmEthSignerSet(ctx, confirm, bridgeID)
	default:
		return sdkerrors.Wrapf(types.ErrConfirmUnsupported, "confirm type %s: %T", confirm.GetType(), confirm)
	}

	if err != nil {
		return err
	}

	signatureAddr, err := types.EcRecover(checkpoint, confirm.GetSignature())
	if err != nil {
		return sdkerrors.Wrapf(types.ErrSignatureInvalid, "signature verification failed: %w", err)
	}

	if signatureAddr != ethAddress {
		return fmt.Errorf(
			"signature address doesn't match the provided ethereum address (%s ≠ %s)",
			signatureAddr, ethAddress,
		)
	}

	// TODO: set confirm to store?
	// TODO: emit event
	return nil
}

func (k Keeper) ConfirmBatchTx(
	ctx sdk.Context, confirm *types.ConfirmBatch, bridgeID tmbytes.HexBytes,
) ([]byte, error) {
	// TODO:
	batchTx, found := k.GetBatchTx(ctx, common.Address{}, nil)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTxNotFound, "batch tx")
	}

	transfers := make([]types.TransferTx, len(batchTx.Transactions))

	for i, txID := range batchTx.Transactions {
		transfers[i], found = k.GetTransferTx(ctx, txID)
		if !found {
			return nil, sdkerrors.Wrapf(types.ErrTxNotFound, "transfer tx with ID %s is on batch but not on store", txID)
		}
	}

	return batchTx.GetCheckpoint(bridgeID, transfers)
}

func (k Keeper) ConfirmLogicCallTx(
	ctx sdk.Context, confirm *types.ConfirmLogicCall, bridgeID tmbytes.HexBytes,
) ([]byte, error) {
	logicCallTx, found := k.GetLogicCallTx(ctx, confirm.InvalidationID, confirm.InvalidationNonce)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTxNotFound, "logic call tx")
	}

	return logicCallTx.GetCheckpoint(bridgeID, confirm.InvalidationID, confirm.InvalidationNonce)
}

func (k Keeper) ConfirmEthSignerSet(
	ctx sdk.Context, confirm *types.ConfirmSignerSet, bridgeID tmbytes.HexBytes,
) ([]byte, error) {

	signerSet, found := k.GetEthSignerSet(ctx, confirm.Nonce)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrSignerSetNotFound, "nonce %d", confirm.Nonce)
	}

	return signerSet.GetCheckpoint(bridgeID)
}
