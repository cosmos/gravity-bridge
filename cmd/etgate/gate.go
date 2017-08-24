package main

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"
    "context"
    "encoding/json"
    "io/ioutil"
    "math/big"
    "bytes"

//    "golang.org/x/crypto/ripemd160"

    "github.com/spf13/cobra"

    "github.com/ethereum/go-ethereum/ethclient"
//    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
 
    abci "github.com/tendermint/abci/types"
    cmn "github.com/tendermint/tmlibs/common"
    basecmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
    bctypes "github.com/tendermint/basecoin/types"
    mintclient "github.com/tendermint/tendermint/rpc/client"
    tmtypes "github.com/tendermint/tendermint/types"
 //   "github.com/tendermint/tmlibs/merkle"

    "github.com/tendermint/go-wire"

    "../../plugins/etgate"
)

var ( 
    one = big.NewInt(1)
    chunksize = big.NewInt(16)
)


type gateway struct {
    ethclient *ethclient.Client
    mintclient *mintclient.HTTP
    mintkey *basecmd.Key
    query []etgate.Query
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
)

func init() {
    flags := []basecmd.Flag2Register {
        {&testnetFlag, "testnet", false, "Ropsten network: pre-configured test network"},
        {&nodeaddrFlag, "nodeaddr", "tcp://localhost:46657", "Node address for tendermint chain"},
        {&chainIDFlag, "chain-id", "etgate-chain", "Chain ID"},
    }

    basecmd.RegisterPersistentFlags(GateCmd, flags)

    startFlags := []basecmd.Flag2Register {
        {&datadirFlag, "datadir", filepath.Join(os.Getenv("HOME"), ".ethereum"), "Data directory for the databases and keystore"},
        {&ipcpathFlag, "ipcpath", "geth.ipc", "Filename for IPC socket/pipe within the datadir"},
    }

    basecmd.RegisterFlags(GateStartCmd, startFlags)

    GateCmd.AddCommand(GateStartCmd)
    GateCmd.AddCommand(GateInitCmd)
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
    if len(args) < 1 {
        return errors.New("Usage: etgate gate init [--testnet] contractsfile")
    }
    consfile := args[0]

    data, err := getConsfile(consfile)
    if err != nil {
        return err
    }

    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return err
    }

    g := gateway {
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        mintkey: mintkey,
    }

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
    return nil
}

func gateStartCmd(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
        return errors.New("Usage: etgate gate start [--testnet] [--datadir ~/.ethereum] [--ipcpath geth.ipc] [--rpcpath localhost:1234] contractsfile")
    } 
    consfile := args[0]

    var clientpath string
    var datadir string
    if testnetFlag {
        datadir = filepath.Join(datadirFlag, "testnet")
    } else {
        datadir = datadirFlag
    }

    clientpath = filepath.Join(datadir, ipcpathFlag)

    data, err := getConsfile(consfile)
    if err != nil {
        return err
    }   
 
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


    gateway, err := newGateway(clientpath, cons)
    if err != nil {
        return err
    }


    gateway.start()

    cmn.TrapSignal(func() {})

    return nil
}

func newGateway(ipc string, cons []etgate.Contract) (*gateway, error) {
    ethclient, err := ethclient.Dial(ipc)
    if err != nil {
        return nil, err
    }

    codemap := etgate.GetCodemap()
   
    queries := []etgate.Query{}

    for _, con := range cons {
        query := codemap[con.Code].Query(con.Address)
        queries = append(queries, query...)
    }

    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return nil, err
    }

    return &gateway{
        ethclient: ethclient, 
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        mintkey: mintkey,
        query: queries,
    }, nil
}

func (g *gateway) start() {
    go g.mintloop()
    go g.ethloop()
}

