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

    bctypes "github.com/tendermint/basecoin/types"

    "./abi"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
)

const (
    _CODE = "code"
)

type Code struct {
    abi abi.ABI 
    run func(*ETGateStateMachine, interface{}) abci.Result
    format map[string]reflect.Type
}

func (code Code) Run(sm *ETGateStateMachine, name string, log types.Log) abci.Result {
    v := reflect.New(code.format[name]) 
    if err := code.abi.Unpack(v, name, log); err != nil {
        return abci.ErrInternalError.AppendLog("Error unpacking log")
    }
    return code.run(sm, v)
}

func (code Code) Query(addr common.Address) []ethereum.FilterQuery {
    res := make([]ethereum.FilterQuery, len(code.abi.Events))
    for _, event := range code.abi.Events {
        res = append(res, ethereum.FilterQuery {
            Addresses: []common.Address{addr},
            Topics: [][]common.Hash{{event.Id()}},
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
        "test": &Code {
            run: runTestCode,
            format: map[string]reflect.Type {
                "Event": reflect.TypeOf(TestEvent{}),
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

type TestEvent struct {
    N *big.Int
    B bool
    S string
}

type TokenDeposit struct {
    From common.Address
    To common.Address
    Value *big.Int
    Name string
}

func runTestCode(sm *ETGateStateMachine, test interface{}) abci.Result {
    switch test := test.(type) {
    case TestEvent:
        codeKey := toKey(_ETGATE, _CODE, "test", test.N.String())
        save(sm.store, codeKey, test)
    default:
        return abci.ErrInternalError.AppendLog("Type error in runTestCode")
    }
    return abci.OK
}

func runTokenCode(sm *ETGateStateMachine, token interface{}) abci.Result {
    switch token := token.(type) {
    case TokenDeposit:
        acc := bctypes.GetAccount(sm.store, token.To[:])
        if acc == nil {
            return abci.ErrInternalError.AppendLog("Destination address does not exist")
        }
        coins := bctypes.Coins{bctypes.Coin{Denom: token.Name, Amount: token.Value.Int64()}}
        acc.Balance = acc.Balance.Plus(coins)
        bctypes.SetAccount(sm.store, token.To[:], acc)
    default:
        return abci.ErrInternalError.AppendLog("Type error in runTokenCode")
    }
    return abci.OK
}
