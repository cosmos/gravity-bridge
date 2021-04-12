package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO: rename functions to Send / Receive
// TODO: test with IBC vouchers

// AddToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
//
//
// CONTRACT: amount and fee must be valid Ethereum ERC20 token or a Cosmos coin
// (i.e with or without the gravity prefix)
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, ethereumReceiver common.Address, amount, fee sdk.Coin) (uint64, error) {
	if amount.Denom != fee.Denom {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "coin denom doesn't match with fee denom (%s ≠ %s)", amount.Denom, fee.Denom)
	}

	// Add the fees to the transfer coins in order to escrow them on the ModuleAccount
	coinsToEscrow := sdk.NewCoins(amount.Add(fee))

	// If the coin is a gravity voucher, burn the coins. If not, check if there is a deployed ERC20 contract representing it.
	// If there is, lock the coins.

	var (
		tokenContract    common.Address
		tokenContractHex string
		found            bool
	)

	if types.IsEthereumERC20Token(amount.Denom) {
		tokenContractHex = types.GravityDenomToERC20Contract(amount.Denom)
		tokenContract = common.HexToAddress(tokenContractHex)

		// If it is an ethereum-originated asset we burn it
		// send coins to module in prep for burn
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToEscrow); err != nil {
			return 0, err
		}

		// burn vouchers to send them back to ETH
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToEscrow); err != nil {
			panic(fmt.Errorf("couldn't burn ERC20 vouchers %s: %w", tokenContractHex, err))
		}
	} else {
		// coin is a native Cosmos coin, fetch the contract if exists
		tokenContract, found = k.GetERC20ContractFromCoinDenom(ctx, amount.Denom)
		if !found {
			// TODO: what if there is no corresponding contract? will it be "generated" on ethereum
			// upon receiving?
			// FIXME: this will fail if the cosmos tokens are relayed for the first time and they don't have a counterparty contract
			// Fede: Also there's the question of how do we handle IBC denominations from a security perspective. Do we assign them the same
			// contract? My guess is that each new contract assigned to a cosmos coin should be approved by governance
			return 0, sdkerrors.Wrapf(types.ErrContractNotFound, "denom %s", amount.Denom)
		}

		tokenContractHex = tokenContract.String()

		// lock coins in module
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToEscrow); err != nil {
			return 0, err
		}
	}

	// get next tx id from keeper
	// TODO: why is this an integer instead of a hash? The nonce should be uint64 and the id a hash
	txID := k.GetTxID(ctx)
	txID++

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.

	// we use only the token contract as the denom for outgoing transactions
	// to avoid unnecessary parsing on the orchestrator
	fee.Denom = tokenContractHex
	amount.Denom = tokenContractHex

	tx := types.TransferTx{
		Sender:            sender.String(),
		EthereumRecipient: ethereumReceiver.String(),
		Erc20Token:        amount,
		Erc20Fee:          fee,
	}

	// set the outgoing transfer tx in the pool
	k.SetTransferTx(ctx, txID, tx)

	// TODO: add the transfer tx to the unbatched transaction pool
	// k.appendToUnbatchedTxIndex(ctx, txID, tokenContractHex, fee)

	// set the incremented tx ID
	k.setTxID(ctx, txID)

	// TODO: fix events / add more attrs
	nonceStr := strconv.FormatUint(txID, 64)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawalReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			// TODO: why are the nonce and the txID the same?
			sdk.NewAttribute(types.AttributeKeyOutgoingTxID, nonceStr),
			sdk.NewAttribute(types.AttributeKeyNonce, nonceStr),
		),
	)

	return txID, nil
}

