package main

import (
    "os"

    "github.com/spf13/cobra"

    keycmd "github.com/tendermint/go-crypto/cmd"
    "github.com/tendermint/light-client/commands"
    "github.com/tendermint/light-client/commands/proofs"
    "github.com/tendermint/light-client/commands/txs"
    ecmd "./commands"
    "github.com/tendermint/tmlibs/cli"
)

var ETCli = &cobra.Command {
    Use: "etcli",
    Short: "ETGate light client",
}

func main() {
    commands.AddBasicFlags(ETCli)

    pr := proofs.RootCmd
    pr.AddCommand(proofs.TxCmd)
    pr.AddCommand(proofs.KeyCmd)
    //pr.AddCommand(ecmd.AccountQueryCmd)

//    proofs.TxPresenters.Register("etgate", i)
    tr := txs.RootCmd
    tr.AddCommand(ecmd.WithdrawTxCmd)

    ETCli.AddCommand(
        commands.InitCmd,
        commands.ResetCmd,
        keycmd.RootCmd,
        pr,
        tr,
    )

    cmd := cli.PrepareMainCmd(ETCli, "ETC", os.ExpandEnv("$HOME/.etgate/client"))
    cmd.Execute()
}
