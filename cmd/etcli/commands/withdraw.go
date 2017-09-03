package commands

import (
    "path/filepath"
    "os"
    "io/ioutil"
    "fmt"
    "math/big"
    "bytes"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/bgentry/speakeasy"

    "github.com/tendermint/light-client/commands"
    txcmd "github.com/tendermint/light-client/commands/txs"
    keycmd "github.com/tendermint/go-crypto/cmd"
//    "github.com/tendermint/tendermint/rpc/client"
    bclicmd "github.com/tendermint/basecoin/cmd/basecli/commands"
    bctypes "github.com/tendermint/basecoin/types"

//    data "github.com/tendermint/go-wire/data"
    keys "github.com/tendermint/go-crypto/keys"
    "github.com/tendermint/go-wire"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"

    etcmd "../../../commands"
    "../../../plugins/etgate"
    "../../../contracts"
)

const (
    FlagTo       = "to"
    FlagValue    = "value"
    FlagToken    = "token"
    FlagSequence = "sequence"
    FlagKey      = "key"
    FlagDatadir  = "datadir"
    FlagTestnet  = "testnet"
    FlagAddress  = "address"
)

var valueFlag int64

func init() {
    flags := WithdrawTxCmd.Flags()
    flags.String(FlagTo, "", "Destination ethereum address")
    flags.Int64(FlagValue, 0, "Value of coins to send")
    flags.String(FlagToken, "", "Token ethereum address")
    flags.String(FlagKey, "", "Ethereum key json file path")
    flags.String(FlagDatadir, filepath.Join(os.Getenv("HOME"), ".ethereum"), "Data directory for the databases and keystore")
    flags.Bool(FlagTestnet, false, "Ropsten network: pre-configured test network")
    flags.String(FlagAddress, "", "ETGate contract address on Ethereum chain")
}

var WithdrawTxCmd = &cobra.Command {
    Use: "withdraw",
    RunE: withdrawCmd,
}

func getInfo() (keys.Info, error) {
    manager := keycmd.GetKeyManager()
    info, err := manager.Get(viper.GetString("name"))
    return info, err

}

func withdrawCmd(cmd *cobra.Command, args []string) error {
    originChainID := "etgate-chain" // for now

    node := commands.GetNode()
    info, err := getInfo()
    if err != nil {
        return err
    }
/*
    acc, err := etcmd.GetAccWithClient(node, info.Address[:])
    if err != nil {
        return err 
    }

    sequence := acc.Sequence + 1
    */

    key := fmt.Sprintf("etgate,withdraw,%s", /*change it later*/originChainID)
    query, err := etcmd.QueryWithClient(node, []byte(key))
    if err != nil {
        return err
    }
    var seq uint64
    if len(query.Value) == 0 {
        seq = 0
    } else {
        if err = wire.ReadBinaryBytes(query.Value, &seq); err != nil {
            return err
        }
    }

    token := viper.GetString(FlagToken)

    inner := etgate.ETGateWithdrawTx {
        To: common.HexToAddress(viper.GetString(FlagTo)),
        Value: uint64(valueFlag),
        Token: common.HexToAddress(token),
        ChainID: originChainID, // change it later
        Sequence: seq+1,
    }
    acc, err := etcmd.GetAccWithClient(node, info.Address[:])
    if err != nil {
        return err
    }

    feeCoins := bctypes.Coin{Denom: "mycoin", Amount: 1}

    enctoken := ""
    for _, s := range token[2:] {
        enctoken = enctoken + string(s+32) // coin denom dosent work on numerics
    }
    ethCoins := bctypes.Coin{Denom: enctoken, Amount: int64(valueFlag)}


    tx := &bctypes.AppTx {
        Gas: 0,
        Fee: feeCoins,
        Name: "ETGATE",
        Input: bctypes.NewTxInput(info.PubKey, bctypes.Coins{ethCoins}, acc.Sequence+1),
        Data: wire.BinaryBytes(struct {
            etgate.ETGateTx `json:"unwrap"`
        }{inner}),
    }

//    fmt.Printf("%+v\n%+v\n", info.PubKey, tx.String())

    res, err := bclicmd.BroadcastAppTx(tx)

    if err != nil {
        return err
    }

    fmt.Printf("%s\n", txcmd.OutputTx(res))

    password, err := speakeasy.Ask("Enter the password for keyfile: ")
    if err != nil {
        return err
    }

    keyBytes, err := ioutil.ReadFile(viper.GetString(FlagKey))
    if err != nil {
        return err
    }
/*
    priv, err := keystore.DecryptKey(keyBytes, password)
    if err != nil {
        return err
    }
*/

// TODO: move boilerplate to commands
    var datadir string
    if viper.GetBool(FlagTestnet) {
        datadir = filepath.Join(viper.GetString(FlagDatadir), "testnet")
    } else {
        datadir = viper.GetString(FlagDatadir)
    }

    ethclient, err := ethclient.Dial(filepath.Join(datadir, "geth.ipc"))

    if err != nil {
        return err
    }

    contract, err := contracts.NewETGate(common.HexToAddress(viper.GetString(FlagAddress)), ethclient)
    if err != nil {
        return err
    }

    _, err = bind.NewTransactor(bytes.NewReader(keyBytes), password)     
    if err != nil {
        return err
    }

    fmt.Println("Waiting for the header to be uploaded...")

    for {
        withdrawable, err := contract.Withdrawable(
            nil, 
            big.NewInt(int64(res.Height)), 
            /*change it later*/[]byte(originChainID), common.HexToAddress(viper.GetString(FlagTo)), 
            uint64(viper.GetInt64(FlagValue)))
        if err != nil {
            return err
        }
        if withdrawable {
            break
        }
    }

//    contract.Withdraw

    return nil
}
/*
func withdrawCmd(cmd *cobra.Command, args []string) error {
    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", ""))
    si
    withdrawTx := new(etgate.ETGateWithdrawTx)

    found, err := txcmd.LoadJSON(withdrawTx)
    if err != nil {
        return err
    }
    if !found {
        withdrawTx.To = common.HexToAddress(toFlag)
        withdrawTx.Value = valueFlag
        withdratTx.Token = common.HexToAddress(tokenFlag)
        withdrawTx.ChainID = "etgate-chain"
        withdrawTx.Sequence = sequenceFlag
    }

    data := []byte(wire.BinaryBytes(struct {
        WithdrawTx `json:"unwrap"`
    }{withdrawTx}))

    smallCoins := bctypes.Coin{Denom: "mycoin", Amount: 1}

    input := bctypes.NewTxInput(txcmd.GetSigner(), bctypes.Coins{smallCoins}, sequence)

    tx := &bctypes.AppTx {
        Gas: 0,
        Fee: smallCoins,
        Name: "ETGATE",
        Input: input,
        Data: data,
    }

    tx.SignBytes()
}*/
