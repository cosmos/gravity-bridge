package etgate

import (
    "fmt"
    "strings"
    "errors"
    "net/url"
    "strconv"
    "bytes"

    abci "github.com/tendermint/abci/types"
    "github.com/tendermint/basecoin/types"
    "github.com/tendermint/go-wire"
    cmn "github.com/tendermint/tmlibs/common"

    eth "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
//    "github.com/ethereum/go-ethereum/accounts/abi"
//    "github.com/ethereum/go-ethereum/consensus/ethash"
//    "github.com/ethereum/go-ethereum/core"
    "./abi"
    "../../contracts"
)

const (
    _ETGATE = "etgate"
    _BLOCKCHAIN = "blockchain"
    _GENESIS = "genesis"
    _CONFIRM = "confirm"
    _BUFFER = "buffer"
    _RECENT = "recent"
//    _CONTRACT = "contract"
    _DEPOSIT = "deposit"
    _WITHDRAW = "withdraw"

    confirmation = 12

)

var depositabi abi.ABI

func init() {
    deposit, err := abi.JSON(strings.NewReader(contracts.ETGateABI))
    if err != nil {
        panic(err)
    }
    depositabi = deposit
}

type Contract struct {
    Address common.Address
    Code string
}

type Header struct {
    ParentHash common.Hash
    Hash common.Hash
    Number uint64
    ReceiptHash common.Hash
    Time uint64
}

type ETGatePluginState struct {
    // @[:etgate, :contract, Address] <~
    // @[:etgate, :blockchain, :buffer, Hash] <~ 
    // @[:etgate, :blockchain, :confirm, Height] <~
    // @[:etgate, :blockchain, :recent] <~
    // @[:etgate, :egress, Dst, Sequence]
    // @[:etgate, :ingress, Src, Sequence] <~
}

type ETGateTx interface {
    Validate() abci.Result
}

var _ = wire.RegisterInterface (
    struct { ETGateTx }{},
//    wire.ConcreteType{ETGateRegisterContractTx{}, ETGateTxTypeRegisterContract},
    wire.ConcreteType{ETGateUpdateChainTx{}, ETGateTxTypeUpdateChain},
    wire.ConcreteType{ETGateWithdrawTx{}, ETGateTxTypeWithdraw},
    wire.ConcreteType{ETGateDepositTx{}, ETGateTxTypeDeposit},
)

type ETGateUpdateChainTx struct {
    Headers [][]byte
}

func (tx ETGateUpdateChainTx) Validate() abci.Result {
    // TODO: ethash.VerifyHeader?
    return abci.OK
}
/*
type ETGateRegisterContractTx struct {
    Contract
}*/
/*
func (tx ETGateRegisterContractTx) Validate() abci.Result {
    codemap := GetCodemap()
    if _, exists := codemap[tx.Code]; !exists {
        return abci.ErrInternalError.AppendLog(cmn.Fmt("Invalid code"))
    }
    return abci.OK
}
*/
type ETGateDepositTx struct {
    LogProof
}

func (tx ETGateDepositTx) Validate() abci.Result {
    return abci.OK
}

type ETGateWithdrawTx struct {
    To [20]byte 
    Value uint64
    Token [20]byte
    ChainID string
    Sequence uint64
}

func (tx ETGateWithdrawTx) Validate() abci.Result {
    return abci.OK
}

//func SaveNewETGatePacket(state types.KVStore, src, dst string, )

const (
    ETGateTxTypeUpdateChain = byte(0x01)
    //ETGateTxTypeRegisterContract = byte(0x02)
    ETGateTxTypeWithdraw = byte(0x03)
    ETGateTxTypeDeposit = byte(0x04)    

    ETGateCodeConflictingChain = abci.CodeType(1001)
    ETGateCodeMissingAncestor = abci.CodeType(1002)
    ETGateCodeExistingHeader = abci.CodeType(1003)
    ETGateCodeInvalidHeader = abci.CodeType(1004)
    ETGateCodeInvalidLogProof = abci.CodeType(1005)
    ETGateCodeLogHeaderNotFound = abci.CodeType(1006)
    ETGateCodeDepositAlreadyExists = abci.CodeType(1007)
    ETGateCodeLogUnpackingError = abci.CodeType(1008)
    ETGateCodeInvalidEventName = abci.CodeType(1009)
    ETGateCodeAlreadyWithdrawn = abci.CodeType(1010)
    ETGateCodeInvalidCoins = abci.CodeType(1011)
    ETGateCodeInvalidSequence = abci.CodeType(1012)
)

type ETGatePlugin struct {
}

func New() *ETGatePlugin {
    return &ETGatePlugin{}
}

