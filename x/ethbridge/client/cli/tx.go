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
		Use:   "create-claim [nonce] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim_type]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(6),
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

			var claimType types.ClaimType
			if value, ok := types.StringToClaimType[args[5]]; ok {
				claimType = value
			} else {
				return types.ErrInvalidClaimType()
			}

			ethBridgeClaim := types.NewEthBridgeClaim(nonce, ethereumSender, cosmosReceiver, validator, amount, claimType)
			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurn is the CLI command for burning some of your coins and triggering an event
func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [ethereum-chain-id] [token-contract-address] [cosmos-sender-address] [ethereum-receiver-address] [amount]",
		Short: "burn some coins!",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			token := types.NewEthereumAddress(args[1])

			cosmosSender, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			ethereumReceiver := types.NewEthereumAddress(args[3])

			amount, err := sdk.ParseCoins(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgBurn(ethereumChainID, token, cosmosSender, ethereumReceiver, amount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLock is the CLI command for locking some of your coins and triggering an event
func GetCmdLock(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "lock [cosmos-sender-address] [ethereum-receiver-address] [amount]",
		Short: "lock some coins!",
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

			msg := types.NewMsgLock(cosmosSender, ethereumReceiver, amount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
