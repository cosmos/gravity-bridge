package etgate 

import (
//    sdk "github.com/cosmos/cosmos-sdk"
    "github.com/cosmos/cosmos-sdk/stack"
    "github.com/cosmos/cosmos-sdk/state"
    wire "github.com/tendermint/go-wire"

    "github.com/ethereum/go-ethereum/common"
)

/*
- i: info
    - l: last finalized, header
    - g: genesis, header
    - w: withdraw, uint
    - t: transfer, uint
- b: buffer, hash => header
- f: final, uint => header
- w: withdraw, uint => Withdraw
- t: transfer, uint => Transfer
*/

const (
    prefixInfo   = "i"
    infoLast     = "l"
    infoGenesis  = "g"
    infoWithdraw = "w"
    infoTransfer = "t"

    prefixBuffer   = "b"
    prefixFinal    = "f"
    prefixWithdraw = "w"
    prefixTransfer = "t"
)

var (
    
)

type Header struct {
    ParentHash common.Hash
    Hash common.Hash
    Number uint64
    ReceiptHash common.Hash
}

type ChainSet struct {
    InfoSet   *state.Set
    BufferSet *state.Set
    FinalSet  *state.Set
}

func NewChainSet(store state.SimpleDB) ChainSet {
    infoSpace   := stack.PrefixedStore(prefixInfo,   store)
    bufferSpace := stack.PrefixedStore(prefixBuffer, store)
    finalSpace  := stack.PrefixedStore(prefixFinal,  store)
    return ChainSet {
        InfoSet:   state.NewSet(infoSpace),
        BufferSet: state.NewSet(bufferSpace),
        FinalSet:  state.NewSet(finalSpace),
    }
}

func (c ChainSet) GetGenesis() (genesis Header, err error) {

}

func (c ChainSet) LastFinalized() (lf uint, err error) {
    d := set.InfoSet.Get([]byte(infoLast))
    if len(d) == 0 {
        return 0, ErrNotInitialized()
    }

    err = wire.ReadBinaryBytes(d, &lf)
    return lf, err
}

// second argument == nth parent
func (c ChainSet) GetAncestor(header Header, genesis Header) (ancestor Header, err error) {
    for i = 0; i < confirmation; i++ {
        ancestor, err = c.Parent(header)
        if err != nil {
            return
        }
        if ancestor.Hash == genesis.Hash {
            return // genesis is finalized, so we take next
        }
        header = ancestor
    }
    return 
}

func (c ChainSet) GetHeader(number uint) (header Header, err error) { 
    return c.
}   

func (c ChainSet) Parent(header Header) (parent Header, err error) {
    d := set.BufferSet.Get([]byte(header.ParentHash))
    if len(d) == 0 {
        return parent, ErrParentNotFound(header)
    }

    err = wire.ReadBinaryBytes(d, &parent)
    return parent, err
}

func (c ChainSet) ToBuffer(header Header) error {
    h := header.Hash.Hex()
    if set.BufferSet.Exists([]byte(h)) {
        return ErrAlreadyBuffered(header)
    }
    set.BufferSet.Set([]byte(h), header)
    return nil
}

func (c ChainSet) Finalize(header Header) error {
    n := strconv.FormatUint(header.Number, 10)
    if set.FinalSet.Exists([]byte(n)) {
        return ErrAlreadyFinalized(header)
    }
    set.BufferSet.Set([]byte(n), header)
    return nil  
}

// Must check IsInitialized first
func (c ChainSet) Initialize(header Header) {
    set.ToBuffer(header)
    set.Finalize(header)
    set.InfoSet.Set([]byte(infoGenesis), header)
}

func (c ChainSet) IsInitialized() bool {
    return set.InfoSet.Exists([]byte(infoGenesis))
}

func (c ChainSet) UpdateBuffer(header Header) error {
    if c.BufferSet.Exists([]byte(header.Hash)) {
        return ErrConflictingChain(header.Hash)
    }
    c.BufferSet.Set([]byte(header.Hash), header)
}

func (c ChainSet) setWithdraw(tx WithdrawTx) error {

}

func (c ChainSet) setTransfer(tx TransferTx) error {

}
