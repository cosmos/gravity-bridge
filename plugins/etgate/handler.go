package etgate

import (
    sdk "github.com/cosmos/cosmos-sdk"
    "github.com/cosmos/cosmos-sdk/errors"
    "github.com/cosmos/cosmos-sdk/modules/auth"
    "github.com/cosmos/cosmos-sdk/modules/base"
    "github.com/cosmos/cosmos-sdk/modules/coin"
    "github.com/cosmos/cosmos-sdk/modules/fee"
    "github.com/cosmos/cosmos-sdk/modules/ibc"
    "github.com/cosmos/cosmos-sdk/modules/nonce"
    "github.com/cosmos/cosmos-sdk/modules/roles"
    "github.com/cosmos/cosmos-sdk/stack"
    "github.com/cosmos/cosmos-sdk/state"

)

const (
    NameETGate = "etgate"

    confirmation = 4 // for dev, 12 in production
)

var (
    depositabi abi.ABI
)

func init() {
    deposit, err := abi.JSON(strings.NewReader(contracts.ETGateABI))
    if err != nil {
        panic(err)
    }
    depositabi = deposit
}

// TODO: handle them manually, we need to know valset
type Handler struct {
    stack.PassInitValidate
}

func (Handler) AssertDispatcher() {}

var _ stack.Dispatchable = Handler{}

func (Handler) Name() string {
    return NameETGate
}

func (Handler) AssertDispather() {}

// https://github.com/cosmos/cosmos-sdk/blob/develop/modules/ibc/handler.go

func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Checker) (res sdk.DeliverResult, err error) {
    err := tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch tx.Unwrap().(type) {
    case UpdateTx:
        return res, nil
    // case ValChangeTx:
    case DepositTx:
        return res, nil
    case WithdrawTx:
        return res, nil
    case TransferTx:
        return res, nil
    }

    return res, errors.ErrUnknownTxType(tx.Unwrap())
}

func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    err := tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch t := tx.Unwrap().(type) {
    case UpdateTx:
        return h.updateTx(ctx, store, t)
    // case ValChangeTx:
    case DepositTx:
        return h.depositTx(ctx, store, t, next)
    case WithdrawTx:
        return h.withdrawTx(ctx, store, t)
    case TransferTx:
        return h.transferTx(ctx, store, t, next)
    }


    return res, errors.ErrUnknownTxType(tx.Unwrap())
}

func validateHeaders(headers []Header) error {
    for i, h := range headers[1:] {
        if h.ParentHash != headers[i-1].Hash {
            return errors.New("Non-continuous header list")
        }
    }
    return nil
}

func decodeHeader(headerb []byte) (Header, error) {
    var header eth.Header
    if err := rlp.DecodeBytes(headerb, &header); err != nil {
        return Header{}, err
    }
    return Header {
        ParentHash:  header.ParentHash,
        Hash:        header.Hash(),
        ReceiptHash: header.ReceiptHash,
        Number:      header.Number.Uint64(),
        Time:        header.Time.Uint64()
    }, nil
}

func decodeHeaders(headersb [][]byte) (headers []Header, error) {
    headers = make([]Header, len(headersb))
    for i, headerb := range headersb {
        header, err := decodeHeader(headerb)
        if err != nil {
            return err
        }
        headers[i] = header
    }
    return headers, nil
}

func updateHeaders(headers []Header, set ChainSet) error {
    // push all headers to buffer
    for _, h := range headers {
        if err := set.ToBuffer(h); err != nil {
            return err
        }        
    }

    genesis, err := set.GetGenesis()
    if err != nil {
        return err
    }

    headAnc, err := set.GetAncestor(headers[0], genesis)
    if err != nil {
        return err
    }

    lf, err := set.LastFinalized()
    if err != nil {
        return err
    }

    // check headAnc is the direct child of lf
    if lf != headAnc.Number-1 {
        return Err
    }
    
    lastAnc, err := set.GetAncestor(headers[len(headers)-1], genesis)
    if err != nil {
        return err
    }

    // finalize lastAnc to headAnc
    for {
        set.Finalize(lastAnc)
        if lastAnc.Hash == headAnc.Hash {
            break
        }
        lastAnc, err = set.Parent(lastAnc)
        if err != nil {
            return err
        }
    }

    return nil
}



func (h Handler) updateTx(ctx sdk.Context, store state.SimpleDB, t UpdateTx) (res sdk.DeliverResult, err error) {
    // TODO: check sender is a validator

    set := NewChainSet(store)

    if !set.IsInitialized() {
        return res, ErrNotInitialized()
    }

    headers, err := decodeHeaders(tx.Headers)
    if err != nil {
        return res, err
    }

    err = validateHeaders(headers)
    if err != nil {
        return res, err
    }

    err = updateHeaders(headers, NewChainSet(store))
    if err != nil {
        return res, err
    }

    return
}

func (h Handler) depositTx(ctx sdk.Context, store state.SimpleDB, t DepositTx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    log, err := tx.Proof.Log()
    if err != nil {
        return res, ErrInvalidLogProof(log, err)
    }

    set := NewChainSet(store)

    header, err := set.GetHeader(tx.Proof.Number)
    if err != nil {
        return res, err
    }

    if !tx.Proof.IsValid(header.ReceiptHash) {
        return res, ErrInvalidLogProof(log, header)
    }

    deposit := new(etend.DepositTx)
    if err != depositabi.Unpack(deposit, "Deposit", log); err != nil {
        return res, ErrLogUnpackingError(log, err)
    }

    if c.DepositExists(deposit) {
        return res, ErrDepositExists(deposit)
    }

    if deposit.DestChain == ctx.ChainID() {
        return res, ErrInvalidDestChain(deposit.DestChain)
    }

    packet := ibc.CreatePacketTx {
        DestChain: string(deposit.DestChain),
        Permissions: 
        Tx: deposit
    }

    ibcCtx := ctx.WithPermissions(ibc.AllowIBC(NameETGate)) // NameETEnd?
    _, err := next.DeliverTx(ibcCtx, store, packet.Wrap())
    if err != nil {
        return err
    }
}

func (h Handler) withdrawTx(ctx sdk.Context, store state.SimpleDB, tx WithdrawTx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    setWithdraw()
}

func (h Handler) transferTx(ctx sdk.Context, store state.SimpleDB, tx TransferTx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    setTransfer()
    for _, coin := range tx.Value {
        
    }
}

func (h Handler) setGenesis(store state.SimpleDB, value string) (log string, err error) {
    // TODO: check sender is a validator

    set := NewChainSet(store)

    if set.IsInitialized() {
        return res, ErrAlreadyInitialized()
    }

    var header Header
    err = data.FromJSON([]byte(value), &header)
    if err != nil {
        return "", err
    }

    set.ToBuffer(header)
    set.Finalize(header)
    set.Initialize(header)

    return res, nil
}

func (h Handler) InitState(l log.Logger, store state.SimpleDB, module, key, value string, cb sdk.InitStater) (log string, err error) {
    if module != NameETGate {
        return "", errors.ErrUnknownModule(module)
    }

    switch key {
    // should be the block header that the contract is deployed
    // it can be not real genesis (height 0)
    case "genesis":   
        return setGenesis(store, value)
    }
    return "", 
}
