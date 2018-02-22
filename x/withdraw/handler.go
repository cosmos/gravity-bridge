package withdraw 

import (
    "reflect"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(with WithdrawTxMapper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case WithdrawTx:
            return handleWithdrawTx(ctx, with, msg)
        default:
            errMsg := "Unrecognized withdraw Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleWithdrawTx(ctx sdk.Context, with WithdrawTxMapper, msg sdk.Msg) sdk.Result {

}
