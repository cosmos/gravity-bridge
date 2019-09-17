package cli

import (
	"strconv"

	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetCmdCreateEthBridgeClaim is the CLI command for creating a claim on an ethereum prophecy
func GetCmdCreateEthBridgeClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-claim [nonce] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			nonce, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			ethereumSender := types.NewEthereumAddress(args[1])

			cosmosReceiver, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(args[4])
			if err != nil {
				return err
			}

			ethBridgeClaim := types.NewEthBridgeClaim(nonce, ethereumSender, cosmosReceiver, validator, amount)
			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
