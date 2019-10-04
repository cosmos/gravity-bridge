package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/peggy/x/ethbridge/querier"
	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetEthBridgeProphecy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "prophecy [ethereum-chain-id] [bridge-contract] [nonce] [symbol] [token-contract] [ethereum-sender]",
		Short: "Query prophecy",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ethereumChainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			bridgeContract := types.NewEthereumAddress(args[1])

			nonce, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			symbol := args[3]
			if strings.TrimSpace(symbol) == "" {
				return errors.New("Error: must specify a token symbol, including 'eth' for Ethereum")
			}

			tokenContract := types.NewEthereumAddress(args[4])

			ethereumSender := types.NewEthereumAddress(args[5])

			bz, err := cdc.MarshalJSON(types.NewQueryEthProphecyParams(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender))
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
