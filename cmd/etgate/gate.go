package main

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"
    "context"
    "encoding/json"
    "encoding/hex"
    "io/ioutil"
    "math/big"
    "strings"
//    "bytes"

//    "golang.org/x/crypto/ripemd160"

    "github.com/spf13/cobra"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/ethclient"
//    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    ecrypto "github.com/ethereum/go-ethereum/crypto"
 
//    abci "github.com/tendermint/abci/types"
    cmn "github.com/tendermint/tmlibs/common"
    basecmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
//    bctypes "github.com/tendermint/basecoin/types"
    mintclient "github.com/tendermint/tendermint/rpc/client"
    tmtypes "github.com/tendermint/tendermint/types"
 //   "github.com/tendermint/tmlibs/merkle"

    "github.com/tendermint/go-wire"
    "github.com/tendermint/go-crypto"

    "../../plugins/etgate"
    "../../commands"
    "../../contracts"
    "../../plugins/etgate/abi"

    secp256k1 "github.com/btcsuite/btcd/btcec"
)

var ( 
    one = big.NewInt(1)
    chunksize = big.NewInt(16)
)


type gateway struct {
    ethclient *ethclient.Client
    mintclient *mintclient.HTTP
    ethauth *bind.TransactOpts
    mintkey *basecmd.Key
//    query []etgate.Query
}

var GateCmd = &cobra.Command {
    Use: "gate",
    Short: "Relay ethereum logs to tendermint",
}

var GateStartCmd = &cobra.Command {
    Use: "start",
    Short: "Start etgate relayer to relay ethereum logs to tendermint",
    RunE: gateStartCmd,
}

var GateInitCmd = &cobra.Command {
    Use: "init",
    Short: "Register ethereum contract",
    RunE: gateInitCmd,
}


var (
    testnetFlag bool
    datadirFlag string
    ipcpathFlag string
    nodeaddrFlag string
    chainIDFlag string
    addressFlag string
    genesisFlag string
    depositABI abi.ABI
)

func init() {
    flags := []basecmd.Flag2Register {
        {&testnetFlag, "testnet", false, "Ropsten network: pre-configured test network"},
        {&nodeaddrFlag, "nodeaddr", "tcp://localhost:46657", "Node address for tendermint chain"},
        {&chainIDFlag, "chain-id", "etgate-chain", "Chain ID"},
        {&addressFlag, "address", "", "ETGate contract address on Ethereum chain"},
    }

    basecmd.RegisterPersistentFlags(GateCmd, flags)

    initFlags := []basecmd.Flag2Register {
        {&genesisFlag, "genesis", "", "Path to genesis file"},
    }

    basecmd.RegisterPersistentFlags(GateInitCmd, initFlags)

    startFlags := []basecmd.Flag2Register {
        {&datadirFlag, "datadir", filepath.Join(os.Getenv("HOME"), ".ethereum"), "Data directory for the databases and keystore"},
        {&ipcpathFlag, "ipcpath", "geth.ipc", "Filename for IPC socket/pipe within the datadir"},
    }

    basecmd.RegisterFlags(GateStartCmd, startFlags)

    GateCmd.AddCommand(GateStartCmd)
    GateCmd.AddCommand(GateInitCmd)

    var err error
    depositABI, err = abi.JSON(strings.NewReader(contracts.ETGateABI))
    if err != nil {
        panic(err)
    }

}

func getConsfile(consfile string) ([]map[string]interface{}, error) {
    if !filepath.IsAbs(consfile) {
        wd, err := os.Getwd()
        if err != nil {
            return nil, err
        }
        
        consfile = filepath.Join(wd, consfile)
    }

    plan, err := ioutil.ReadFile(consfile)
    if err != nil {
        return nil, err
    }
    
    var data []map[string]interface{}
    if err := json.Unmarshal(plan, &data); err != nil {
        return nil, err
    }

    return data, nil
}

