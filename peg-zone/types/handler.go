package types

func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper) {
	r.AddRoute(WitnessTx, DepositMsgHandler(accts))
	r.AddRoute(SendTx, SettleMsgHandler(accts))
	r.AddRoute(WithdrawTx, WithdrawMsgHandler(accts))
	r.AddRoute(SignTx, CreateOperatorMsgHandler(accts))
}

// Handle all peggy type messages.
func NewHandler(ck CoinKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SendMsg:
			return handleSendMsg(ctx, ck, msg)
		case IssueMsg:
			return handleIssueMsg(ctx, ck, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SendMsg.
func handleSendMsg(ctx sdk.Context, ck CoinKeeper, msg SendMsg) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	for _, in := range msg.Inputs {
		_, err := ck.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err.Result()
		}
	}

	for _, out := range msg.Outputs {
		_, err := ck.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{} // TODO
}

// Handle IssueMsg.
func handleIssueMsg(ctx sdk.Context, ck CoinKeeper, msg IssueMsg) sdk.Result {
	panic("not implemented yet")
}

type sendTxHandler struct {
}

type withdrawTxHandler struct {
}

type signTxHandler struct {
}
