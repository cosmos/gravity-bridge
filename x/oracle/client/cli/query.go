package cli

import (
	"fmt"
	"strconv"

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
			nonce, stringError := strconv.Atoi(args[0])
			if stringError != nil {
				return stringError
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/prophecy/%s", queryRoute, nonce), nil)
			if err != nil {
				fmt.Printf("could not find with given nonce %s \n", string(nonce))
				return nil
			}

			var out oracle.BridgeProphecy
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
