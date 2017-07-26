package etgate

import (
//    "fmt"
    "strings"
    "errors"
    "net/url"

    abci "github.com/tendermint/abci/types"
    "github.com/tendermint/basecoin/types"
    "github.com/tendermint/go-wire"
    cmn "github.com/tendermint/tmlibs/common"

    eth "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
//    "github.com/ethereum/go-ethereum/accounts/abi"
//    "github.com/ethereum/go-ethereum/consensus/ethash"
//    "github.com/ethereum/go-ethereum/core"
)

const (
    _ETGATE = "etgate"
    _BLOCKCHAIN = "blockchain"
    _GENESIS = "genesis"
    _CONFIRM = "confirm"
    _BUFFER = "buffer"
    _CONTRACT = "contract"

    confirmation = 12
)

type ETGatePluginState struct {

}

type ETGateTx interface {
    Validate() abci.Result
}

var _ = wire.RegisterInterface (
    struct { ETGateTx }{},
    wire.ConcreteType{ETGateRegisterContractTx{}, ETGateTxTypeRegisterContract},
    wire.ConcreteType{ETGateUpdateChainTx{}, ETGateTxTypeUpdateChain},
//    wire.ConcreteType{ETGatePacketCreateTx{}, ETGateTxTypePacketCreate},
    wire.ConcreteType{ETGatePacketPostTx{}, ETGateTxTypePacketPost},
)

type ETGateUpdateChainTx struct {
    Header eth.Header
}

func (tx ETGateUpdateChainTx) Validate() abci.Result {
    // TODO: ethash.VerifyHeader?
    return abci.OK
}

type ETGateRegisterContractTx struct {
    Address string 
    Code string
}

func (tx ETGateRegisterContractTx) Validate() abci.Result {
    if !common.IsHexAddress(tx.Address) {
        return abci.ErrInternalError.AppendLog("Invalid address format")
    }
    codemap := get()
    if _, exists := codemap[tx.Code]; !exists {
        return abci.ErrInternalError.AppendLog(cmn.Fmt("Invalid code"))
    }
    return abci.OK
}

type ETGatePacketPostTx struct {
    Name string
    Proof LogProof
}

func (tx ETGatePacketPostTx) Validate() abci.Result {
    return abci.OK
}

//type ETGatePacketPostTx struct {
//}

const (
    ETGateTxTypeUpdateChain = byte(0x01)
    ETGateTxTypeRegisterContract = byte(0x02)
//    ETGateTxTypePacketCreate = byte(0x03)
    ETGateTxTypePacketPost = byte(0x04)    

    ETGateCodeConflictingChain = abci.CodeType(1001)
)

type ETGatePlugin struct {
}

func New() *ETGatePlugin {
    return &ETGatePlugin{}
}

func (gp *ETGatePlugin) RunTx(store types.KVStore, ctx types.CallContext, txBytes []byte) (abci.Result) {
    var tx ETGateTx
    
    if err := wire.ReadBinaryBytes(txBytes, &tx); err != nil {
        return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
    }

    res := tx.Validate()
    if res.IsErr() {
        return res.PrependLog("Validate failed: ")
    }

    sm := &ETGateStateMachine{store, ctx, abci.OK}

    switch tx := tx.(type) {
    case ETGateUpdateChainTx:
        sm.runUpdateChainTx(tx)
    case ETGateRegisterContractTx:
        sm.runRegisterContractTx(tx)
//    case ETGatePacketCreateTx:
//        sm.runPacketCreateTx(tx)
    case ETGatePacketPostTx:
        sm.runPacketPostTx(tx)
    }

    return sm.res
}

type ETGateStateMachine struct {
    store types.KVStore
    ctx types.CallContext
    res abci.Result
}

func (sm *ETGateStateMachine) runUpdateChainTx(tx ETGateUpdateChainTx) {
    hash := tx.Header.Hash().Str()
    bufferKey := toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, hash)

    ancestor := tx.Header
    for i := 0; i < confirmation; i++ {
        bufferKey = toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, ancestor.ParentHash.Str())
        exists, err := load(sm.store, bufferKey, &ancestor)
        if err != nil {
            sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading ancestor header: %+v", err.Error()))
            return
        }
        if !exists {
            sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Missing ancestor header"))
            return
        }
    }

    confirmKey := toKey(_ETGATE, _BLOCKCHAIN, _CONFIRM, ancestor.Number.String())
    if exists(sm.store, confirmKey) {
        sm.res.Code = ETGateCodeConflictingChain
        sm.res.Log = "Conflicting chain"
        return
    }

    bufferKey = toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, hash)
    save(sm.store, bufferKey, tx.Header)
    save(sm.store, confirmKey, ancestor)
}

func (sm *ETGateStateMachine) runRegisterContractTx(tx ETGateRegisterContractTx) {
    
    var code string
    conKey := toKey(_ETGATE, _CONTRACT, tx.Address)
    exists, err := load(sm.store, conKey, &code)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading code of address: %+v", tx.Address))
        return
    }
    if !exists {
        save(sm.store, conKey, tx.Code)
    } else if code != tx.Code {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Contract already registered with diffrent code: %s", code))
        return
    }
    // does nothing if contract already registered
}

func (sm *ETGateStateMachine) runPacketPostTx(tx ETGatePacketPostTx) {
    var code string
    log := tx.Proof.Log()
    conKey := toKey(_ETGATE, _CONTRACT, log.Address.Str())
    exists, err := load(sm.store, conKey, &code)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading code of log: %s", log.Address.Str()))
        return
    }
    if !exists { 
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Contract not registered: %s", log.Address.Str()))
        return
    }
   
    codemap := get()

    res := codemap[code].Run(sm, tx.Name, log)
    if res.IsErr() {
        sm.res = res.PrependLog("runPacketPostTx failed: ")
        return
    }
}

func (gp *ETGatePlugin) Name() string{
    return "ETGATE"
}

func (gp *ETGatePlugin) SetOption(store types.KVStore, key string, value string) (log string) {
    return ""
}

func (gp *ETGatePlugin) InitChain(store types.KVStore, vals []*abci.Validator) {
}

func (gp *ETGatePlugin) BeginBlock(store types.KVStore, hash []byte, header *abci.Header) {
}

func (gp *ETGatePlugin) EndBlock(store types.KVStore, height uint64) (res abci.ResponseEndBlock) {
    return
}

// https://github.com/tendermint/basecoin/blob/master/plugins/ibc/ibc.go

// Returns true if exists, false if nil.
func exists(store types.KVStore, key []byte) (exists bool) {
    value := store.Get(key)
    return len(value) > 0
}

// Load bytes from store by reading value for key and read into ptr.
// Returns true if exists, false if nil.
// Returns err if decoding error.
func load(store types.KVStore, key []byte, ptr interface{}) (exists bool, err error) {
    value := store.Get(key)
    if len(value) > 0 {
        err = wire.ReadBinaryBytes(value, ptr)
        if err != nil {
            return true, errors.New(
                cmn.Fmt("Error decoding key 0x%X = 0x%X: %v", key, value, err.Error()),
            )
        }
        return true, nil
    } else {
        return false, nil
    }
}

// Save bytes to store by writing obj's go-wire binary bytes.
func save(store types.KVStore, key []byte, obj interface{}) {
    store.Set(key, wire.BinaryBytes(obj))
}

func toKey(parts ...string) []byte {
    escParts := make([]string, len(parts))
    for i, part := range parts {
        escParts[i] = url.QueryEscape(part)
    }
    return []byte(strings.Join(escParts, ","))
}
