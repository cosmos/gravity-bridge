package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (k Keeper) ConfirmEvent(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	if err = k.ValidateConfirmSignature(ctx, confirm, validatorAddr); err != nil {
		return err
	}
	return k.SetConfirm(ctx, confirm, validatorAddr)
}

func (k Keeper) SetConfirm(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	switch confirm := confirm.(type) {
	case *types.ConfirmBatch:
		k.SetConfirmBatch(ctx, confirm, validatorAddr)
	case *types.ConfirmLogicCall:
		k.SetConfirmLogicCall(ctx, confirm, validatorAddr)
	case *types.ConfirmSignerSet:
		k.SetConfirmSignerSet(ctx, confirm, validatorAddr)
	default:
		return sdkerrors.Wrapf(types.ErrConfirmUnsupported, "confirm type %s: %T", confirm.GetType(), confirm)
	}
	return nil
}

func (k Keeper) ValidateConfirmSignature(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	ethAddress := k.GetEthAddress(ctx, validatorAddr)
	if (ethAddress == common.Address{}) {
		return sdkerrors.Wrap(types.ErrValidatorEthAddressNotFound, validatorAddr.String())
	}

	bridgeID := k.GetBridgeID(ctx)
	var checkpoint []byte
	switch confirm := confirm.(type) {
	case *types.ConfirmBatch:
		checkpoint, err = k.CheckpointBatchTx(ctx, confirm, bridgeID)
	case *types.ConfirmLogicCall:
		checkpoint, err = k.CheckpointLogicCallTx(ctx, confirm, bridgeID)
	case *types.ConfirmSignerSet:
		checkpoint, err = k.CheckpointEthSignerSet(ctx, confirm, bridgeID)
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
	return nil
}

func (k Keeper) CheckpointBatchTx(
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

func (k Keeper) CheckpointLogicCallTx(
	ctx sdk.Context, confirm *types.ConfirmLogicCall, bridgeID tmbytes.HexBytes,
) ([]byte, error) {
	logicCallTx, found := k.GetLogicCallTx(ctx, confirm.InvalidationID, confirm.InvalidationNonce)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTxNotFound, "logic call tx")
	}

	return logicCallTx.GetCheckpoint(bridgeID, confirm.InvalidationID, confirm.InvalidationNonce)
}

func (k Keeper) CheckpointEthSignerSet(
	ctx sdk.Context, confirm *types.ConfirmSignerSet, bridgeID tmbytes.HexBytes,
) ([]byte, error) {

	signerSet, found := k.GetEthSignerSet(ctx, confirm.Nonce)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrSignerSetNotFound, "nonce %d", confirm.Nonce)
	}

	return signerSet.GetCheckpoint(bridgeID)
}

// SetConfirmBatch sets a confirmation signature for a given validator into the store
func (k Keeper) SetConfirmBatch(ctx sdk.Context, confirm *types.ConfirmBatch, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetBatchConfirmKey(confirm, valAddr), confirm.Signature)
}

// GetBatchConfirm returns the signature for a given validator from the store
func (k Keeper) GetConfirmBatch(ctx sdk.Context, valAddr sdk.ValAddress, contractAddr common.Address) hexutil.Bytes {
	return ctx.KVStore(k.storeKey).Set(types.GetBatchConfirmKey(confirm))
}

// SetConfirmLogicCall sets a confirmation signature for a given validator into the store
func (k Keeper) SetConfirmLogicCall(ctx sdk.Context, confirm *types.ConfirmLogicCall, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetLogCallConfirmKey(confirm, valAddr), confirm.Signature)
}

// func (k Keeper) GetConfirmBatch(ctx sdk.Context, valAddr sdk.ValAddress) hexutil.Bytes

// SetConfirmSignerSet sets a confirmation signature for a given validator into the store
func (k Keeper) SetConfirmSignerSet(ctx sdk.Context, confirm *types.ConfirmSignerSet, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetSignerSetConfirmKey(confirm, valAddr), confirm.Signature)
}

// func (k Keeper) GetConfirmBatch(ctx sdk.Context, valAddr sdk.ValAddress) hexutil.Bytes
