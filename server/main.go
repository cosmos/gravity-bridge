package main

import (
    "fmt"
    "os"
    "errors"
    "net/http"
    "encoding/json"
    "encoding/hex"
    "path/filepath"

    "github.com/spf13/viper"

    "github.com/gorilla/mux"

    bctypes "github.com/tendermint/basecoin/types"
    bclicmd "github.com/tendermint/basecoin/cmd/basecli/commands"

    "github.com/tendermint/go-crypto/keys"
    keyserver "github.com/tendermint/go-crypto/keys/server"
    "github.com/tendermint/go-crypto/keys/cryptostore"
    "github.com/tendermint/go-crypto/keys/storage/filestorage"

    "github.com/tendermint/go-wire"

//    ctypes "github.com/tendermint/tendermint/rpc/core/types"

    lc "github.com/tendermint/light-client"
    "github.com/tendermint/light-client/commands"
    "github.com/tendermint/light-client/proofs"
    proofcmd "github.com/tendermint/light-client/commands/proofs"

    rpcclient "github.com/tendermint/tendermint/rpc/client"

    "github.com/tendermint/merkleeyes/iavl"

    "github.com/ethereum/go-ethereum/common"

    etcmd "github.com/cosmos/peggy/commands"
    "github.com/cosmos/peggy/plugins/etgate"
)

func GetKeyManager() keys.Manager {
    return cryptostore.New(
        cryptostore.SecretBox,
        filestorage.New(filepath.Join(os.Getenv("HOME"), ".etgate", "client", "keys")),
        keys.MustLoadCodec("english"),
    )
}

func init() {
    viper.Set("home", filepath.Join(os.Getenv("HOME"), ".etgate", "client"))
    viper.Set("node", "tcp://localhost:12347")
    viper.Set("chain-id", "etgate-chain")
    manager = GetKeyManager()
    node = commands.GetNode()
}

var (
    manager keys.Manager
    node    rpcclient.Client
)


type ResultJSON struct {
    Result interface{} `json:"result"`
    Error  string      `json:"error"`
}

func Result(w http.ResponseWriter, result interface{}) {
    json.NewEncoder(w).Encode(ResultJSON{result, ""})
}

func Error(w http.ResponseWriter, tag string, err error) {
    json.NewEncoder(w).Encode(ResultJSON{nil, tag + err.Error()})
}

type PostWithdrawRequest struct {
    Name       string `json:"name"`
    Passphrase string `json:"passphrase"`
//  Origin     string `json:"origin"`  
    To         string `json:"to"`
    Value      int64  `json:"value"`
    Token      string `json:"token"`
    ChainID    string `json:"chainid"`
}

type PostWithdrawResponse struct {
    Height               uint64         `json:"height"`
    IavlProofLeafHash    string         `json:"iavlProofLeafHash"` // [20]byte
    IavlProofInnerHeight []int8         `json:"iavlProofInnerHeight"` // [20]byte
    IavlProofInnerSize   []int          `json:"iavlProofInnerSize"`
    IavlProofInnerLeft   []string       `json:"iavlProofInnerLeft"` // [][20]byte
    IavlProofInnerRight  []string       `json:"iavlProofInnerRight"` // [][20]byte
    IavlProofRootHash    string         `json:"iavlProofRootHash"` // [20]byte
    To                   string         `json:"to"`
    Value                uint64         `json:"value"`
    Token                string         `json:"token"`
    Chain                string         `json:"chain"`
    Seq                  string         `json:"seq"`
}

func Debug(tag string, data interface{}) {
    fmt.Printf("%s: %+v\n", tag, data)
}

func toHex(s []byte) string {
    return fmt.Sprintf("%040s", hex.EncodeToString(s))
}   

