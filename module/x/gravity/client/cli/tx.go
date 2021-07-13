package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
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
		Use:     "send-to-etheruem [ethereum-reciever] [send-coins] [fee-coins]",
		Aliases: []string{"send", "transfer"},
		Args:    cobra.ExactArgs(3),
		Short:   "Send tokens from cosmos chain to connected etheruem chain",
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
				return fmt.Errorf("must be a valid etheruem address got %s", args[0])
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
		Use:   "cancel-send-to-etheruem [id]", // TODO(levi) this argument name is vague (but matches what we call it everywhere)
		Args:  cobra.ExactArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			if from == nil {
				return fmt.Errorf("must pass from flag")
			}

			var ( // args
				id uint64 // TODO(levi) init from args[0]
			)

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
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				denom  string         // TODO(levi) init and validate from args[0]
				signer sdk.AccAddress // TODO(levi) init and validate from args[1]
			)

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
		Use:   "set-delegate-keys [validator-address] [orchestrator-address] [ethereum-address]",
		Args:  cobra.ExactArgs(3),
		Short: "Set gravity delegate keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromHex(args[0])
			if err != nil {
				return err
			}
			orcAddr, err := sdk.AccAddressFromHex(args[1])
			if err != nil {
				return err
			}
			ethAddr := args[2]

			msg := types.NewMsgDelegateKeys(valAddr, orcAddr, ethAddr)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// func CmdUnsafeETHPrivKey() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "gen-eth-key",
// 		Short: "Generate and print a new ecdsa key",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			key, err := ethCrypto.GenerateKey()
// 			if err != nil {
// 				return sdkerrors.Wrap(err, "can not generate key")
// 			}
// 			k := "0x" + hex.EncodeToString(ethCrypto.FromECDSA(key))
// 			println(k)
// 			return nil
// 		},
// 	}
// }

// func CmdUnsafeETHAddr() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "eth-address",
// 		Short: "Print address for an ECDSA eth key",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			privKeyString := args[0][2:]
// 			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			// You've got to do all this to get an Eth address from the private key
// 			publicKey := privateKey.Public()
// 			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 			if !ok {
// 				log.Fatal("error casting public key to ECDSA")
// 			}
// 			ethAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()
// 			println(ethAddress)
// 			return nil
// 		},
// 	}
// }
