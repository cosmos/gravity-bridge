package etgate

import (
    "fmt"

    abci "github.com/tendermint/abci/types"
    "github.com/tendermint/basecoin/types"
    "github.com/tendermint/go-wire"

    eth "github.com/ethereum/go-ethereum/core/types"
//    "github.com/ethereum/go-ethereum/consensus/ethash"
    "github.com/ethereum/go-ethereum/core"
)

const (
    _ETGATE = "etgate"
    _GENESIS = "genesis"
    _CONFIRM = "confirm"
    _BUFFER = "buffer"

    confirmation = 12
)

type ETGatePluginState struct {

}

type ETGateTx interface {
    Validate() abci.Result
}

type ETGateUpdateChainTx struct {
    Header eth.Header
}

func (tx ETGateUpdateChainTx) Validate() abci.Result {
    // TODO: ethash.VerifyHeader?
    return
}

type ETGateRegisterTokenTx struct {

}

type ETGateDepositTokenTx struct {
    Proof LogProof
}

func (tx ETGateDepositTokenTx) Validate() abci.Result {
    return
}

//type ETGateWithdrawTokenTx struct {
//}

const (
    ETGateTxTypeUpdateChainTx = byte(0x01)
//    ETGateTxTypeRegisterTokenTx = byte(0x02)
    ETGateTxTypeDepositTokenTx = byte(0x03)
//    ETGateTxTypeWithdrawTokenTx = byte(0x04)    

    ETGateCodeConfilctingChain = abci.CodeType(1001)
)

type ETGatePlugin struct {
    genesis core.Genesis
}

func (gp *ETGatePlugin) Name() string {
    return "etgate"
}

func New() *ETGatePlugin {
    return &ETGatePlugin{}
}

func (gp *ETGatePlugin) RunTx(store types.KVStore, ctx types.CallContext, txBytes []byte) (abci.Result) {
    var tx ETGateTx
    
    if err := wire.ReadBinaryBytes(txBytes, &tx); err != nil {
        return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
    }
    
    sm := &ETGateStateMachine{store, ctx, abci.OK}

    switch tx := tx.(type) {
    case ETGateUpdateChainTx:
        sm.runUpdateChainTx(gp, tx)
    case ETGateDepositTokenTx:
        sm.runDepositTokenTx(gp, tx)
//    case ETGateWithdrawTokenTx:
//        sm.runWithdrawTokenTx(tx)
    }

    return sm.res
}

type ETGateStateMachine struct {
    store types.KVStore
    ctx types.CallContext
    res abci.Result
}

func (sm *ETGateStateMachine) runUpdateChainTx(gp *ETGatePlugin, tx ETGateUpdateChainTx) {
    hash := tx.Header.Hash()
    bufferKey := toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, hash)

    ancestor := tx.Header
    for i := 0; i < confirmation; i++ {
        bufferKey = toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, ancestor.Header.ParentHash)
        exists, err := load(sm.store, bufferKey, &ancestor)
        if err != nil {
            sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading ancestor header: %+v", err.Error()))
        }
        if !exists {
            sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Missing ancestor header"))
        }
    }

    confirmKey := toKey(_ETGATE, _BLOCKCHAIN, _CONFIRM, ancestor.Number)
    if exists(sm.store, confirmKey) {
        sm.res.Code = ETGateCodeConflictingChain
        sm.res.Log = "Conflicting chain"
        return
    }

    bufferKey = toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, hash)
    save(sm.store, bufferKey, tx.Header)
    save(sm.store, confirmKey, ancestor)
}

func (gp *ETGatePlugin) InitChain(store types.KVStore, vals []*abci.Validator) {
}

func (gp *ETGatePlugin) BeginBlock(store types.KVStore, vals []*abci.Validator) {
}

func (gp *ETGatePlugin) EndBlock(store types.KVStore, height uint64) (res abci.ResponseEndBlock) {
    return
}


