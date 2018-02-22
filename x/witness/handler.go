package withdraw 

import (
    "reflect"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(with WitnessTxMapper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case WitnessTx:
            return handleWitnessTx(ctx, with, msg)
        default:
            errMsg := "Unrecognized withdraw Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleWitnessTx(ctx sdk.Context, with WithdrawTxMapper, msg sdk.Msg) sdk.Result {

}
