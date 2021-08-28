package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func GetTxCmd(storeKey string) *cobra.Command {
	gravityTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gravity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gravityTxCmd.AddCommand(
		CmdSendToEthereum(),
		CmdCancelSendToEthereum(),
		CmdRequestBatchTx(),
		CmdSetDelegateKeys(),
	)

	return gravityTxCmd
}

func CmdSendToEthereum() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-to-ethereum [ethereum-reciever] [send-coins] [fee-coins]",
		Aliases: []string{"send", "transfer"},
		Args:    cobra.ExactArgs(3),
		Short:   "Send tokens from cosmos chain to connected ethereum chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			if from == nil {
				return fmt.Errorf("must pass from flag")
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("must be a valid ethereum address got %s", args[0])
			}

			// Get amount of coins
			sendCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			feeCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSendToEthereum(from, common.HexToAddress(args[0]).Hex(), sendCoin, feeCoin)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelSendToEthereum() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-ethereum [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel ethereum send by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			if from == nil {
				return fmt.Errorf("must pass from flag")
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelSendToEthereum(id, from)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatchTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-tx [denom] [signer]",
		Args:  cobra.ExactArgs(2),
		Short: "Request batch transaction for denom by signer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			signer, err := sdk.AccAddressFromHex(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestBatchTx(denom, signer)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSetDelegateKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-delegate-keys [validator-address] [orchestrator-address] [ethereum-address] [ethereum-signature]",
		Args:  cobra.ExactArgs(4),
		Short: "Set gravity delegate keys",
		Long: `Set a validator's Ethereum and orchestrator addresses. The validator must
sign over a binary Proto-encoded DelegateKeysSignMsg message. The message contains
the validator's address and operator account current nonce.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			orcAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			ethAddr, err := parseContractAddress(args[2])
			if err != nil {
				return err
			}

			ethSig, err := hexutil.Decode(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegateKeys(valAddr, orcAddr, ethAddr, ethSig)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
