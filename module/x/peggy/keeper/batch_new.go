package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// pickUnbatchedTx find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTx(ctx sdk.Context, maxElements int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IteratePoolTxByFee(ctx, func(txID uint64, tx *types.OutgoingTx) bool {
		erc20Amount, err := types.ERC20FromPeggyCoin(tx.Amount)
		txOut := &types.OutgoingTransferTx{
			Id:          txID,
			Sender:      tx.Sender,
			DestAddress: tx.DestAddr,
			Erc20Token:  erc20Amount,
		}
		selectedTx = append(selectedTx, txOut)
		err = k.removeFromUnbatchedTXIndex(ctx, tx.BridgeFee, txID)
		return err != nil || len(selectedTx) == maxElements
	})
	return selectedTx, err
}
