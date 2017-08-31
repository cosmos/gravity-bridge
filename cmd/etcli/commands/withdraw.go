package commands

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

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

    etcmd "../../../commands"
    "../../../plugins/etgate"
)

const (
    FlagTo       = "to"
    FlagToken    = "token"
    FlagSequence = "sequence"
)

var valueFlag int64

func init() {
    flags := WithdrawTxCmd.Flags()
    flags.String(FlagTo, "", "Destination ethereum address")
    flags.Int64VarP(&valueFlag, "value", "v",  0, "Value of coins to send")
    flags.String(FlagToken, "", "Token ethereum address")
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

    key := fmt.Sprintf("etgate,withdraw,etgate-chain")
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
        ChainID: "etgate-chain", // change it later
        Sequence: seq+1,
    }
    fmt.Printf("%+v\n", inner)
    acc, err := etcmd.GetAccWithClient(node, info.Address[:])
    if err != nil {
        return err
    }

    feeCoins := bctypes.Coin{Denom: "mycoin", Amount: 1}

    enctoken := ""
    for i, s := range token {
        if i == 1 {
            enctoken = enctoken + string(s)
            continue
        }
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

    return txcmd.OutputTx(res)
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