func gateInitCmd(cmd *cobra.Command, args []string) error {
/*    if len(args) < 1 {
        return errors.New("Usage: etgate gate init [--testnet] contractsfile")
    }
    consfile := args[0]

    data, err := getConsfile(consfile)
    if err != nil {
        return err
    }
*/
    g, err := newGateway()
    if err != nil {
        return err
    }
/*
    // Delete this part
    for _, con := range data {
        addr, code := con["address"].(string), con["code"].(string)
        if !common.IsHexAddress(addr) {
            return errors.New("Invalid address format")
        }

        registerTx := etgate.ETGateRegisterContractTx {
            etgate.Contract {
                Address: common.HexToAddress(addr),
                Code: code,
            },
        }
        if err := g.appTx(registerTx); err != nil {
            return err
        }
    }
    // ^
*/

    genesisBytes, err := ioutil.ReadFile(genesisFlag)
    if err != nil {
        return err
    }

    chainGenDoc := new(tmtypes.GenesisDoc)
    if err = json.Unmarshal(genesisBytes, chainGenDoc); err != nil {
        return err
    }
   
    validatorsBytes := []byte{}
    votingPowers := []*big.Int{}
    for _, val := range chainGenDoc.Validators {
        pub_, err := getSecp256k1Pub(val.PubKey)
        if err != nil {
            return err
        }

        pub, err := secp256k1.ParsePubKey(pub_[:], secp256k1.S256())
        if err != nil {
            return err
        }

        validatorsBytes = append(validatorsBytes, pub.SerializeUncompressed()[:]...)
        votingPowers = append(votingPowers, big.NewInt(val.Amount))
    }

    validatorsHex := "["

    for i := 0; i < len(validatorsBytes); i++ {
        validatorsHex = validatorsHex + "\"0x" + hex.EncodeToString(validatorsBytes[i:i+1]) + "\", "
    }
    
    validatorsHex = validatorsHex+ "]"

    fmt.Printf("%+v\n%+v\n", validatorsHex, votingPowers)

    address, _, _, err := contracts.DeployETGate(g.ethauth, g.ethclient, []byte("etgate-chain"), validatorsBytes, votingPowers)
    if err != nil {
        return err
    }

    fmt.Printf("ETGate contract is deployed on %s", address.Hex())

    return nil
}

func gateStartCmd(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
        return errors.New("Usage: etgate gate start [--testnet] [--datadir ~/.ethereum] [--ipcpath geth.ipc] [--rpcpath localhost:1234] contractsfile")
    } 
   

    gateway, err := newGateway()
    if err != nil {
        return err
    }


    gateway.start()

    cmn.TrapSignal(func() {})

    return nil
}

func newGateway() (*gateway, error) {

    var clientpath string
    var datadir string
    if testnetFlag {
        datadir = filepath.Join(datadirFlag, "testnet")
    } else {
        datadir = datadirFlag
    }

    clientpath = filepath.Join(datadir, ipcpathFlag)

    ethclient, err := ethclient.Dial(clientpath)
    if err != nil {
        return nil, err
    }
/*
    codemap := etgate.GetCodemap()
   
    queries := []etgate.Query{}

    var cons []etgate.Contract
    for _, con := range data {
        addr := con["address"].(string)
        if !common.IsHexAddress(addr) {
            return errors.New("Invalid address format")
        }
        code, exists := con["code"]
        if !exists {
            return errors.New("Invalid code format")
        }
        cons = append(cons, etgate.Contract {
            Address: common.HexToAddress(addr),
            Code: code.(string),
        })
    }

    for _, con := range cons {
        query := codemap[con.Code].Query(con.Address)
        queries = append(queries, query...)
    }
*/
    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return nil, err
    }
 /*   if len(mintkey.PrivKey) != 32 {
        return nil, errors.New("Tendermint keyfile is not secp256k1")
    }
*/

    priv, err := getSecp256k1Priv(mintkey.PrivKey)
    if err != nil {
        return nil, err
    }

    ecdsa, err := ecrypto.ToECDSA(priv[:])
    fmt.Printf("%+v\n%+v\n", ecdsa, priv)
    if err != nil {
        return nil, err
    }

    fmt.Printf("%+v\n", ecrypto.PubkeyToAddress(ecdsa.PublicKey).Hex())

//    pub_, _ := getSecp256k1Pub(mintkey.PubKey)
//    _, pub := secp256k1.PrivKeyFromBytes(secp256k1.S256(), priv[:])
    //ecrypto.PubkeyToAddress(ecdsa.PublicKey)
/*    fmt.Printf("ecrypto pubkey: %v\n", ecdsa.PublicKey)
    fmt.Printf("crypto pubkey:    %v\n", pub)*/

    ethauth := bind.NewKeyedTransactor(ecdsa)

    ethauth.GasLimit = big.NewInt(4700000)

    /*   data, err := getConsfile(consfile)
    if err != nil {
        return err
    }    
   */ 
    return &gateway{
        ethclient: ethclient, 
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        ethauth: ethauth, 
        mintkey: mintkey,
//        query: queries,
    }, nil
}

func getSecp256k1Priv(priv crypto.PrivKey) (crypto.PrivKeySecp256k1, error) {
    switch inner := priv.Unwrap().(type) {
    case crypto.PrivKeySecp256k1:
        return inner, nil
    default:
        return crypto.PrivKeySecp256k1{}, errors.New("PrivKey is not secp256k1")
    }
}

