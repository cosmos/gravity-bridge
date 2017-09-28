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
    stack.PassInitState
}

var _ stack.Dispatchable = Handler{}

func (Handler) Name() string {
    return NameETGate
}

func (Handler) AssertDispather() {}

// https://github.com/cosmos/cosmos-sdk/blob/develop/modules/ibc/handler.go

func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx) (res sdk.DeliverResult, err error) {
    err := tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch tx.Unwrap().(type) {
    case ChainTx:
        return res, nil
    // case UpdateValidatorTx:
    case DepositTx:
        return res, nil
    }

    return res, errors.ErrUnknownTxType(tx.Unwrap())
}

func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx) (res sdk.DeliverResult, err error) {
    err := tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch t := tx.Unwrap().(type) {
    case InitTx:
        return h.registerTx(ctx, store, t)
    case UpdateTx:
        return h.updateTx(ctx, store, t)
    // case ValChangeTx:
    case DepositTx:
        return h.depositTx(ctx, store, t)
    // case WithdrawTx? How can I dispatch an incoming IBC packet?
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

func (h Handler) registerTx(ctx sdk.Context, store state.SimpleDB, t InitTx) (res sdk.DeliverResult, err error) {
    // TODO: check sender is a validator

    set := NewChainSet(store)

    if set.IsInitialized() {
        return res, errAlreadyInitialized()
    }

    header, err := decodeHeader(tx.Header)
    if err != nil {
        return res, err
    }

    set.ToBuffer(header)
    set.Finalize(header)
    set.Initialize(header)

    return res, nil
}

func (h Handler) updateTx(ctx sdk.Context, store state.SimpleDB, t UpdateTx) (res sdk.DeliverResult, err error) {
    // TODO: check sender is a validator

    set := NewChainSet(store)

    if !set.IsInitialized() {
        return res, errNotInitialized()
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

func LoadState
