package witness

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/examples/democoin/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case oracle.Msg:
			return k.ork.Handle(func(ctx sdk.Context, p oracle.Payload) sdk.Error {
				switch p := p.(type) {
				case PayloadLock:
					return handlePayloadLock(ctx, k, p)
				default:
					errMsg := "Unrecognized witness oracle.Payload type: " + reflect.TypeOf(p).Name()
					return sdk.ErrUnknownRequest(errMsg)
				}
			}, ctx, msg, sdk.CodespaceUndefined) // TODO: set codespace
		default:
			errMsg := "Unrecognized witness Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlePayloadLock(ctx sdk.Context, k Keeper, p PayloadLock) sdk.Error {
	// TODO: issue tokens
	// blocked by sdk #1194
	/*
		k.ck.Issue(p.Coins, p.DestAddr)
	*/
	return nil
}
