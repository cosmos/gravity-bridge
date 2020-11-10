package cli

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/cosmos/cosmos-sdk/types/errors"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	peggyTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Peggy transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	peggyTxCmd.AddCommand(flags.PostCommands(
		CmdWithdrawToETH(cdc),
		CmdRequestBatch(cdc),
		CmdUpdateEthAddress(cdc),
		CmdValsetRequest(cdc),
		GetUnsafeTestingCmd(cdc),
	)...)

	return peggyTxCmd
}

func GetUnsafeTestingCmd(cdc *codec.Codec) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "unsafe_testing",
		Short:                      "helpers for testing. not going into production",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand(flags.PostCommands(
		CmdUnsafeETHPrivKey(),
		CmdUnsafeETHAddr(),
	)...)

	return testingTxCmd
}

// GetCmdUpdateEthAddress updates the network about the eth address that you have on record.
func CmdUpdateEthAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-eth-addr [eth_private_key]",
		Short: "Update your Ethereum address which will be used for signing executables for the `multisig set`",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosAddr := cliCtx.GetFromAddress()

			privKeyString := args[0][2:]

			// Make Eth Signature over validator address
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				return err
			}

			hash := ethCrypto.Keccak256(cosmosAddr.Bytes())
			signature, err := types.NewEthereumSignature(hash, privateKey)
			if err != nil {
				return sdkerrors.Wrap(err, "signing cosmos address with Ethereum key")
			}
			// You've got to do all this to get an Eth address from the private key
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return sdkerrors.Wrap(err, "casting public key to ECDSA")
			}
			ethAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA)

			msg := types.NewMsgSetEthAddress(types.EthereumAddress(ethAddress), cosmosAddr, hex.EncodeToString(signature))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdValsetRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request",
		Short: "Trigger a new `multisig set` update request on the cosmos side",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			msg := types.NewMsgValsetRequest(cosmosAddr)

			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdUnsafeETHPrivKey() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-eth-key",
		Short: "Generate and print a new ecdsa key",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := ethCrypto.GenerateKey()
			if err != nil {
				return errors.Wrap(err, "can not generate key")
			}
			k := "0x" + hex.EncodeToString(ethCrypto.FromECDSA(key))
			println(k)
			return nil
		},
	}
}

func CmdUnsafeETHAddr() *cobra.Command {
	return &cobra.Command{
		Use:   "eth-address",
		Short: "Print address for an ECDSA eth key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			privKeyString := args[0][2:]
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				log.Fatal(err)
			}
			// You've got to do all this to get an Eth address from the private key
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}
			ethAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			println(ethAddress)
			return nil
		},
	}
}

func CmdWithdrawToETH(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [from_key_or_cosmos_address] [to_eth_address] [amount] [bridge_fee]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the Ethereum bridge contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}
			bridgeFee, err := sdk.ParseCoin(args[3])
			if err != nil {
				return sdkerrors.Wrap(err, "bridge fee")
			}

			// Make the message
			msg := types.MsgSendToEth{
				Sender:      cosmosAddr,
				DestAddress: types.NewEthereumAddress(args[1]),
				Amount:      amount,
				BridgeFee:   bridgeFee,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdRequestBatch(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "build-batch [voucher_denom]",
		Short: "Build a new batch on the cosmos side for pooled withdrawal transactions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			denom, err := types.AsVoucherDenom(args[0])
			if err != nil {
				return sdkerrors.Wrap(err, "denom")
			}

			msg := types.MsgRequestBatch{
				Requester: cosmosAddr,
				Denom:     denom,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
