package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetEthBridgeProphecy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-prophecy identifier",
		Short: "get prophecy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			bz, err := cdc.MarshalJSON(ethbridge.NewQueryEthProphecyParams(id))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, ethbridge.QueryEthProphecy)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf(err.Error())
				return nil
			}

			var out types.QueryEthProphecyResponse
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
