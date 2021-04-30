package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// TODO: rename functions to Send / Receive
// TODO: test with IBC vouchers

// AddTransferToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
//
//
// CONTRACT: amount and fee must be valid Ethereum ERC20 token or a Cosmos coin
// (i.e with or without the gravity prefix)
func (k Keeper) AddTransferToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, ethereumReceiver common.Address, amount, fee sdk.Coin) (tmbytes.HexBytes, error) {
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
			return nil, err
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
			return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "denom %s", amount.Denom)
		}

		tokenContractHex = tokenContract.String()

		// lock coins in module
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coinsToEscrow); err != nil {
			return nil, err
		}
	}

	// get next tx nonce from keeper
	nonce := k.GetTransferTxNonce(ctx)
	nonce++

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.

	// we use only the token contract as the denom for outgoing transactions
	// to avoid unnecessary parsing on the orchestrator
	fee.Denom = tokenContractHex
	amount.Denom = tokenContractHex

	tx := types.TransferTx{
		Nonce:             nonce,
		Sender:            sender.String(),
		EthereumRecipient: ethereumReceiver.String(),
		Erc20Token:        amount,
		Erc20Fee:          fee,
	}

	// set the outgoing transfer tx in the pool
	txID := k.SetTransferTx(ctx, tx)

	k.IndexTransferTxByFee(ctx, tokenContract, fee, txID)

	// set the incremented tx ID
	k.setTransferTxNonce(ctx, nonce)

	nonceStr := strconv.FormatUint(nonce, 64)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferPooled,
			sdk.NewAttribute(types.AttributeKeyTxID, txID.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyEthRecipient, ethereumReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, amount.Denom),             // cosmos or gravity denom
			sdk.NewAttribute(types.AttributeKeyTokenContract, tokenContractHex), // ERC20 contract address
			sdk.NewAttribute(types.AttributeKeyNonce, nonceStr),
		),
	)

	k.Logger(ctx).Info(
		"outgoing transfer",
		"id", txID.String(),
		"nonce", nonceStr,
	)

	return txID, nil
}

// RemoveFromOutgoingPoolAndRefund
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) RemoveFromOutgoingPoolAndRefund(ctx sdk.Context, txID tmbytes.HexBytes, sender sdk.AccAddress) error {
	// check that we actually have a tx with that id and what it's details are
	tx, found := k.GetTransferTx(ctx, txID)
	if !found {
		return sdkerrors.Wrapf(types.ErrTxNotFound, "transfer tx with id %s", txID)
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
			types.EventTypeTransferCanceled,
			sdk.NewAttribute(types.AttributeKeyTxID, txID.String()),
			sdk.NewAttribute(types.AttributeKeyRefundDenom, refund.Denom), // cosmos or gravity denom
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		),
	)

	return nil
}

// IndexTransferTxByFee
func (k Keeper) IndexTransferTxByFee(ctx sdk.Context, tokenContract common.Address, fee sdk.Coin, txID tmbytes.HexBytes) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransferTxFeeKey)
	key := types.GetTransferTxFeeKey(tokenContract.String(), fee.Amount.Uint64())

	var txIDs types.TransactionIDs

	bz := store.Get(key)
	if len(bz) != 0 {
		k.cdc.MustUnmarshalBinaryBare(bz, &txIDs)
	}

	txIDs.IDs = append(txIDs.IDs, txID)

	store.Set(key, k.cdc.MustMarshalBinaryBare(&txIDs))
}

// IterateTransferPoolByFee iterates in DESC order over the contracts and performs a callback with
// transactions that share the same contract and fee amount
func (k Keeper) IterateTransferPoolByFee(ctx sdk.Context, tokenContract common.Address, cb func(fee uint64, txIDs []tmbytes.HexBytes) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.TransferTxFeeKey, []byte(tokenContract.String())...))

	iter := store.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var txIDs types.TransactionIDs
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &txIDs)

		fee := uint64(0)
		if cb(fee, txIDs.IDs) {
			return
		}
	}
}

func (k Keeper) GetTransferTxNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastTransferTxID)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) setTransferTxNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastTransferTxID, sdk.Uint64ToBigEndian(nonce))
}
