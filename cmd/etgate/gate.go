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

type gateway struct {
    ethclient *ethclient.Client
    mintclient *mintclient.HTTP
    mintkey *commands.Key
    query ethereum.FilterQuery
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
            Address: addr,
            Code: code,
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
    
    var addrs []common.Address
    for _, con := range data {
        addr := con["address"].(string)
        if !common.IsHexAddress(addr) {
            return errors.New("Invalid address format")
        }
        addrs = append(addrs, common.HexToAddress(addr))
    }

    gateway, err := newGateway(filepath.Join(datadir, ipcpathFlag), addrs)
    if err != nil {
        return err
    }

    gateway.start()

    cmn.TrapSignal(func() {})

    return nil
}

func newGateway(ipc string, addrs []common.Address) (*gateway, error) {
    ethclient, err := ethclient.Dial(ipc)
    if err != nil {
        return nil, err
    }

    fmt.Printf("%+v\n", addrs)

    query := ethereum.FilterQuery {
//        Addresses: addrs,
    }

    fmt.Printf("%+v\n", query)

    return &gateway{
        ethclient: ethclient, 
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        query: query,
    }, nil
}

func (g *gateway) start() {
    go g.loop()    
}

func (g *gateway) loop() {
    logs := make(chan types.Log)
    heads := make(chan *types.Header)

    logsub, err := g.ethclient.SubscribeFilterLogs(context.Background(), g.query, logs)
    if err != nil {
        panic("Failed to subscribe to log events")
    }

    headsub, err := g.ethclient.SubscribeNewHead(context.Background(), heads)
    if err != nil {
        panic("Failed to subscribe to new headers")
    }  

    defer logsub.Unsubscribe()
    defer headsub.Unsubscribe()

    for {
        select {
        case head := <-heads:
            updateTx := etgate.ETGateUpdateChainTx {
                Header: *head,        
            }
            if err := g.appTx(updateTx); err != nil {
                fmt.Printf("Error in UpdataChain: ", err)
                continue
            }
            fmt.Printf("Submitted header: %d\n", head.Number)

        case log := <-logs:
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
            
            postTx := etgate.ETGatePacketPostTx {
                Proof: *proof,
            }

            if err := g.appTx(postTx); err != nil {
                fmt.Println("Error submitting logs: ", err)
                continue
            }

            fmt.Printf("Submitted log\n")
        }
    }
}

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
/*    if testnetFlag {
        tx.Input.Signature = g.mintkey.Sign(tx.SignBytes("ethereum-ropsten"))
    } else {
        tx.Input.Signature = g.mintkey.Sign(tx.SignBytes("ethereum-mainnet"))
    }
*/
    fmt.Println(chainIDFlag)
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
