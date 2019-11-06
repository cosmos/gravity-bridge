package cli

import (
	"errors"
	"strconv"
	"strings"

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
		Use:   "create-claim [ethereum-chain-id] [bridge-contract] [nonce] [symbol] [token-contract] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

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
				return errors.New("must specify a token symbol/denomination, including 'eth' for Ethereum")
			}

			tokenContract := types.NewEthereumAddress(args[4])

			ethereumSender := types.NewEthereumAddress(args[5])

			cosmosReceiver, err := sdk.AccAddressFromBech32(args[6])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[7])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(args[8])
			if err != nil {
				return err
			}

			ethBridgeClaim := types.NewEthBridgeClaim(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender, cosmosReceiver, validator, amount)
			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurnEth is the CLI command for burning some of your eth and triggering an event
func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [cosmos-sender-address] [ethereum-receiver-address] [amount]",
		Short: "This should be used to burn cETH or cERC20. It will burn your coins on the Cosmos Chain, removing them from your account and deducting them from the supply. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original ETH/ERC20 to you from the Ethereum contract!",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			ethereumReceiver := types.NewEthereumAddress(args[1])

			amount, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgBurn(cosmosSender, ethereumReceiver, amount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
