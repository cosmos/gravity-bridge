package witness

import (
    "reflect"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(wmap WitnessMsgMapper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case LockMsg:
            return handleLockMsg(ctx, wmap, msg)
        default:
            errMsg := "Unrecognized withdraw Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleLockMsg(ctx sdk.Context, wtx WitnessMsgMapper, msg LockMsg) sdk.Result {
    data := wtx.GetWitnessData(ctx, msg)
    for _, w := range data.Witnesses {
        if w == msg.Signer {
            return ErrWitnessReplay()
        }
    }
    data.Witnesses = append(data.Witnesses, msg.Signer)
    wtx.SetWitnessData(ctx, msg, data)
    return sdk.Result{}
}
