package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/peggy/x/ethbridge/querier"
	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetEthBridgeProphecy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "prophecy [chain-id] [nonce] [ethereum-sender]",
		Short: "Query prophecy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			ethereumSender := types.NewEthereumAddress(args[2])

			bz, err := cdc.MarshalJSON(types.NewQueryEthProphecyParams(chainID, nonce, ethereumSender))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, querier.QueryEthProphecy)
			res, _, err := cliCtx.QueryWithData(route, bz)
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
