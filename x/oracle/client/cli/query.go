package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
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

			bz, err := cdc.MarshalJSON(oracle.NewQueryProphecyParams(id))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryProphecy)
			res, err := cliCtx.QueryWithData(route, bz)
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