func getSecp256k1Pub(pub crypto.PubKey) (crypto.PubKeySecp256k1, error) {
    switch inner := pub.Unwrap().(type) {
    case crypto.PubKeySecp256k1:
        return inner, nil
    default:
        return crypto.PubKeySecp256k1{}, errors.New("PubKey is not secp256k1")
    }
}

func (g *gateway) start() {
    go g.mintloop()
    go g.ethloop()
}

func (g *gateway) mintloop() {
    // Get last submitted withdrawal's sequence from ethereum    

    for {
        status, err := g.mintclient.Status()
        if err != nil {
            fmt.Printf("Failed to get status: \"%s\"\n", err)
            continue
        }
        height := status.LatestBlockHeight
        
        _, err = g.mintclient.Commit(height)
        if err != nil {
            fmt.Printf("Failed to get commit: \"%s\"\n", err)
            continue
        }

        key := fmt.Sprintf("etgate,withdraw,%s", /*change it later*/"etgate-chain")
        query, err := commands.QueryWithClient(g.mintclient, []byte(key))
        if err != nil {
            fmt.Printf("Failed to query last withdrawal: \"%s\"\n", err)
            continue
        }
        if len(query.Value) == 0 {
            continue
        }

        // submit withdrawals
    }
}

func (g *gateway) ethloop() {
    heads := make(chan *types.Header) 
    headsub, err := g.ethclient.SubscribeNewHead(context.Background(), heads)
    if err != nil {
        panic("Failed to subscribe to new headers")
    }  

    defer headsub.Unsubscribe()

    for {
        select {
        case head := <-heads:
            // Check if there is no submitted headers
            // TODO: remove this code at production phase
            if g.recentHeader() == nil {
                header, err := rlp.EncodeToBytes(head)
                if err != nil {
                    fmt.Printf("Failed to encode header: \"%s\"\n", err)
                    continue
                }
                updateTx := etgate.ETGateUpdateChainTx {
                    Headers: [][]byte{ header },
                }
                if err := g.appTx(updateTx); err != nil {
                    fmt.Printf("Error sending updateTx: \"%s\"\n", err)
                    continue
                }
                continue
            }

            // Check if the header already exists
            key := fmt.Sprintf("etgate,blockchain,buffer,%v", head.Hash().Hex())
            query, err := commands.QueryWithClient(g.mintclient, []byte(key))
            if err != nil {
                fmt.Printf("Failed to query: \"%s\"", err)
                continue
            }          
            
            
            if len(query.Value) != 0 {
                fmt.Printf("Header already submitted: %v. Skipping.\n", head.Number)
                g.post(g.recentHeader())
                continue
            }

            // Check the the header's parent exists
            // If it does, updateTx only the current header.
            key = fmt.Sprintf("etgate,blockchain,buffer,%v", head.ParentHash.Hex())
            query, err = commands.QueryWithClient(g.mintclient, []byte(key))
            if err != nil {
                fmt.Printf("Failed to query: \"%s\"", err)
                continue
            }
   
            if len(query.Value) != 0 {    
                header, err := rlp.EncodeToBytes(head)
                if err != nil {
                    fmt.Printf("Failed to encode header: \"%s\".\n", err)
                    continue
                }
                updateTx := etgate.ETGateUpdateChainTx {
                    Headers: [][]byte{ header },
                }
                if err := g.appTx(updateTx); err != nil {
                    fmt.Printf("Error sending updateTx: \"%s\". Skipping.\n", err)
                    continue
                }
                fmt.Printf("Submitted header: %v\n", head.Number)
                g.post(g.recentHeader())
                continue
            }

            // Sync from recent header to current header
            recent := g.recentHeader()
            g.sync(recent, head.Number)
        }
    }
}

func (g *gateway) recentHeader() *big.Int { // in mintchain
    key := fmt.Sprintf("etgate,blockchain,recent")
    query, err := commands.QueryWithClient(g.mintclient, []byte(key))
    
    if err != nil {
        panic("Error querying recent chain header")
    }

    recent := new(big.Int)
    if len(query.Value) == 0 {
        recent = nil 
    } else {
        var recentNumber uint64
        if err = wire.ReadBinaryBytes(query.Value, &recentNumber); err != nil {
            panic(err)
        }/*
        var recentInt int64
        recentInt, err = strconv.ParseInt(recentNumber, 10, 64)
        if err != nil {
            panic(err)
        }*/
        recent.SetUint64(recentNumber)
    }

    return recent
}

