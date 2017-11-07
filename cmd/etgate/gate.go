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
    "time"


    "github.com/spf13/cobra"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    ecrypto "github.com/ethereum/go-ethereum/crypto"
 
    cmn "github.com/tendermint/tmlibs/common"
    basecmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
    mintclient "github.com/tendermint/tendermint/rpc/client"
    tmtypes "github.com/tendermint/tendermint/types"

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

var GateGenValidatorCmd = &cobra.Command {
    Use: "genval",
    Short: "Generate secp256k1 validator",
    RunE: gateGenValidatorCmd,
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
    GateCmd.AddCommand(GateGenValidatorCmd)

    var err error
    depositABI, err = abi.JSON(strings.NewReader(contracts.ETGateABI))
    if err != nil {
        panic(err)
    }

}
/*
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
*/

func gateGenValidatorCmd(cmd *cobra.Command, args []string) error {
    privKey := crypto.GenPrivKeySecp256k1()
    pubKey := privKey.Wrap().PubKey().Unwrap()
    var addr common.Address
    copy(addr[:], pubKey.Address())
    fmt.Printf("Priv:\t%v\nPub:\t%v\nAddr:\t%v\n", strings.ToUpper(hex.EncodeToString(privKey[:])), pubKey.KeyString(), strings.ToUpper(addr.Hex()[2:]))

    return nil
}

func gateInitCmd(cmd *cobra.Command, args []string) error {

    g, err := newGateway()
    if err != nil {
        return err
    }

    g.ethauth.GasLimit = big.NewInt(4700000)

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

    address, _, _, err := contracts.DeployETGate(g.ethauth, g.ethclient, "etgate-chain", validatorsBytes, votingPowers)
    if err != nil {
        return err
    }

    fmt.Printf("ETGate contract is deployed on %s\n", address.Hex())

    return nil
}

