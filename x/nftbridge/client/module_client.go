package client

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/peggy/x/nftbridge/client/cli"
	"github.com/cosmos/peggy/x/nftbridge/client/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	// Group nftbridge queries under a subcommand
	nftBridgeQueryCmd := &cobra.Command{
		Use:   "nftbridge",
		Short: "Querying commands for the nftbridge module",
	}

	nftBridgeQueryCmd.AddCommand(flags.GetCommands(
		cli.GetCmdGetNFTBridgeProphecy(storeKey, cdc),
	)...)

	return nftBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nftBridgeTxCmd := &cobra.Command{
		Use:   "nftbridge",
		Short: "NFTBridge transactions subcommands",
	}

	nftBridgeTxCmd.AddCommand(flags.PostCommands(
		cli.GetCmdCreateNFTBridgeClaim(cdc),
		cli.GetCmdBurnNFT(cdc),
		cli.GetCmdLockNFT(cdc),
	)...)

	return nftBridgeTxCmd
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	rest.RegisterRESTRoutes(cliCtx, r, storeName)
}
