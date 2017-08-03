package etgate 

import (
    "os"
    "fmt"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "reflect"
    "strconv"

    abci "github.com/tendermint/abci/types"

    bctypes "github.com/tendermint/basecoin/types"

    "./abi"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
)

const (
    _CODE = "code"
)

type Query struct {
   ethereum.FilterQuery
   Name string
}

type Code struct {
    abi abi.ABI 
    run func(*ETGateStateMachine, common.Address, Payload) abci.Result
    format map[string]reflect.Type
}

func savePayload(sm *ETGateStateMachine, addr common.Address, payload Payload) abci.Result {
    packetKeyIngress := toKey(_ETGATE, _INGRESS,
        addr.Hex(),
        strconv.FormatUint((payload).Sequence(), 10),
    )
    if exists(sm.store, packetKeyIngress) {
        var res abci.Result
        res.Code = ETGateCodePacketAlreadyExists
        res.Log = "Packet already exists"
        return res
    }
    save(sm.store, packetKeyIngress, payload)
    return abci.OK
}

func (code Code) Run(sm *ETGateStateMachine, name string, log types.Log) abci.Result {
    f, ok := code.format[name]
    if !ok {
        var res abci.Result
        res.Code = ETGateCodeInvalidEventName
        res.Log = "Invalid event name"
        return res
    }
    v := reflect.New(f) 
    if !v.CanInterface() {
        return abci.ErrInternalError.AppendLog("Cannot convert reflect value to interface")
    }
    payload := v.Interface().(Payload)

    if err := code.abi.Unpack(payload, name, log); err != nil {
        return abci.ErrInternalError.AppendLog("Error unpacking log")
    }

    if res := savePayload(sm, log.Address, payload); res.IsErr() {
        return res
    }

    return code.run(sm, log.Address, payload)
}

func (code Code) Query(addr common.Address) []Query {
    res := []Query{}
    for _, event := range code.abi.Events {
        res = append(res, Query {
            Name: event.Name,
            FilterQuery: ethereum.FilterQuery {
                Addresses: []common.Address{ addr },
                Topics: [][]common.Hash{{ event.Id() }},
            },
        })
    }
    return res
}

var codemap map[string]*Code

func GetCodemap() map[string]*Code {
    if codemap != nil {
        return codemap
    }

    plan, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "abimap.json"))
    if err != nil {
        panic(fmt.Errorf("Error reading abimap: ", err))
    }

    var rawmap map[string]*json.RawMessage
    if err := json.Unmarshal(plan, &rawmap); err != nil {
        panic(fmt.Errorf("Error reading abimap: ", err))
    }

    abimap := map[string]abi.ABI{}
    for k, v := range rawmap {
        var abi abi.ABI
        if err := abi.UnmarshalJSON(*v); err != nil {
            panic(fmt.Errorf("Error reading abimap: ", err))
        }
        abimap[k] = abi
    }

    codemap = genCodemap(abimap)
    return codemap
}

func genCodemap(abimap map[string]abi.ABI) map[string]*Code {
    codes := map[string]*Code {
        "data": &Code {
            run: runDataCode,
            format: map[string]reflect.Type {
                "Submit": reflect.TypeOf(DataSubmit{}),
            },
        },
        "token": &Code {
            run: runTokenCode,
            format: map[string]reflect.Type {
                "Deposit": reflect.TypeOf(TokenDeposit{}),
            },
        },
    }

    for k, _ := range codes {
        codes[k].abi = abimap[k]
    }

    return codes
}

type Payload interface {
    Sequence() uint64
}

type DataSubmit struct {
    Data []byte
    Seq uint64
}

func (ds DataSubmit) Sequence() uint64 {
    return ds.Seq
}

type TokenDeposit struct {
    To common.Address // Assuming (basecoin)commands.Address == (ethereum)common.Address == [20]byte
    Value uint64
    Addr common.Address // token contract address
    Seq uint64
}

func (td TokenDeposit) Sequence() uint64 {
    return td.Seq
}


func runDataCode(sm *ETGateStateMachine, addr common.Address, data Payload) abci.Result {
    switch data.(type) {
    case *DataSubmit:
        return abci.OK
    default:
        return abci.ErrInternalError.AppendLog("Type error in runDataCode")
    }
    return abci.OK
}


func runTokenCode(sm *ETGateStateMachine, addr common.Address, token Payload) abci.Result {
    switch token := token.(type) {
    case *TokenDeposit:
        acc := bctypes.GetAccount(sm.store, token.To[:])
        if acc == nil {
            return abci.ErrInternalError.AppendLog("Destination address does not exist")
        }
        coins := bctypes.Coins{bctypes.Coin{Denom: fmt.Sprintf("%s,%s", token.Addr.Hex(), addr.Hex()), Amount: int64(token.Value)}} // TODO: check uint64->int64 information loss
        acc.Balance = acc.Balance.Plus(coins)
        bctypes.SetAccount(sm.store, token.To[:], acc)
    default:
        return abci.ErrInternalError.AppendLog("Type error in runTokenCode")
    }
    return abci.OK
}
