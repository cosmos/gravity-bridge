package main

import (
    "os"

    "github.com/spf13/cobra"

    "github.com/tendermint/tmlibs/cli"

    "github.com/tendermint/basecoin/cmd/basecoin/commands"
    "github.com/tendermint/basecoin/types"

    "../../plugins/etgate"

)

func main() {
    var RootCmd = &cobra.Command {
        Use: "etgate",
        Short: "ethereum log relaying plugin for basecoin"
    }

    RootCmd.AddCommand(
        commands.InitCmd,
        commands.StartCmd,
        commands.UnsafeResetAllCmd,
        commands.VersionCmd,
    )

    commands.RegisterStartPlugin("etgate", func() types.Plugin { return etgate.New() })
    cmd := cli.PrepareMainCmd(RootCmd, "ETG", os.ExpandEnv("$HOME/.etgate"))
    cmd.Excute()
}

