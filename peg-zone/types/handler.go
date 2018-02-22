package types

func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper, wits WitnessTxMapper, witd WithdrawTxMapper) {
	r.AddRoute(WitnessTx, DepositMsgHandler(accts))
	r.AddRoute(SendTx, SettleMsgHandler(accts))
	r.AddRoute(WithdrawTx, WithdrawMsgHandler(accts))
	r.AddRoute(SignTx, CreateOperatorMsgHandler(accts))
}

// Handle all peggy type messages.
func NewHandler(wtx WitnessTxMapper, val ValidatorMapper) sdk.Handler {
	return func(ctx sdk.Context, ck sdk.CoinKeeper, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case WitnessTx:
			return handleWitnessTx(ctx, wtx, val, msg)
		case SendTx:
			return handleSendTx(ctx, ck, msg)
		case WithdrawTx:
			return handleWithdrawTx(ctx, msg)
		case SignTx:
			return handleSignTx(ctx, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SendMsg.
func handleSendTx(ctx sdk.Context, ck CoinKeeper, msg SendMsg) sdk.Result {
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
func handleWitnessTx(ctx sdk.Context, wtx WitnessTxMapper, val ValidatorMapper, msg WitnessTx) sdk.Result {
    data := wtx.GetWitnessData(ctx, msg)
    for _, d := range data {
    }
}

// Handle IssueMsg.
func handleWithdrawTx(ctx sdk.Context, ck CoinKeeper, msg IssueMsg) sdk.Result {
	panic("not implemented yet")
}

// Handle IssueMsg.
func handleSignTx(ctx sdk.Context, ck CoinKeeper, msg IssueMsg) sdk.Result {
	panic("not implemented yet")
}

type sendTxHandler struct {
}

type withdrawTxHandler struct {
}

type signTxHandler struct {
}