func PostWithdraw(w http.ResponseWriter, request *http.Request) {
    if request.Body == nil {
        Error(w, "Request body is empty", nil)
    }
    var req PostWithdrawRequest
    if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
        Error(w, "Error decoding request: ", err)
        return
    }

    originChainID := "etgate-chain" // for now

    sequenceKey := fmt.Sprintf("etgate,withdraw,%s", originChainID)
    query, err := etcmd.QueryWithClient(node, []byte(sequenceKey))
    if err != nil {
        Error(w, "Error querying sequence: ", err)
        return
    }
    var seq uint64
    if len(query.Value) == 0 {
        seq = 0
    } else {
        if err = wire.ReadBinaryBytes(query.Value, &seq); err != nil {
            Error(w, "Error reading sequence: ", err)
            return
        }
    }
    

    info, err := manager.Get(req.Name)
    if err != nil {
        Error(w, "Error getting key info: ", err)
        return
    }

    acc, err := etcmd.GetAccWithClient(node, info.Address[:])
    if err != nil {
        Error(w, "Error getting account sequence: ", err)
        return
    }

    inner := etgate.ETGateWithdrawTx {
        To: common.HexToAddress(req.To),
        Value: uint64(req.Value),
        Token: req.Token,
        ChainID: originChainID,
        Sequence: seq+1,
    }

    feeCoins := bctypes.Coin{Denom: "mycoin", Amount:1}
    enctoken := ""
    for _, s := range req.Token[2:] {
        enctoken = enctoken + string(s+32)
    }
    ethCoins := bctypes.Coin{Denom: enctoken, Amount: int64(req.Value)}

    coins := bctypes.Coins{ethCoins, feeCoins}
    coins.Sort()

    tx := &bctypes.AppTx {
        Gas: 0,
        Fee: feeCoins,
        Name: "ETGATE",
        Input: bctypes.NewTxInput(info.PubKey, coins, acc.Sequence+1),
        Data: wire.BinaryBytes(struct {
            etgate.ETGateTx `json:"unwrap"`
        }{inner}),
    }

    // https://github.com/tendermint/light-client/blob/master/commands/txs/helpers.go 

    apptx := bclicmd.WrapAppTx(tx)

    err = apptx.ValidateBasic()
    if err != nil {
        Error(w, "", err)
        return
    }


    err = manager.Sign(req.Name, req.Passphrase, apptx)
    if err != nil {
        Error(w, "", err)
        return
    }

    packet, err := apptx.TxBytes()
    if err != nil {
        Error(w, "", err)
        return
    }

    _, err = node.BroadcastTxCommit(packet)
    if err != nil {
        Error(w, "", err)
        return
    }

    // returns needed arguments for withdraw()
    key := fmt.Sprintf("etgate,withdraw,%s,%v", originChainID, seq+1)
    Debug("key", key)
    query, err = etcmd.QueryWithClient(node, []byte(key))
    if err != nil {
        Error(w, "", err)
        return
    }
    Debug("query", query)

    proof, err := iavl.ReadProof(query.Proof)
    if err != nil {
        Error(w, "", err)
        return
    }
    Debug("proof", proof)

    var iavlProofInnerHeight []int8
    var iavlProofInnerSize []int
    var iavlProofInnerLeft []string
    var iavlProofInnerRight []string
    for _, in := range proof.InnerNodes {
        iavlProofInnerHeight = append(iavlProofInnerHeight, in.Height)
        iavlProofInnerSize   = append(iavlProofInnerSize,   in.Size)
        iavlProofInnerLeft   = append(iavlProofInnerLeft,   toHex(in.Left))
        iavlProofInnerRight  = append(iavlProofInnerRight,  toHex(in.Right))
    }

    response := PostWithdrawResponse {
        Height:               query.Height,
        IavlProofLeafHash:    toHex(proof.LeafHash),
        IavlProofInnerHeight: iavlProofInnerHeight,
        IavlProofInnerSize:   iavlProofInnerSize,
        IavlProofInnerLeft:   iavlProofInnerLeft,
        IavlProofInnerRight:  iavlProofInnerRight,
        IavlProofRootHash:    toHex(proof.RootHash),
        To:                   req.To,
        Value:                uint64(req.Value),
        Token:                req.Token,
        Chain:                originChainID,
        Seq:                  fmt.Sprintf("%v", seq+1),
    }

    Debug("response", response)

    Result(w, response)    
}
func GetAccount(w http.ResponseWriter, request *http.Request) {
    params := mux.Vars(request)
    info, err := manager.Get(params["name"])
    if err != nil {
        Error(w, "Error getting key info: ", err)
        return
    }

    fmt.Printf("%+v\n", info.Address)

    // https://github.com/cosmos/cosmos-sdk/blob/master/cmd/basecli/commands/query.go

    addr, err := proofs.ParseHexKey(info.Address.String())
    fmt.Printf("%+v\n", addr)
    if err != nil {
        Error(w, "Error parsing key address: ", err)
        return
    }
    key := bctypes.AccountKey(addr)

    acc := new(bctypes.Account)
    _, err = proofcmd.GetAndParseAppProof(key, &acc)
    if lc.IsNoDataErr(err) {
        Error(w, "", errors.New("Account bytes are empty for address " + params["name"]))
        return
    } else if err != nil {
        Error(w, "Error getting address balance: ", err)
        return
    }
    
    Result(w, acc.Balance)
}

func main() {
    r := mux.NewRouter()
    
    // Keys
    k := keyserver.New(manager, "secp256k1")
    r.HandleFunc("/keys", k.GenerateKey).Methods("POST")
    r.HandleFunc("/keys", k.ListKeys).Methods("GET")
    r.HandleFunc("/keys/{name}", k.GetKey).Methods("GET")
    r.HandleFunc("/keys/{name}", k.UpdateKey).Methods("POST", "PUT")
    r.HandleFunc("/keys/{name}", k.DeleteKey).Methods("DELETE")

    // Account
    r.HandleFunc("/query/account/{name}", GetAccount).Methods("GET")

    // Txs
    r.HandleFunc("/withdraw", PostWithdraw).Methods("POST")

    // Frontend
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    http.Handle("/", r)
    http.ListenAndServe(":12349", nil)
}