// RemoveFromOutgoingPoolAndRefund
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) RemoveFromOutgoingPoolAndRefund(ctx sdk.Context, txID uint64, sender sdk.AccAddress) error {
	// check that we actually have a tx with that id and what it's details are
	tx, found := k.GetTransferTx(ctx, txID)
	if !found {
		return sdkerrors.Wrapf(types.ErrOutgoingTxNotFound, "tx id %d", txID)
	}

	// TODO: check if the transaction is currently on a batch and remove it

	// poolTx := k.GetPoolTransactions(ctx)
	// for _, pTx := range poolTx {
	// 	if pTx.Id == txID {
	// 		found = true
	// 	}
	// }
	// if !found {
	// 	return sdkerrors.Wrapf(types.ErrInvalid, "Id %d is in a batch", txID)
	// }

	// k.removeFromUnbatchedTxIndex(ctx, txID, tx.Erc20Fee)

	// delete the tx from the transfer tx pool
	k.DeleteTransferTx(ctx, txID)

	// reissue the amount and the fee
	refund := tx.Erc20Token.Add(tx.Erc20Fee)

	// Coins native to cosmos are unlocked and transferred to the original sender
	if types.IsCosmosCoin(refund.Denom) {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(refund)); err != nil {
			return err
		}
	} else {
		// Ethereum ERC20 vouchers that were prev burned need to be re-minted and
		// transferred to the original sender
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(refund)); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", refund)
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(refund)); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute("denom", refund.Denom), // TODO: create attr
			sdk.NewAttribute(types.AttributeKeyOutgoingTxID, strconv.FormatUint(txID, 64)),
		),
	)

	return nil
}

// // appendToUnbatchedTXIndex add at the end when tx with same fee exists
// func (k Keeper) appendToUnbatchedTxIndex(ctx sdk.Context, txID uint64, tokenContract string, fee sdk.Coin) {
// 	store := ctx.KVStore(k.storeKey)
// 	idxKey := types.GetFeeSecondIndexKey(fee)
// 	var idSet types.IDSet
// 	if store.Has(idxKey) {
// 		bz := store.Get(idxKey)
// 		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
// 	}
// 	idSet.Ids = append(idSet.Ids, txID)
// 	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
// }

// // appendToUnbatchedTXIndex add at the top when tx with same fee exists
// func (k Keeper) prependToUnbatchedTxIndex(ctx sdk.Context, txID uint64, tokenContract string, fee sdk.Coin) {
// 	store := ctx.KVStore(k.storeKey)
// 	idxKey := types.GetFeeSecondIndexKey(fee)
// 	var idSet types.IDSet
// 	if store.Has(idxKey) {
// 		bz := store.Get(idxKey)
// 		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
// 	}
// 	idSet.Ids = append([]uint64{txID}, idSet.Ids...)
// 	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
// }

// // removeFromUnbatchedTXIndex removes the tx from the index and makes it implicit no available anymore
// func (k Keeper) removeFromUnbatchedTxIndex(ctx sdk.Context, txID uint64, fee sdk.Coin) error {
// 	store := ctx.KVStore(k.storeKey)
// 	idxKey := types.GetFeeSecondIndexKey(fee)
// 	var idSet types.IDSet
// 	bz := store.Get(idxKey)
// 	if bz == nil {
// 		return sdkerrors.Wrap(types.ErrUnknown, "fee")
// 	}
// 	k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
// 	for i := range idSet.Ids {
// 		if idSet.Ids[i] == txID {
// 			idSet.Ids = append(idSet.Ids[0:i], idSet.Ids[i+1:]...)
// 			if len(idSet.Ids) != 0 {
// 				store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
// 			} else {
// 				store.Delete(idxKey)
// 			}
// 			return nil
// 		}
// 	}
// 	return sdkerrors.Wrap(types.ErrUnknown, "tx id")
// }

// // IterateOutgoingPoolByFee iterates over the outgoing pool which is sorted by fee
// func (k Keeper) IterateOutgoingPoolByFee(ctx sdk.Context, contract string, cb func(txID uint64, transferTx types.TransferTx) bool) {
// 	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SecondIndexOutgoingTxFeeKey)
// 	iter := prefixStore.ReverseIterator(prefixRange([]byte(contract)))
// 	defer iter.Close()
// 	for ; iter.Valid(); iter.Next() {
// 		var ids types.IDSet
// 		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)
// 		// cb returns true to stop early
// 		for _, id := range ids.Ids {
// 			tx, err := k.getPoolEntry(ctx, id)
// 			if err != nil {
// 				panic("Invalid id in tx index!")
// 			}
// 			if cb(id, tx) {
// 				return
// 			}
// 		}
// 	}
// }

// TODO: explain difference with batch ID?

func (k Keeper) GetTxID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastTxPoolID)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) setTxID(ctx sdk.Context, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastTxPoolID, sdk.Uint64ToBigEndian(txID))
}
