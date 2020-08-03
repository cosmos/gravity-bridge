package cli

import (
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	peggyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the peggy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	peggyQueryCmd.AddCommand(flags.GetCommands(
		// GetCmdResolveName(storeKey, cdc),
		// GetCmdWhois(storeKey, cdc),
		// GetCmdNames(storeKey, cdc),
		CmdGetCurrentValset(storeKey, cdc),
		CmdGetValsetRequest(storeKey, cdc),
		CmdGetValsetConfirm(storeKey, cdc),
	)...)

	return peggyQueryCmd
}

// // GetCmdResolveName queries information about a name
// func GetCmdResolveName(storeKey string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "resolve [name]",
// 		Short: "resolve name",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			name := args[0]

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", storeKey, name), nil)
// 			if err != nil {
// 				fmt.Printf("could not resolve name - %s \n", name)
// 				return nil
// 			}

// 			var out types.QueryResResolve
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }

// // GetCmdWhois queries information about a domain
// func GetCmdWhois(storeKey string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "whois [name]",
// 		Short: "Query whois info of name",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			name := args[0]

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", storeKey, name), nil)
// 			if err != nil {
// 				fmt.Printf("could not resolve whois - %s \n", name)
// 				return nil
// 			}

// 			var out types.Whois
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }

func CmdGetCurrentValset(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current-valset",
		Short: "Query current valset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/currentValset", storeKey), nil)
			if err != nil {
				fmt.Printf("could not get valset")
				return nil
			}

			var out types.Valset
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetValsetRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request [nonce]",
		Short: "Get requested valset with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			nonce := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, nonce), nil)
			if err != nil {
				return err
			}

			var out types.Valset
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetValsetConfirm(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [bech32 validator address]",
		Short: "Get valset confirmation with a particular nonce from a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetConfirm/%s/%s", storeKey, args[0], args[1]), nil)
			if err != nil {
				fmt.Printf("could not get valset")
				return nil
			}

			var out []types.MsgValsetConfirm
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// // GetCmdNames queries a list of all names
// func GetCmdNames(storeKey string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "names",
// 		Short: "names",
// 		// Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names", storeKey), nil)
// 			if err != nil {
// 				fmt.Printf("could not get query names\n")
// 				return nil
// 			}

// 			var out types.QueryResNames
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }
