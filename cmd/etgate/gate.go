package main

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"
    "context"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "math/big"
    "time"
    "strconv"

    "github.com/spf13/cobra"

    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/trie"
//    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/rlp"
 
    abci "github.com/tendermint/abci/types"
    cmn "github.com/tendermint/tmlibs/common"
    "github.com/tendermint/basecoin/cmd/basecoin/commands"
    bctypes "github.com/tendermint/basecoin/types"
    mintclient "github.com/tendermint/tendermint/rpc/client"
    tmtypes "github.com/tendermint/tendermint/types"

    "github.com/tendermint/go-wire"

    "../../plugins/etgate"
)

var ( 
    chunksize = big.NewInt(1024)
)

type gateway struct {
    ethclient *ethclient.Client
    mintclient *mintclient.HTTP
    mintkey *commands.Key
    query []ethereum.FilterQuery
    buffer map[uint64][]*etgate.LogProof
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
    flags := []commands.Flag2Register {
        {&testnetFlag, "testnet", false, "Ropsten network: pre-configured test network"},
        {&nodeaddrFlag, "nodeaddr", "tcp://localhost:46657", "Node address for tendermint chain"},
        {&chainIDFlag, "chain-id", "etgate-chain", "Chain ID"},
    }

    commands.RegisterPersistentFlags(GateCmd, flags)

    startFlags := []commands.Flag2Register {
        {&datadirFlag, "datadir", filepath.Join(os.Getenv("HOME"), ".ethereum"), "Data directory for the databases and keystore"},
        {&ipcpathFlag, "ipcpath", "geth.ipc", "Filename for IPC socket/pipe within the datadir"},
    }

    commands.RegisterFlags(GateStartCmd, startFlags)

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

    mintkey, err := commands.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
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
        return errors.New("Usage: etgate gate start [--testnet] [--datadir ~/.ethereum] [--ipcpath geth.ipc] contractsfile")
    } 
    consfile := args[0]

    var datadir string
    if testnetFlag {
        datadir = filepath.Join(datadirFlag, "testnet")
    } else {
        datadir = datadirFlag
    }

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

    gateway, err := newGateway(filepath.Join(datadir, ipcpathFlag), cons)
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
   
    queries := []ethereum.FilterQuery{}

    for _, con := range cons {
        query := codemap[con.Code].Query(con.Address)
        queries = append(queries, query...)
    }
    fmt.Printf("%v\n", queries)

    mintkey, err := commands.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return nil, err
    }

    return &gateway{
        ethclient: ethclient, 
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        mintkey: mintkey,
        query: queries,
        buffer: make(map[uint64][]*etgate.LogProof),
    }, nil
}

func (g *gateway) start() {
    go g.loop()
}
/*
func (g *gateway) logloop() {
    logs := [](chan types.Log){}
    for _, q := range g.query {
        log := make(chan types.Log)
        sub, err := g.ethclient.SubscribeFilterLogs(context.Background(), q, log)
        if err != nil {
            fmt.Println("Failed to subscribe log events")
            panic(err)
        }
        defer sub.Unsubscribe()
        logs = append(logs, log)
    }
    cases := make([]reflect.SelectCase, len(logs))

    for i, ch := range logs {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
    }

    for {
        _, data, ok := reflect.Select(cases)
        if !ok {
            panic("Channel closed unexpectedly")
        }
        log := data.Interface().(types.Log)
        proof, err := g.NewLogProof(log)
        if err != nil {
            fmt.Println("Error in NewLogProof: ", err)
            fmt.Printf("%+v\n", log)
            header, _ := g.ethclient.HeaderByHash(context.Background(), log.BlockHash)
            fmt.Printf("%+v\n", header)
            continue
        }
      
        if !proof.IsValid() {
            fmt.Println("Invalid log proof generation")
            fmt.Printf("%+v\n", log)
            header, _ := g.ethclient.HeaderByHash(context.Background(), log.BlockHash)
            fmt.Printf("%+v\n", header)
            continue
        }
        
        g.buffer[log.BlockNumber] = append(g.buffer[log.BlockNumber], proof)

        fmt.Println("Received log")

        
        postTx := etgate.ETGatePacketPostTx {
            Name: "",
            Proof: *proof,
        }

        if err := g.appTx(postTx); err != nil {
            fmt.Println("Error submitting logs: ", err)
            continue
        }

        fmt.Printf("Submitted log\n")
    } 
}
*/

func (g *gateway) loop() {
    heads := make(chan *types.Header) 
    headsub, err := g.ethclient.SubscribeNewHead(context.Background(), heads)
    if err != nil {
        panic("Failed to subscribe to new headers")
    }  

    defer headsub.Unsubscribe()

    for {
        select {
        case head := <-heads:
            fmt.Println("Received new header")
            // Check if the header already exists
            key := fmt.Sprintf("etgate,blockchain,buffer,%s", head.Hash().Str())
            query, err := queryWithClient(g.mintclient, []byte(key))
            if err != nil {
                fmt.Printf("Failed to query: \"%s\"", err)
                continue
            }
            
            if len(query.Value) != 0 {
                continue
            }

            // Check the the header's parent exists
            // If it does, updateTx current header.
            key = fmt.Sprintf("etgate,blockchain,buffer,%s", head.ParentHash.Str())
            query, err = queryWithClient(g.mintclient, []byte(key))
            if err != nil {
                fmt.Printf("Failed to query: \"%s\"", err)
                continue
            }

            if len(query.Value) != 0 {
                updateTx := etgate.ETGateUpdateChainTx {
                    Header: *head,
                }
                if err := g.appTx(updateTx); err != nil {
                    fmt.Printf("Error sending updateTx: \"%s\". Skipping.\n", err)
                    continue
                }
                fmt.Printf("Submitted header: %v\n", head.Number)
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

    var recent *big.Int
    if len(query.Value) == 0 {
        recent = big.NewInt(0)
    } else {
        wire.ReadBinaryBytes(query.Value, recent)
    }

    return recent
}

func (g *gateway) sync(recent *big.Int, to *big.Int) { // (recent, to]
    one := big.NewInt(1)
    recent.Add(recent, one)

    for recent.Cmp(to) == -1 {
        header, err := g.ethclient.HeaderByNumber(context.Background(), recent)
        if err != nil {
            fmt.Printf("Error retrieving headers: \"%s\". Retrying...\n", err)
            time.Sleep(time.Second)
            continue
        }

        key := fmt.Sprintf("etgate,blockchain,buffer,%s", header.Hash().Str())
        query, err := queryWithClient(g.mintclient, []byte(key))
        if err != nil {
            fmt.Printf("Error querying header: \"%s\", Retrying...\n", err)
            time.Sleep(time.Second)
            continue
        }
        if len(query.Value) != 0 {
            fmt.Printf("Header is already submitted. Skipping.")
            continue
        }

        updateTx := etgate.ETGateUpdateChainTx {
            Header: *header,
        }
        if err := g.appTx(updateTx); err != nil {
            fmt.Printf("Error sending updateTx: \"%s\". Skipping.\n", err)
            recent.Add(recent, one)
            continue
        }

        fmt.Printf("Submitted header: %v\n", recent)

        recent.Add(recent, one)
    }
}

func (g *gateway) submit(blocknumber *big.Int) error {
    header, err := g.ethclient.HeaderByNumber(context.Background(), blocknumber)
    if err != nil {
        return err
    }

    updateTx := etgate.ETGateUpdateChainTx {
        Header: *header,
    }
    if err := g.appTx(updateTx); err != nil {
        return err
    }

    g.post(blocknumber)

    return nil
}

func (g *gateway) post(blocknumber *big.Int) {
    for _, q := range g.query {
        logs, err := g.ethclient.FilterLogs(context.Background(), q)
        if err != nil {
            fmt.Printf("Error retrieving logs: \"%s\". Skipping.\n", err)
            continue
        }

        for _, log := range logs {
            proof, err := g.NewLogProof(log)
            if err != nil {
                fmt.Printf("Error generating logproof: \"%s\". Skipping.\n", err)
                continue
            }
            if !proof.IsValid() {
                fmt.Println("Invalid logproof. Skipping.")
                continue
            }
            
            seq, err := g.getSequence()
            if err != nil {
                panic(fmt.Sprintf("Error getting sequence: \"%s\".\n", err))
            }   

            postTx := etgate.ETGatePacketPostTx {
                Proof: *proof,
                Sequence: seq,
            }
            if err := g.appTx(postTx); err != nil {
                fmt.Printf("Error sending postTx: \"%s\". Skipping.\n", err)
                continue
            }
        }
    }
}

func (g *gateway) getSequence() (uint64, error) {
    key := fmt.Sprintf("etgate,contract")
    query, err := queryWithClient(g.mintclient, []byte(key))
    if err != nil {
        return 0, err
    }
    if len(query.Value) == 0 {
        return 0, nil
    }

    seq, err := strconv.ParseUint(string(query.Value), 10, 64)
    if err != nil {
        return 0, err
    }

    return seq, nil
}
/*
func (g *gateway) headsync(recent *big.Int, to *big.Int) { // after recent, including to
    recent.Add(recent, big.NewInt(1))
    for recent.Cmp(to) == -1 {
        temp := big.NewInt(0)
        headers := make([]types.Header, chunksize.Int64())
        for i := big.NewInt(0); i.Cmp(chunksize) == -1 && temp.Add(recent, i).Cmp(to) != 1; i.Add(i, big.NewInt(1)) {
            header, err := g.ethclient.HeaderByNumber(context.Background(), temp.Add(i, recent))
            if err != nil {
                panic(err)
            }
            headers[i.Int64()] = *header
        }
        updateTx := etgate.ETGateUpdateChainTx {
            Header: headers,
        }
        if err := g.appTx(updateTx); err != nil {
            fmt.Println("Error in headsync")
            panic(err)
        }
        fmt.Printf("Submitted headers: %v\n", temp.Add(recent, chunksize))
        recent = g.recentHeader()
    }
}

func (g *gateway) headloop() {
    heads := make(chan *types.Header) 
    headsub, err := g.ethclient.SubscribeNewHead(context.Background(), heads)
    if err != nil {
        panic("Failed to subscribe to new headers")
    }  

    defer headsub.Unsubscribe()

    for {
        select {
            case head := <-heads:
            before := big.NewInt(1)
            before.Sub(head.Number, before)
            recent := g.recentHeader()
            if recent.Cmp(before) == 0 { // recent submitted header == right before
                updateTx := etgate.ETGateUpdateChainTx {
                    Headers: []types.Header{ *head },
                }

                if err := g.appTx(updateTx); err != nil {
                    fmt.Printf("Error in UpdateChain: ", err)
                    continue
                }
                fmt.Printf("Submitted header: %v\n", head.Number)
            } else if recent.Cmp(before) == -1 { // recent submitted header < right before
                g.headsync(recent, head.Number)
            } // nothing to do if recent submitted header >= received
            recent = g.recentHeader()

            var packets []etgate.Packet

            for q := range g.query {
                q.FromBlock = recent
                q.ToBlock = recent
                logs, err := g.ethclient.FilterLogs(context.Background(), q)
                if err != nil {
                    fmt.Println("Error querying filterquery")
                    fmt.Println(err)
                    continue
                }

                for log := range logs {
                    proof := g.NewLogProof(log)
                    key := fmt.Sprintf("etgate,contract,")
                    queryWithClient(g.mintclient, []byte(key))
                }
            }
        }
    }
}*/

func (g *gateway) NewLogProof(log types.Log) (*etgate.LogProof, error) {
    var logReceipt *types.Receipt
    keybuf := new(bytes.Buffer)
    trie := new(trie.Trie)

    count, err := g.ethclient.TransactionCount(context.Background(), log.BlockHash)
    if err != nil {
        return nil, err
    }


    for i := 0; i < int(count); i++ {
        tx, err := g.ethclient.TransactionInBlock(context.Background(), log.BlockHash, uint(i))
        if err != nil {
            return nil, err
        }

        receipt, err := g.ethclient.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
            return nil, err
        }

        keybuf.Reset()
        rlp.Encode(keybuf, uint(i))
        bytes, err := rlp.EncodeToBytes(receipt)
        if err != nil {
            return nil, err
        }
        trie.Update(keybuf.Bytes(), bytes)

        if log.TxIndex == uint(i) {
            logReceipt = receipt
        }
    }

    if logReceipt == nil {
        return nil, errors.New("Error in NewLogProof")   
    }

    keybuf.Reset()
    rlp.Encode(keybuf, log.TxIndex)


    header, err := g.ethclient.HeaderByHash(context.Background(), log.BlockHash)
    if err != nil {
        return nil, err
    }

    return &etgate.LogProof {
        ReceiptHash: header.ReceiptHash,
        Receipt: logReceipt,
        TxIndex: log.TxIndex,
        Index: log.Index,
        Proof: trie.Prove(keybuf.Bytes()),  
    }, nil
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
