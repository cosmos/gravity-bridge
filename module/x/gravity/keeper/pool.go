package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// createSendToEthereum
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) createSendToEthereum(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	// If the coin is a gravity voucher, burn the coins. If not, check if there is a deployed ERC20 contract representing it.
	// If there is, lock the coins.

	isCosmosOriginated, tokenContract, err := k.DenomToERC20Lookup(ctx, totalAmount.Denom)
	if err != nil {
		return 0, err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, totalInVouchers); err != nil {
		return 0, err
	}

	// If it is no a cosmos-originated asset we burn
	if !isCosmosOriginated {
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, totalInVouchers); err != nil {
			panic(err)
		}
	}

	// get next tx id from keeper
	nextID := k.incrementLastSendToEthereumIDKey(ctx)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.

	// set the outgoing tx in the pool index
	k.setUnbatchedSendToEthereum(ctx, &types.SendToEthereum{
		Id:                nextID,
		Sender:            sender.String(),
		EthereumRecipient: counterpartReceiver,
		Erc20Token:        types.NewSDKIntERC20Token(amount.Amount, tokenContract),
		Erc20Fee:          types.NewSDKIntERC20Token(fee.Amount, tokenContract),
	})

	return nextID, nil
}

// cancelSendToEthereum
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) cancelSendToEthereum(ctx sdk.Context, id uint64, s string) error {
	sender, _ := sdk.AccAddressFromBech32(s)

	var send *types.SendToEthereum
	for _, ste := range k.getUnbatchedSendToEthereums(ctx) {
		if ste.Id == id {
			send = ste
		}
	}
	if send == nil {
		// NOTE: this case will also be hit if the transaction is in a batch
		return sdkerrors.Wrap(types.ErrInvalid, "id not found in send to ethereum pool")
	}

	if sender.String() != send.Sender {
		return fmt.Errorf("can't cancel a message you didn't send")
	}

	totalToRefund := send.Erc20Token.GravityCoin()
	totalToRefund.Amount = totalToRefund.Amount.Add(send.Erc20Fee.Amount)
	totalToRefundCoins := sdk.NewCoins(totalToRefund)
	isCosmosOriginated, _ := k.ERC20ToDenomLookup(ctx, send.Erc20Token.Contract)

	// If it is not cosmos-originated the coins are minted
	if !isCosmosOriginated {
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, totalToRefundCoins); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", totalToRefundCoins)
		}
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, totalToRefundCoins); err != nil {
		return sdkerrors.Wrap(err, "sending coins from module account")
	}

	k.deleteUnbatchedSendToEthereum(ctx, send.Id, send.Erc20Fee)
	return nil
}

func (k Keeper) setUnbatchedSendToEthereum(ctx sdk.Context, ste *types.SendToEthereum) {
	ctx.KVStore(k.storeKey).Set(types.MakeSendToEthereumKey(ste.Id, ste.Erc20Fee), k.cdc.MustMarshal(ste))
}

func (k Keeper) deleteUnbatchedSendToEthereum(ctx sdk.Context, id uint64, fee types.ERC20Token) {
	ctx.KVStore(k.storeKey).Delete(types.MakeSendToEthereumKey(id, fee))
}

func (k Keeper) iterateUnbatchedSendToEthereumsByContract(ctx sdk.Context, contract common.Address, cb func(*types.SendToEthereum) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.SendToEthereumKey}, contract.Bytes()...)).ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ste types.SendToEthereum
		k.cdc.MustUnmarshal(iter.Value(), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (k Keeper) IterateUnbatchedSendToEthereums(ctx sdk.Context, cb func(*types.SendToEthereum) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.SendToEthereumKey}).ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ste types.SendToEthereum
		k.cdc.MustUnmarshal(iter.Value(), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (k Keeper) getUnbatchedSendToEthereums(ctx sdk.Context) []*types.SendToEthereum {
	var out []*types.SendToEthereum
	k.IterateUnbatchedSendToEthereums(ctx, func(ste *types.SendToEthereum) bool {
		out = append(out, ste)
		return false
	})
	return out
}

func (k Keeper) incrementLastSendToEthereumIDKey(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LastSendToEthereumIDKey})
	var id uint64 = 0
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	newId := id + 1
	bz = sdk.Uint64ToBigEndian(newId)
	store.Set([]byte{types.LastSendToEthereumIDKey}, bz)
	return newId
}