func (gp *ETGatePlugin) RunTx(store types.KVStore, ctx types.CallContext, txBytes []byte) (abci.Result) {
    var tx ETGateTx
    
    if err := wire.ReadBinaryBytes(txBytes, &tx); err != nil {
        return abci.ErrBaseEncodingError.AppendLog("AError decoding tx: " + err.Error())
    }

    res := tx.Validate()
    if res.IsErr() {
        return res.PrependLog("Validate failed: ")
    }

    sm := &ETGateStateMachine{store, ctx, abci.OK}

    switch tx := tx.(type) {
    case ETGateUpdateChainTx:
        sm.runUpdateChainTx(tx)
    //case ETGateRegisterContractTx:
     //   sm.runRegisterContractTx(tx)
    case ETGateWithdrawTx:        
        sm.runWithdrawTx(tx)
    case ETGateDepositTx:
        sm.runDepositTx(tx)
    }

    return sm.res
}

type ETGateStateMachine struct {
    store types.KVStore
    ctx types.CallContext
    res abci.Result
}

func (sm *ETGateStateMachine) getAncestor(header Header) (Header, abci.Result) {
    var genesis uint64
    genesisKey := toKey(_ETGATE, _BLOCKCHAIN, _GENESIS)
    exists, err := load(sm.store, genesisKey, &genesis)
    if !exists {
        return Header{}, abci.ErrInternalError.AppendLog(cmn.Fmt("Genesis not exists: %s", err))
    }
    if err != nil {
        return Header{}, abci.ErrInternalError.AppendLog(cmn.Fmt("Error loading genesis: %s", err))
    }

    for i := 0; i < confirmation; i++ {
        if header.Number == genesis {
            return header, abci.OK
        }
        bufferKey := toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, header.ParentHash.Hex())
        exists, err := load(sm.store, bufferKey, &header)
        if err != nil {
            return Header{}, abci.ErrInternalError.AppendLog(cmn.Fmt("Error loading parent header: %s", err))
        }
        if !exists {
            var res abci.Result
            res.Code = ETGateCodeMissingAncestor
            res.Log = "Missing ancestor"
            return Header{}, res
        }
    }
    return header, abci.OK
}

func (sm *ETGateStateMachine) runUpdateChainTx(tx ETGateUpdateChainTx) {
    for _, headerb := range tx.Headers {
        var header eth.Header
        if err := rlp.DecodeBytes(headerb, &header); err != nil {
            sm.res.Code = ETGateCodeInvalidHeader
            sm.res.Log = "Invalid header"
            return
        }
        res := sm.updateHeader(Header {
            ParentHash: header.ParentHash,
            Hash: header.Hash(),
            ReceiptHash: header.ReceiptHash,
            Number: header.Number.Uint64(),
            Time: header.Time.Uint64(),
        })
        if res.IsErr() {
            sm.res = res.PrependLog(cmn.Fmt("In %vth header: ", header.Number))
            return
        }
    }
}

func (sm *ETGateStateMachine) updateHeader(header Header) abci.Result {
    var res abci.Result

    bufferKey := toKey(_ETGATE, _BLOCKCHAIN, _BUFFER, header.Hash.Hex())
    recentKey := toKey(_ETGATE, _BLOCKCHAIN, _RECENT)

    if !exists(sm.store, bufferKey) {
        save(sm.store, bufferKey, header)
    }

    // TODO: use InitChain to submit genesis, delete this code
    genesisKey := toKey(_ETGATE, _BLOCKCHAIN, _GENESIS)
    if !exists(sm.store, genesisKey) { // genesis
        confirmKey := toKey(_ETGATE, _BLOCKCHAIN, _CONFIRM, strconv.FormatUint(header.Number, 10))
        save(sm.store, recentKey, header.Number)
        save(sm.store, genesisKey, header.Number)
        save(sm.store, confirmKey, header)
        return abci.OK
    }


    ancestor, res := sm.getAncestor(header)
    if res.IsErr() {
        return res
    }
    
    var confirmed Header
    confirmKey := toKey(_ETGATE, _BLOCKCHAIN, _CONFIRM, strconv.FormatUint(ancestor.Number, 10))
    exists, err := load(sm.store, confirmKey, &confirmed)
    if err != nil {
        return abci.ErrInternalError.AppendLog(cmn.Fmt("Loading confirmed header: %s", err))
    }
    
    if exists {    
        if confirmed.Hash != ancestor.Hash {
            res.Code = ETGateCodeConflictingChain
            res.Log = "Conflicting chain"
            return res
        } 
    } else {
        save(sm.store, confirmKey, ancestor)
    }

    save(sm.store, recentKey, ancestor.Number)
  
    return abci.OK
}
/*
func (sm *ETGateStateMachine) runRegisterContractTx(tx ETGateRegisterContractTx) {
    
    var code string
    conKey := toKey(_ETGATE, _CONTRACT, tx.Address.Str())
    exists, err := load(sm.store, conKey, &code)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading code of address: %+v", tx.Address))
        return
    }
    if !exists {
        save(sm.store, conKey, tx.Code)
    } else if code != tx.Code {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Contract already registered with diffrent code: %s", code)) // TODO: change to code
        return
    }
    // does nothing if contract already registered
}
*/
func (sm *ETGateStateMachine) runDepositTx(tx ETGateDepositTx) {
   
    log, err := tx.Log()
    if err != nil {
        sm.res.Code = ETGateCodeInvalidLogProof
        sm.res.Log = "Invalid log proof"
        return
    }
   
    var header Header
    confirmKey := toKey(_ETGATE, _BLOCKCHAIN, _CONFIRM, strconv.FormatUint(tx.Number, 10))
    exi, err := load(sm.store, confirmKey, &header)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading corresponding header to submitted log"))
        return
    }
    if !exi {
        sm.res.Code = ETGateCodeLogHeaderNotFound
        sm.res.Log = "Log header not found"
        return
    }
    if !tx.IsValid(header.ReceiptHash) {
        sm.res.Code = ETGateCodeInvalidLogProof
        sm.res.Log = "Invalid log proof"
        return
    }
