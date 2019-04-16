package cli

import (
	"fmt"

	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetCmdGetProphecy queries information about a name
func GetCmdGetProphecy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-prophecy nonce",
		Short: "get prophecy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/prophecy/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Printf(err.Error())
				return nil
			}

			var out oracle.BridgeProphecy
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