// sync does not post logs
func (g *gateway) sync(recent *big.Int, to *big.Int) { // (recent, to]
    temp := big.NewInt(0) // or whatever

    var err error

    for recent.Cmp(to) == -1 {
        var headers [][]byte
        if recent.Cmp(temp.Sub(to, chunksize)) == -1 {
            headers, err = g.getHeaders(temp.Add(recent, one), recent.Add(recent, chunksize))
            if err != nil {
                panic(err)
            }
        } else {
            headers, err = g.getHeaders(temp.Add(recent, one), recent.Add(recent, one))
            if err != nil {
                panic(err)
            }
        }

        updateTx := etgate.ETGateUpdateChainTx {
            Headers: headers,
        }
        if err = g.appTx(updateTx); err != nil {
            fmt.Printf("Error sending updateTx: \"%s\". Skipping.\n", err)
            continue
        }

        fmt.Printf("Submitted headers: %v\n", recent)
    }
}

func (g *gateway) getHeaders(from *big.Int, to *big.Int) ([][]byte, error) { // [from, to]
    res := [][]byte{}

    for i := from; i.Cmp(to) != 1; i.Add(i, one) {
        header, err := g.ethclient.HeaderByNumber(context.Background(), i)
        if err != nil {
            fmt.Printf("Error retrieving headers: \"%s\".\n", err)
            return nil, err
        }

        key := fmt.Sprintf("etgate,blockchain,buffer,%s", header.Hash().Hex())
        query, err := commands.QueryWithClient(g.mintclient, []byte(key))
        if err != nil {
            fmt.Printf("Error querying header: \"%s\".\n", err)
            return nil, err
        }
        if len(query.Value) != 0 {
            fmt.Printf("Header already submitted: %v. Skipping.\n", i)
            continue
        }

        headerb, err := rlp.EncodeToBytes(header)
        if err != nil {
            fmt.Printf("Error querying header: \"%s\".\n", err)
            return nil, err
        }

        res = append(res, headerb)
    }
    return res, nil
}

func (g *gateway) post(blocknumber *big.Int) {
    header, err := g.ethclient.HeaderByNumber(context.Background(), blocknumber)
    if err != nil {
        fmt.Printf("Error retrieving header: \"%s\". Skipping.\n", err)
        return
    }
    fmt.Printf("%v\n", header.Hash().Hex())
    
        
    depositQuery := ethereum.FilterQuery {
        Addresses: []common.Address{common.HexToAddress(addressFlag)},
        Topics: [][]common.Hash{{depositABI.Events["Deposit"].Id()}},
    }
    for _, q := range /*g.query*/ []ethereum.FilterQuery{depositQuery} {
        q.FromBlock = blocknumber
        q.ToBlock = blocknumber
        logs, err := g.ethclient.FilterLogs(context.Background(), q)
        if err != nil {
            fmt.Printf("Error retrieving logs: \"%s\". Skipping.\n", err)
            continue
        }

        for _, log := range logs {
            fmt.Printf("%v\n", log.BlockHash.Hex())
            proof, err := g.newLogProof(log)
            if err != nil {
                fmt.Printf("Error generating logproof: \"%s\". Skipping.\n", err)
                continue
            }
            logt, err := proof.Log()
            if err != nil {
                fmt.Printf("aa: %s", err)
                continue
            }
            fmt.Printf("%+v\n", logt)
            if !proof.IsValid(header.ReceiptHash) {
                fmt.Println("Invalid logproof. Skipping.")
                fmt.Printf("%+v\n", log)
                continue
            }
            
            if err != nil {
                panic(fmt.Sprintf("Error getting sequence: \"%s\".\n", err))
            }   

            postTx := etgate.ETGateDepositTx {
                proof,
            }
            if err := g.appTx(postTx); err != nil {
                fmt.Printf("Error sending postTx: \"%s\". Skipping.\n", err)
                continue
            }
            fmt.Println("Submitted log")
        }
    }
}

func (g *gateway) newLogProof(log types.Log) (etgate.LogProof, error) {
    var receipts []*types.Receipt

    count, err := g.ethclient.TransactionCount(context.Background(), log.BlockHash)
    if err != nil {
        return etgate.LogProof{}, err
    }

    for i := 0; i < int(count); i++ {
        tx, err := g.ethclient.TransactionInBlock(context.Background(), log.BlockHash, uint(i))
        if err != nil {
            return etgate.LogProof{}, err
        }

        receipt, err := g.ethclient.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
            return etgate.LogProof{}, err
        }

        receipts = append(receipts, receipt)
    }

    return etgate.NewLogProof(receipts, log.TxIndex, log.Index, log.BlockNumber)
}

func (g *gateway) appTx(etgateTx etgate.ETGateTx) error {
	return commands.AppTx(g.mintclient, g.mintkey, etgateTx, chainIDFlag)
}