/*
    var code string
    conKey := toKey(_ETGATE, _CONTRACT, log.Address.Str())
    exists, err = load(sm.store, conKey, &code)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog(cmn.Fmt("Loading code of log: %s", log.Address.Str()))
        return
    }
    if !exists { 
        sm.res.Code = ETGateCodeUnregisteredContract
        sm.res.Log = "Unregistered Contract"
        return
    }

    codemap := GetCodemap()

    res := codemap[code].Run(sm, tx.Name, log)
    if res.IsErr() {
        sm.res = res.PrependLog("runDepositTx failed: ")
        return
    }
*/
  
    deposit := new(Deposit)
    if err := depositabi.Unpack(deposit, "Deposit", log); err != nil {
        sm.res.Code = ETGateCodeLogUnpackingError
        sm.res.Log = "Log unpacking error"
        return
    }

    depositKey := toKey(_ETGATE, _DEPOSIT, 
        deposit.Chain, 
        strconv.FormatUint(deposit.Seq, 10))

    if exists(sm.store, depositKey) {
        sm.res.Code = ETGateCodeDepositAlreadyExists
        sm.res.Log = "Deposit already exists"
        return 
    }
    
    save(sm.store, depositKey, deposit)
}

type Deposit struct {
    To [20]byte
    Value uint64
    Token common.Address
    Chain string
    Seq uint64
}

func WithdrawValue(tx ETGateWithdrawTx) ([]byte, error) {
    // store as simplest as possible
    buf, n, err := new(bytes.Buffer), new(int), new(error)
    wire.WriteByteSlice(tx.To[:], buf, n, err)
    wire.WriteUint64(tx.Value, buf, n, err)
    wire.WriteByteSlice(tx.Token[:], buf, n, err)
    if *err != nil {
        return []byte{}, *err
    }
    return buf.Bytes(), nil
}

func (sm *ETGateStateMachine) runWithdrawTx(tx ETGateWithdrawTx) {
    tokenHex := "0x"
    for i, _ := range tx.Token {
        tokenHex = tokenHex + string(tx.Token[i]-32)
    }
    if len(sm.ctx.Coins) != 1 || uint64(sm.ctx.Coins[0].Amount) != tx.Value || common.HexToAddress(sm.ctx.Coins[0].Denom) != tx.Token {
        sm.res.Code = ETGateCodeInvalidCoins
        sm.res.Log = fmt.Sprintf("Invalid coins, %v", len(sm.ctx.Coins) /*!= 1, uint64(sm.ctx.Coins[0].Amount) != tx.Value, common.HexToAddress(sm.ctx.Coins[0].Denom) != tx.Token*/)
        return
    }

    sequenceKey := toKey(_ETGATE, _WITHDRAW, tx.ChainID)

    var seq uint64
    exi, err := load(sm.store, sequenceKey, &seq)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog("Error loading sequence")
        return
    }

    if !(exi && tx.Sequence == seq+1 || !exi && tx.Sequence == 1) {
        sm.res.Code = ETGateCodeInvalidSequence
        sm.res.Log = "Invalid sequence"
    }

    withdrawKey := toKey(_ETGATE, _WITHDRAW, tx.ChainID, cmn.Fmt("%v", tx.Sequence))

    if exists(sm.store, withdrawKey) {
        sm.res.Code = ETGateCodeAlreadyWithdrawn
        sm.res.Log = "Already withdrawn"
        return
    }

    val, err := WithdrawValue(tx)
    if err != nil {
        sm.res = abci.ErrInternalError.AppendLog("Error writing withdraw key-value")
        return
    }

    sm.store.Set(withdrawKey, val)
    save(sm.store, sequenceKey, tx.Sequence)
}

func (gp *ETGatePlugin) Name() string{
    return "ETGATE"
}

func (gp *ETGatePlugin) SetOption(store types.KVStore, key string, value string) (log string) {
    return ""
}

func (gp *ETGatePlugin) InitChain(store types.KVStore, vals []*abci.Validator) {
    // TODO: save ethereum genesis block
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