func gateStartCmd(cmd *cobra.Command, args []string) error {
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

    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return nil, err
    }

    priv, err := getSecp256k1Priv(mintkey.PrivKey)
    if err != nil {
        return nil, err
    }

    ecdsa, err := ecrypto.ToECDSA(priv[:])
    if err != nil {
        return nil, err
    }

    fmt.Printf("Using Ethereum address %+v\n", ecrypto.PubkeyToAddress(ecdsa.PublicKey).Hex())

    ethauth := bind.NewKeyedTransactor(ecdsa)
    ethauth.GasPrice = big.NewInt(20000000000)

    return &gateway{
        ethclient: ethclient, 
        mintclient: mintclient.NewHTTP(nodeaddrFlag, "/websocket"),
        ethauth: ethauth, 
        mintkey: mintkey,
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
    contract, err := contracts.NewETGate(common.HexToAddress(addressFlag), g.ethclient)
    if err != nil {
        panic(err)
    }

 /*
    isVal, err := contract.SenderIsValidator(nil)
    if err != nil {
        panic(err)
    }
    if !isVal {
        fmt.Println("This key is not a validator.")
        return
    }
*/
    for {
        time.Sleep(5 * time.Second)
/*
        // update headers mint->eth
        // uncomment this part when started using iavl directly

        chainState, err := contract.ChainState(nil)
        if err != nil {
            panic(err)
        }


        lastHeight := chainState.LastBlockHeight

        status, err := g.mintclient.Status()
        if err != nil {
            fmt.Printf("Failed to get status: \"%s\"\n", err)
            continue
        }
        height := status.LatestBlockHeight
        fmt.Printf("Current height: %d\n", height)

        updated, err := contract.GetUpdated(nil, big.NewInt(int64(height)))
        if err != nil {
            fmt.Printf("Failed to get updated: \"%s\"\n", err)
            continue
        }
        if updated.Cmp(big.NewInt(0)) != 0 {
            // compare hash in production phase
            fmt.Printf("Header already submitted\n")
            continue
        }


        fmt.Printf("debug: %+v, %+v\n", lastHeight.Uint64(), height)

        bigHeight := big.NewInt(int64(height))
        for lastHeight.Add(lastHeight, one); lastHeight.Cmp(bigHeight) != 1; lastHeight.Add(lastHeight, one) {
            commit, err := g.mintclient.Commit(int(lastHeight.Uint64()))
            if err != nil {
                fmt.Printf("Failed to get commit: \"%s\"\n", err)
                break
            }
    
            header := commit.Header

            var timeHashArr [20]byte
            copy(timeHashArr[:], merkle.KVPair{Key: "Time", Value: header.Time}.Hash())
            var blockIDHashArr [20]byte
            copy(blockIDHashArr[:], header.LastBlockID.Hash)
            var partsHeaderHashArr [20]byte
            copy(partsHeaderHashArr[:], header.LastBlockID.PartsHeader.Hash)
            var lastCommitHashArr [20]byte
            copy(lastCommitHashArr[:], header.LastCommitHash)
            var dataHashArr [20]byte
            copy(dataHashArr[:], header.DataHash)
            var validatorsHashArr [20]byte
            copy(validatorsHashArr[:], header.ValidatorsHash)
            var appHashArr [20]byte
            copy(appHashArr[:], header.AppHash)

            tx, err := contract.Update(
                g.ethauth, 
                "etgate-chain", 
                big.NewInt(int64(header.Height)), 
                timeHashArr,
                big.NewInt(int64(header.NumTxs)),
                blockIDHashArr,
                big.NewInt(int64(header.LastBlockID.PartsHeader.Total)),
                partsHeaderHashArr,
                lastCommitHashArr,
                dataHashArr,
                validatorsHashArr,
                appHashArr,
            )
            if err != nil {
                fmt.Printf("Failed to update header: \"%s\", Waiting...\n", err)
                time.Sleep(30 * time.Second)
                break
            }
            fmt.Printf("Submitted header: %d\nTx hash: %v\n", header.Height, tx.Hash().Hex())
            time.Sleep(5 * time.Second)
        }
*/
        // relay withdraw

        key := fmt.Sprintf("etgate,withdraw,%s", "etgate-chain")
        query, err := commands.QueryWithClient(g.mintclient, []byte(key))
        if err != nil {
            fmt.Printf("Failed to query last withdrawal: \"%s\"\n", err)
            continue
        }

        var sequence uint64
        if len(query.Value) == 0 {
            sequence = 0
        } else {
            if err = wire.ReadBinaryBytes(query.Value, &sequence); err != nil {
                fmt.Printf("Error reading sequence from query: \"%s\"\n", err)
                continue
            }
        }
        
        lastWithdraw, err := contract.LastWithdraw(nil)
        if err != nil {
            fmt.Printf("Failed to get lastWithdraw: \"%s\", Waiting...\n", err)
            continue
        }

        fmt.Printf("debug: %d, %d\n", lastWithdraw.Uint64()+1, sequence)
        for i := lastWithdraw.Uint64()+1; i <= sequence; i++ {
            key = fmt.Sprintf("etgate,withdraw,%s,%v", "etgate-chain", i)
            query, err = commands.QueryWithClient(g.mintclient, []byte(key))
            if err != nil {
                fmt.Printf("Failed to query withdrawal: \"%s\", Waiting...\n", err)
                break
            }
            if len(query.Value) == 0 {
                fmt.Printf("Failed to query withdrawal: len(query.Value) == 0, Waiting...")
                break
            }
            var w etgate.ETGateWithdrawTx
            if err = wire.ReadBinaryBytes(query.Value, &w); err != nil {
                fmt.Printf("Failed to read withdraw tx: \"%s\", Waiting...\n", err)
                break
            }
            
            able, err := contract.Withdrawable(
                nil, 
                big.NewInt(int64(w.Height)), 
                common.BytesToAddress(w.To[:]), 
                w.Value, 
                common.HexToAddress(etgate.DecodeToken(w.Token)), 
                []byte(w.ChainID), 
                big.NewInt(int64(i)),
            ) 
            
            if err != nil {
                fmt.Printf("Error calling withdrawable: \"%s\", Waiting...", err)
            }

            if !able {
                fmt.Printf("#%d is not withdrawable, Waiting...\n", i)
                break
            }

            tx, err := contract.Withdraw(
                g.ethauth, 
                big.NewInt(int64(w.Height)), 
                common.BytesToAddress(w.To[:]), 
                w.Value, 
                common.HexToAddress(etgate.DecodeToken(w.Token)),
                []byte(w.ChainID),
                big.NewInt(int64(i)),
            )
            if err != nil {
                fmt.Printf("Failed to submit withdraw: \"%s\", Waiting...\n", err)
                break
            }
            
            fmt.Printf("Submitted withdraw: %d\nTx hash: %v\n", i, tx.Hash().Hex())
            time.Sleep(5 * time.Second)

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

        fmt.Printf("%+v\n", logs)
        for _, log := range logs {
            proof, err := g.newLogProof(log)
            if err != nil {
                fmt.Printf("Error generating logproof: \"%s\". Skipping.\n", err)
                continue
            }
            logt, err := proof.Log()
            if err != nil {
                fmt.Printf("aa: %s\n", err)
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

    return etgate.NewLogProof(receipts, log.TxIndex, log.BlockNumber)
}

func (g *gateway) appTx(etgateTx etgate.ETGateTx) error {
	return commands.AppTx(g.mintclient, g.mintkey, etgateTx, chainIDFlag)
}
