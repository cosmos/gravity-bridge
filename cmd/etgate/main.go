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
        Short: "ethereum log relaying plugin for basecoin",
    }

    RootCmd.AddCommand(
        commands.InitCmd,
        commands.StartCmd,
        commands.RelayCmd,
        GateCmd,
        commands.UnsafeResetAllCmd,
        commands.VersionCmd,
    )
    
    commands.RegisterStartPlugin("ETGATE", func() types.Plugin { return etgate.New() })

    cmd := cli.PrepareMainCmd(RootCmd, "ETGATE", os.ExpandEnv("$HOME/.etgate/server"))
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}

