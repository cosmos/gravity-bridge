package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
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

			bridgeContract := ethbridge.NewEthereumAddress(args[1])

			nonce, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			symbol := args[3]
			tokenContract := ethbridge.NewEthereumAddress(args[4])
			ethereumSender := ethbridge.NewEthereumAddress(args[5])

			bz, err := cdc.MarshalJSON(types.NewQueryNFTProphecyParams(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryEthProphecy)
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
