package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/peggy/x/ethbridge"
	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/spf13/cobra"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetEthBridgeProphecy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-prophecy [nonce] [ethereum-sender]",
		Short: "get prophecy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			nonce, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			ethereumSender := types.NewEthereumAddress(args[1])

			bz, err := cdc.MarshalJSON(ethbridge.NewQueryEthProphecyParams(nonce, ethereumSender))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, ethbridge.QueryEthProphecy)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var out types.QueryEthProphecyResponse
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(out)
		},
	}
}
