package etend 

import (

)

const (
    // https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
    maxInt64 = 9223372036854775807
    maxBigInt64 = big.NewInt(maxInt64)
    maxUint64 = 18446744073709551615 
    maxBigUint64 = big.NewInt(maxUint64)

    maxCoin := coin.Coin {
        Denom: token,
        Amount: maxInt64,
    }
)

type Handler struct {
    stack.PassInitValidate
}

func (Handler) AssertDispatcher() {}

var _ stack.Dispatchable = Handler{}

func (Handler) Name() string {
    return NameETEnd
}

func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB, tx.sdk.Tx, next sdk.Checker) (res sdk.CheckResult, err error) {
    
}

func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    err := tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch t := tx.Unwrap().(type) {
    case DepositTx:
        return res, h.depositTx(ctx, store, t, next)
    case WithdrawTx:
        return res, h.withdrawTx(ctx, store, t, next)
    case TransferTx:
        return res, h.transferTx(ctx, store, t, next)
    }

    return res, errors.ErrUnknownTxType(tx.Unwrap())
}

func toCoins(token string, value []byte) (res coin.Coins, err error) {
    var n, m, v *big.Int
    v.SetBytes(value)
    n.DivMod(v, maxBigInt64, m)

    // It wont be happen... just for case
    if n.Cmp(maxBigUint64) == 1 {
        return res, ErrExceedMax()
    }

    res = make(coin.Coins, n.Uint64())

    for i, _ := range res {
        res[i] = maxCoin
    }

    if m.Uint64() != 0 {
        res = append(res, coin.Coin {
            Denom: token,
            Amount: int64(m.Uint64()),
        })
    }

    return res, nil
}

func (h Handler) changeBalance(ctx sdk.Context, store state.SimpleDB, token string, value []byte, to []byte, next sdk.Deliver) error {
    debtior := sdk.Actor {
        ChainID: ctx.ChainID(),
        App: ,
        Address: to,
    }
    coins, err := toCoins(token, value)
    if err != nil {
        return err
    }

    credit := coin.NewCreditTx(debitor, coins)
    _, err = next.DeliverTx(ctx, store, credit.Wrap())
    return err
}

func (h Handler) ibcDeliver(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Deliver) error {
    packet := ibc.CreatePacketTx {
        DestChain: ETGateChain(),
        Permissions: ,
        Tx: tx,
    }
    ibcCtx := ctx.WithPermissions(ibc.AllowIBC(NameETEnd)) // NameETGate?
    _, err := next.DeliverTx(ibcCtx, store, packet.Wrap())
    return err
}

func (h Handler) depositTx(ctx sdk.Context, store state.SimpleDB, tx DepositTx, next sdk.Deliver) error {
    return h.changeBalance(ctx, store, tx.Token, tx.Value, tx.To, next)
}

func (h Handler) withdrawTx(ctx sdk.Context, store state.SimpleDB, tx WithdrawTx, next sdk.Deliver) error {
    err := h.changeBalance(ctx, store, tx.Token, tx.Value.Neg(), tx.To, next)
    if err != nil {
        return err
    }

    withdraw := etgate.WithdrawTx {
        
    }   
    return ibcDeliver(ctx, store, withdraw, next)    
}

func (h Handler) transferTx(ctx sdk.Context, store state.SimpleDB, tx TransferTx, next sdk.Deliver) error {
    err := h.changeBalance(ctx, store, tx.Token, tx.Value.Neg(), tx.To, next)
    if err != nil {
        return err
    }

    transfer := etgate.TransferTx {
        
    }
    return ibcDeliver(ctx, store, transfer, next)
}

func (h Handler) InitState(l log.Logger, store state.SimpleDB, module, key, value string) (log string, err error) {
    if module != NameETEnd {
        return "", errors.ErrUnknownModule(module)
    }

    switch key {
    case "etgate":
        return setETGate()
    }

    return "", 
}
