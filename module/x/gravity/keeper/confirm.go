package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// ProcessConfirm validates the confirmation then sets it in the store
func (k Keeper) ProcessConfirm(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	if err = k.ValidateConfirm(ctx, confirm, validatorAddr); err != nil {
		return err
	}
	return k.SetConfirm(ctx, confirm, validatorAddr)
}

// SetConfirm sets the confirmation in the store given it's type
func (k Keeper) SetConfirm(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	switch confirm := confirm.(type) {
	case *types.ConfirmBatch:
		if k.HasBatchConfirm(ctx, confirm.Nonce, validatorAddr, common.HexToAddress(confirm.TokenContract)) {
			return sdkerrors.Wrap(types.ErrSignatureDuplicate, "duplicate signature")
		}
		k.SetConfirmBatch(ctx, confirm, validatorAddr)
	case *types.ConfirmLogicCall:
		if k.HasConfirmLogicCall(ctx, confirm.InvalidationId, confirm.InvalidationNonce, validatorAddr) {
			return sdkerrors.Wrap(types.ErrSignatureDuplicate, "duplicate signature")
		}
		k.SetConfirmLogicCall(ctx, confirm, validatorAddr)
	case *types.ConfirmSignerSet:
		if k.HasSignerSetConfirm(ctx, confirm.Nonce, validatorAddr) {
			return sdkerrors.Wrap(types.ErrSignatureDuplicate, "duplicate signature")
		}
		k.SetConfirmSignerSet(ctx, confirm, validatorAddr)
	default:
		return sdkerrors.Wrapf(types.ErrConfirmUnsupported, "confirm type %s: %T", confirm.GetType(), confirm)
	}
	return nil
}

// ValidateConfirm validates the confirmation and checks it's ETH signature for validity
func (k Keeper) ValidateConfirm(ctx sdk.Context, confirm types.Confirm, validatorAddr sdk.ValAddress) (err error) {
	ethAddress := k.GetEthAddress(ctx, validatorAddr)
	if (ethAddress == common.Address{}) {
		return sdkerrors.Wrap(types.ErrValidatorEthAddressNotFound, validatorAddr.String())
	}

	if err = confirm.Validate(); err != nil {
		return err
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

// CheckpointBatchTx returns the abi encoded call and an error
func (k Keeper) CheckpointBatchTx(ctx sdk.Context, confirm *types.ConfirmBatch, bridgeID tmbytes.HexBytes) ([]byte, error) {
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

// CheckpointLogicCallTx returns the abi encoded call and an error
func (k Keeper) CheckpointLogicCallTx(ctx sdk.Context, confirm *types.ConfirmLogicCall, bridgeID tmbytes.HexBytes) ([]byte, error) {
	logicCallTx, found := k.GetLogicCallTx(ctx, confirm.InvalidationId, confirm.InvalidationNonce)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTxNotFound, "logic call tx")
	}
	return logicCallTx.GetCheckpoint(bridgeID, confirm.InvalidationId, confirm.InvalidationNonce)
}

// CheckpointEthSignerSet returns the abi encoded call and an error
func (k Keeper) CheckpointEthSignerSet(ctx sdk.Context, confirm *types.ConfirmSignerSet, bridgeID tmbytes.HexBytes) ([]byte, error) {
	signerSet, found := k.GetEthSignerSet(ctx, confirm.Nonce)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrSignerSetNotFound, "nonce %d", confirm.Nonce)
	}
	return signerSet.GetCheckpoint(bridgeID)
}

//////////////////
// ConfirmBatch //
//////////////////

// SetConfirmBatch sets a confirmation signature for a given validator, nonce and address into the store
func (k Keeper) SetConfirmBatch(ctx sdk.Context, confirm *types.ConfirmBatch, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetBatchConfirmKey(confirm.Nonce, common.HexToAddress(confirm.TokenContract), valAddr), confirm.Signature)
}

// GetBatchConfirm returns the signature for a given validator, nonce and address from the store
func (k Keeper) GetConfirmBatch(ctx sdk.Context, nonce uint64, valAddr sdk.ValAddress, contractAddr common.Address) hexutil.Bytes {
	return ctx.KVStore(k.storeKey).Get(types.GetBatchConfirmKey(nonce, contractAddr, valAddr))
}

// HasBatchConfirm returns if a validator has confirmed a given batch identified by nonce and contract address
func (k Keeper) HasBatchConfirm(ctx sdk.Context, nonce uint64, valAddr sdk.ValAddress, contractAddr common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetBatchConfirmKey(nonce, contractAddr, valAddr))
}

