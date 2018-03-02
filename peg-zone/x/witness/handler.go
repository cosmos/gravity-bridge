package witness

import (
    "bytes"
    "reflect"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/bank"
)

const (
    totalValidators = 2 // temp
)

func NewHandler(wmap WitnessMsgMapper, ck bank.CoinKeeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case WitnessMsg:
            return handleWitnessMsg(ctx, wmap, ck, msg)
        default:
            errMsg := "Unrecognized Witness Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleWitnessMsg(ctx sdk.Context, wmsg WitnessMsgMapper, ck bank.CoinKeeper, msg WitnessMsg) sdk.Result {
    info := msg.Info
    data := wmsg.GetWitnessData(ctx, info)
    if data.Credited {
        return ErrAlreadyCredited().Result()
    }
    /*
    if !isValidator(msg.Signer) {
        return ErrSignerIsNotAValidator().Result()
    }
    */
    for _, w := range data.Witnesses {
        if bytes.Equal(w, msg.Signer) {
            return ErrWitnessReplay().Result()
        }
    }
    data.Witnesses = append(data.Witnesses, msg.Signer)
    if len(data.Witnesses) * 3 >= totalValidators * 2 {
        switch info := info.(type) {
        case LockInfo:
            coin := sdk.Coin {
                Denom: string(info.Token),
                Amount: info.Amount,
            }
            ck.AddCoins(ctx, info.Destination, []sdk.Coin{coin})
            data.Credited = true
        }
    }
    wmsg.SetWitnessData(ctx, info, data)
    return sdk.Result{}
}
