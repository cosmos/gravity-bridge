package etgate 

import (
    "os"
    "fmt"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "reflect"
    "math/big"

    abci "github.com/tendermint/abci/types"

    "./abi"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
)

type Code struct {
    abi abi.ABI 
    run func(abi.ABI, *ETGateStateMachine, types.Log) abci.Result
    format map[string]reflect.Type
}

func (code Code) Run(sm *ETGateStateMachine, name string, log types.Log) abci.Result {
    v := reflect.New(code.format[name]) 
    if err := code.abi.Unpack(v, name, log); err != nil {
        return abci.ErrInternalError.AppendLog("Error unpacking log")
    }
    return code.run(code.abi, sm, log)
}

func (code Code) Query() []ethereum.FilterQuery {
    return nil
}

var codemap map[string]*Code

func get() map[string]*Code {
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
        "test": &Code {
            run: runTestCode,
            format: map[string]reflect.Type {
                "Event": reflect.TypeOf(struct {
                    N *big.Int
                    B bool
                    S string
                }{}),
            },
        },
        "token": &Code {
            run: runTokenCode,
            format: map[string]reflect.Type {
                "Deposit": reflect.TypeOf(struct {
                    _from common.Address // original eth address
                    _to common.Address // destination mint address
                    _value *big.Int
                }{}),
            },
        },
    }

    for k, _ := range codes {
        codes[k].abi = abimap[k]
    }

    return codes
}


func runTestCode(abi abi.ABI, sm *ETGateStateMachine, log types.Log) abci.Result {
//    abi.Unpack
    return abci.OK
}

func runTokenCode(abi abi.ABI, sm *ETGateStateMachine, log types.Log) abci.Result {
    return abci.OK
}