// IterateBatchConfirms iterates over all batch confirmations
func (k Keeper) IterateBatchConfirms(ctx sdk.Context, nonce uint64, contractAddr common.Address, cb func(val sdk.ValAddress, sig hexutil.Bytes) bool) {
	iter := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(append(types.BatchConfirmKey, contractAddr.Bytes()...), sdk.Uint64ToBigEndian(nonce)...),
	).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if cb(sdk.ValAddress(iter.Key()), hexutil.Bytes(iter.Value())) {
			break
		}
	}
}

// GetBatchConfirms returns all the confirmations in map[valaddress]signature format
func (k Keeper) GetBatchConfirms(ctx sdk.Context, nonce uint64, contractAddr common.Address) (out map[string][]byte) {
	k.IterateBatchConfirms(ctx, nonce, contractAddr, func(val sdk.ValAddress, sig hexutil.Bytes) bool {
		out[val.String()] = sig
		return false
	})
	return
}

//////////////////////
// ConfirmLogicCall //
//////////////////////

// SetConfirmLogicCall sets a confirmation signature for a given validator into the store
func (k Keeper) SetConfirmLogicCall(ctx sdk.Context, confirm *types.ConfirmLogicCall, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetLogCallConfirmKey(confirm.InvalidationId, confirm.InvalidationNonce, valAddr), confirm.Signature)
}

// GetLogicCallConfirm sets the confirmation signature for a given validator into the store
func (k Keeper) GetConfirmLogicCall(ctx sdk.Context, invalid []byte, invalnonce uint64, valAddr sdk.ValAddress) hexutil.Bytes {
	return ctx.KVStore(k.storeKey).Get(types.GetLogCallConfirmKey(invalid, invalnonce, valAddr))
}

// HasConfirmLogicCall returns true if the key exists in the store
func (k Keeper) HasConfirmLogicCall(ctx sdk.Context, invalid []byte, invalnonce uint64, valAddr sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetLogCallConfirmKey(invalid, invalnonce, valAddr))
}

// IterateConfirmLogicCalls iterates over all the logic call confirmations in the store
func (k Keeper) IterateConfirmLogicCalls(ctx sdk.Context, invalid []byte, invalnonce uint64, cb func(valAddr sdk.ValAddress, sig hexutil.Bytes) bool) {
	iter := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(append(types.KeyConfirmLogicCall, invalid...), sdk.Uint64ToBigEndian(invalnonce)...),
	).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if cb(sdk.ValAddress(iter.Key()), hexutil.Bytes(iter.Value())) {
			break
		}
	}
}

// GetLogicCallConfirms returns all the confirmations in map[valaddress]signature format
func (k Keeper) GetLogicCallConfirms(ctx sdk.Context, invalid []byte, invalnonce uint64) (out map[string][]byte) {
	k.IterateConfirmLogicCalls(ctx, invalid, invalnonce, func(val sdk.ValAddress, sig hexutil.Bytes) bool {
		out[val.String()] = sig
		return false
	})
	return
}

//////////////////////
// ConfirmSignerSet //
//////////////////////

// SetConfirmSignerSet sets a confirmation signature for a given validator and nonce into the store
func (k Keeper) SetConfirmSignerSet(ctx sdk.Context, confirm *types.ConfirmSignerSet, valAddr sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetSignerSetConfirmKey(confirm.Nonce, valAddr), confirm.Signature)
}

// GetSignerSetConfirm gets a confirmation signature for a given validator and nonce to the store
func (k Keeper) GetSignerSetConfirm(ctx sdk.Context, nonce uint64, valAddr sdk.ValAddress) hexutil.Bytes {
	return ctx.KVStore(k.storeKey).Get(types.GetSignerSetConfirmKey(nonce, valAddr))
}

// HasSignerSetConfirm returns if a validator has confirmed a given signerset identified by a nonce
func (k Keeper) HasSignerSetConfirm(ctx sdk.Context, nonce uint64, valAddr sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetSignerSetConfirmKey(nonce, valAddr))
}

// IterateSignerSetConfirms iterates over all the signer set confirmations in the store
func (k Keeper) IterateSignerSetConfirms(ctx sdk.Context, nonce uint64, cb func(valAddr sdk.ValAddress, sig hexutil.Bytes) bool) {
	iter := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(types.SignersetConfirmKey, sdk.Uint64ToBigEndian(nonce)...),
	).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if cb(sdk.ValAddress(iter.Key()), hexutil.Bytes(iter.Value())) {
			break
		}
	}
}

// GetSignerSetConfirms returns all the confirmations in map[valaddress]signature format
func (k Keeper) GetSignerSetConfirms(ctx sdk.Context, nonce uint64) (out map[string][]byte) {
	k.IterateSignerSetConfirms(ctx, nonce, func(val sdk.ValAddress, sig hexutil.Bytes) bool {
		out[val.String()] = sig
		return false
	})
	return
}
