package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	// Group ethbridge queries under a subcommand
	ethBridgeQueryCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "Querying commands for the ethbridge module",
	}

	ethBridgeQueryCmd.AddCommand(flags.GetCommands(
		GetCmdGetEthBridgeProphecy(storeKey, cdc),
	)...)

	return ethBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ethBridgeTxCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "EthBridge transactions subcommands",
	}

	ethBridgeTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateEthBridgeClaim(cdc),
		GetCmdBurn(cdc),
		GetCmdLock(cdc),
	)...)

	return ethBridgeTxCmd
}
