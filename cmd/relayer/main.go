package main

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"
    "context"
    "bytes"

    "gopkg.in/urfave/cli.v1"

    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/trie"
//    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/rlp"
    // "github.com/ethereum/go-ethereum"
    // "github.com/ethereum/go-ethereum"
    // "github.com/ethereum/go-ethereum"

    cmn "github.com/tendermint/tmlibs/common"

    "../.."
)

type Gateway struct {
    client *ethclient.Client
    query ethereum.FilterQuery
}

var RETRY int = 10

func main() {
    app := cli.NewApp()
    app.Name = "etgateway-relayer"
    app.Usage = "Relay ethereum logs to tendermint blockchain"
    app.Flags = []cli.Flag {
        cli.BoolFlag {
            Name: "testnet",
            Usage: "Ropsten network: pre-configured test network",
        },
        cli.StringFlag {
            Name: "datadir",
            Value: filepath.Join(os.Getenv("HOME"), ".ethereum"),
            Usage: "Data directory for the databases and keystore",
        },
        cli.StringFlag {
            Name: "ipcpath",
            Value: "geth.ipc",
            Usage: "Filename for IPC socket/pipe within the datadir",
        },
    }
    app.Action = action

    if err := app.Run(os.Args); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func action(c *cli.Context) error {
    var datadir string

    if c.NArg() < 1 {
        return errors.New("Usage: etgateway [--testnet] [--datadir ~/.ethereum] [--ipcpath geth.ipc] address")
    }
    addr := c.Args()[0]
    
    if c.Bool("testnet") {
        datadir = filepath.Join(c.String("datadir"), "testnet")
    } else {
        datadir = c.String("datadir")
    }

    if !common.IsHexAddress(addr) {
        return errors.New("Invalid address format")
    }

    gateway, err := newGateway(datadir, c.String("ipcpath"), common.HexToAddress(addr))
    if err != nil {
        return err
    }

    gateway.start()

    cmn.TrapSignal(func() {})

    return nil
}

func newGateway(datadir, ipcpath string, addr common.Address) (*Gateway, error) {
    client, err := ethclient.Dial(filepath.Join(datadir, ipcpath))
    if err != nil {
        return nil, err
    }

    fmt.Printf("%+v\n", addr)

    query := ethereum.FilterQuery {
//        Addresses: []common.Address{addr},
    }

    fmt.Printf("%+v\n", query)

    return &Gateway{client: client, query: query,}, nil
}

func (g *Gateway) start() {
    go g.loop()    
}

func (g *Gateway) loop() {
    logs := make(chan types.Log)

    sub, err := g.client.SubscribeFilterLogs(context.Background(), g.query, logs)
    if err != nil {
        fmt.Println("Failed to subscribe to log events")
    }
    defer sub.Unsubscribe()
    
    var i int

    for {
        select {
            case log := <-logs:
            for i = 0; i < RETRY; i++ {
                proof, err := g.NewLogProof(log)
                if err != nil {
                    fmt.Println("Error in NewLogProof, retrying...: ", err)
                    continue
                }
        
                if !proof.IsValid() {
                    fmt.Println("Invalid log proof generation, retrying...")
                    continue
                }

                fmt.Println("Succeed")
                break
            }
            if i == RETRY {
                fmt.Printf("%+v\n", log)
                header, _ := g.client.HeaderByHash(context.Background(), log.BlockHash)
                fmt.Printf("%+v\n", header)
//              panic("Exceed retry limit")
            }
        }
    }
}

func (g *Gateway) NewLogProof(log types.Log) (*etgate.LogProof, error) {
    var logReceipt *types.Receipt
    keybuf := new(bytes.Buffer)
    trie := new(trie.Trie)

    count, err := g.client.TransactionCount(context.Background(), log.BlockHash)
    if err != nil {
        return nil, err
    }


    for i := 0; i < int(count); i++ {
        tx, err := g.client.TransactionInBlock(context.Background(), log.BlockHash, uint(i))
        if err != nil {
            return nil, err
        }

        receipt, err := g.client.TransactionReceipt(context.Background(), tx.Hash())
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


    header, err := g.client.HeaderByHash(context.Background(), log.BlockHash)
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


