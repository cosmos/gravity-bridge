package witness

import (
    "reflect"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(wmap WitnessMsgMapper, ck bank.CoinKeeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case LockMsg:
            return handleLockMsg(ctx, wmap, ck, msg)
        default:
            errMsg := "Unrecognized withdraw Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleLockMsg(ctx sdk.Context, wtx WitnessMsgMapper, ck bank.CoinKeeper, msg LockMsg) sdk.Result {
    data := wtx.GetWitnessData(ctx, msg)
    if data.credited {
        return ErrAlreadyCredited()
    }
    for _, w := range data.Witnesses {
        if w == msg.Signer {
            return ErrWitnessReplay()
        }
    }
    data.Witnesses = append(data.Witnesses, msg.Signer)
    if len(data.Witnesses) >= 67 {
        coin := sdk.Coin {
            Denom: string(msg.Token),
            Amount: msg.Amount,
        }
        ck.AddCoins(ctx, msg.Destination, coin)
        data.credited = true
    }
    wtx.SetWitnessData(ctx, msg, data)
    return sdk.Result{}
}