func (g *gateway) mintloop() {
    for {
        status, err := g.mintclient.Status()
        if err != nil {
            fmt.Printf("Failed to get status: \"%s\"\n", err)
            continue
        }
        height := status.LatestBlockHeight
        
        commit, err := g.mintclient.Commit(height)
        if err != nil {
            fmt.Printf("Failed to get commit: \"%s\"\n", err)
            continue
        }
        /*
        for _, pc := range commit.Commit.Precommits {
              fmt.Printf("%+v\n", tmtypes.SignBytes("etgate-chain", pc))
        }*/

        fmt.Printf("%+v\n", commit.Header.Hash())
/*
        proofs := merkle.SimpleProofsFromHashables([]merkle.Hashable{commit.Header})
        fmt.Printf("%+v\n", proofs[0])

        */
/*
        header := commit.Header
        headerMap := map[string]interface{}{
            "ChainID": header.ChainID,
            "Height": header.Height,
            "Time": header.Time,
            "NumTxs": header.NumTxs,
            "LastBlockID": header.LastBlockID,
            "LastCommit": header.LastCommitHash,
            "Data": header.DataHash,
            "Validators": header.ValidatorsHash,
            "App": header.AppHash,
        }
        kpPairsH := merkle.MakeSortedKVPairs(headerMap)
        fmt.Printf("%+v\n", kpPairsH)
*/
        for _, data := range []uint{4, 16, 256, 65536, 100000, 4294967296, 18446744073709551615} {
            buf := new(bytes.Buffer)
            n, err := int(0), error(nil)
            wire.WriteUvarint(data, buf, &n, &err)
            fmt.Printf("%d:\t%+v\n", data*2+1, buf.Bytes())
        }

        for {
            status, _ := g.mintclient.Status()
            if status.LatestBlockHeight > height {
                break
            }
        }
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
            query, err := queryWithClient(g.mintclient, []byte(key))
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
            query, err = queryWithClient(g.mintclient, []byte(key))
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
    query, err := queryWithClient(g.mintclient, []byte(key))
    
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
        query, err := queryWithClient(g.mintclient, []byte(key))
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

    for _, q := range g.query {
        q.FromBlock = blocknumber
        q.ToBlock = blocknumber
        logs, err := g.ethclient.FilterLogs(context.Background(), q.FilterQuery)
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

            postTx := etgate.ETGatePacketPostTx {
                Name: q.Name,
                Proof: proof,
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
    acc, err := getAccWithClient(g.mintclient, g.mintkey.Address[:])
    if err != nil {
        return err
    }
    sequence := acc.Sequence + 1

    data := []byte(wire.BinaryBytes(struct {
        etgate.ETGateTx `json:"unwrap"`
    }{etgateTx}))

    smallCoins := bctypes.Coin{Denom: "mycoin", Amount: 1}

    input := bctypes.NewTxInput(g.mintkey.PubKey, bctypes.Coins{smallCoins}, sequence)
    tx := &bctypes.AppTx {
        Gas: 0,
        Fee: smallCoins,
        Name: "ETGATE",
        Input: input,
        Data: data,
    }
    tx.Input.Signature = g.mintkey.Sign(tx.SignBytes(chainIDFlag))
    txBytes := []byte(wire.BinaryBytes(struct {
        bctypes.Tx `json:"unwrap"`
    }{tx}))

    data, log, err := broadcastTxWithClient(g.mintclient, txBytes)
    if err != nil {
        return err
    }

    _, _ = data, log
    return nil
}

func getAccWithClient(httpClient *mintclient.HTTP, address []byte) (*bctypes.Account, error) {

	key := bctypes.AccountKey(address)
	response, err := queryWithClient(httpClient, key)
	if err != nil {
		return nil, err
	}

	accountBytes := response.Value

	if len(accountBytes) == 0 {
		return nil, fmt.Errorf("Account bytes are empty for address: %X ", address) //never stack trace
	}

	var acc *bctypes.Account
	err = wire.ReadBinaryBytes(accountBytes, &acc)
	if err != nil {
		return nil, fmt.Errorf("Error reading account %X error: %v",
			accountBytes, err.Error())
	}

	return acc, nil
}

func queryWithClient(httpClient *mintclient.HTTP, key []byte) (*abci.ResultQuery, error) {
	res, err := httpClient.ABCIQuery("/key", key, true)
	if err != nil {
		return nil, fmt.Errorf("Error calling /abci_query: %v", err)
	}
	if !res.Code.IsOK() {
		return nil, fmt.Errorf("Query got non-zero exit code: %v. %s", res.Code, res.Log)
	}
	return res.ResultQuery, nil
}

func broadcastTxWithClient(httpClient *mintclient.HTTP, tx tmtypes.Tx) ([]byte, string, error) {
    res, err := httpClient.BroadcastTxCommit(tx)
    if err != nil {
        return nil, "", fmt.Errorf("Error on broadcast tx: %v", err)
    }

    if !res.CheckTx.Code.IsOK() {
        r := res.CheckTx
        return nil, "", fmt.Errorf("BroadcastTxCommit got non-zero exit code: %v, %X; %s", r.Code, r.Data, r.Log)
   }

    if !res.DeliverTx.Code.IsOK() {
        r := res.DeliverTx
        return nil, "", fmt.Errorf("BroadcastTxCommit got non-zero exit code: %v, %X; %s", r.Code, r.Data, r.Log)
    }

    return res.DeliverTx.Data, res.DeliverTx.Log, nil
}
