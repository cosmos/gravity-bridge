package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/types/errors"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetTxCmd(storeKey string) *cobra.Command {
	peggyTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Peggy transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	peggyTxCmd.AddCommand([]*cobra.Command{
		CmdWithdrawToETH(),
		CmdRequestBatch(),
		CmdUpdateEthAddress(),
		CmdValsetRequest(),
		GetUnsafeTestingCmd(),
	}...)

	return peggyTxCmd
}

func GetUnsafeTestingCmd() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "unsafe_testing",
		Short:                      "helpers for testing. not going into production",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		CmdUnsafeETHPrivKey(),
		CmdUnsafeETHAddr(),
	}...)

	return testingTxCmd
}

// GetCmdUpdateEthAddress updates the network about the eth address that you have on record.
func CmdUpdateEthAddress() *cobra.Command {
	return &cobra.Command{
		Use:   "update-eth-addr [eth_private_key]",
		Short: "Update your Ethereum address which will be used for signing executables for the `multisig set`",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.ReadTxCommandFlags(client.GetClientContextFromCmd(cmd), cmd.Flags())
			if err != nil {
				return err
			}

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

			msg := types.NewMsgSetEthAddress(ethAddress.String(), cosmosAddr, hex.EncodeToString(signature))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
}

func CmdValsetRequest() *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request",
		Short: "Trigger a new `multisig set` update request on the cosmos side",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.ReadTxCommandFlags(client.GetClientContextFromCmd(cmd), cmd.Flags())
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			msg := types.NewMsgValsetRequest(cosmosAddr)

			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
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

func CmdWithdrawToETH() *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [from_key_or_cosmos_address] [to_eth_address] [amount] [bridge_fee]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the Ethereum bridge contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.ReadTxCommandFlags(client.GetClientContextFromCmd(cmd), cmd.Flags())
			if err != nil {
				return err
			}
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
				Sender:    cosmosAddr.String(),
				EthDest:   args[1],
				Amount:    amount,
				BridgeFee: bridgeFee,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
}

func CmdRequestBatch() *cobra.Command {
	return &cobra.Command{
		Use:   "build-batch [token_contract_address]",
		Short: "Build a new batch on the cosmos side for pooled withdrawal transactions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.ReadTxCommandFlags(client.GetClientContextFromCmd(cmd), cmd.Flags())
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()

			// TODO: better denom searching
			msg := types.MsgRequestBatch{
				Requester: cosmosAddr.String(),
				Denom:     fmt.Sprintf("peggy/%s", args[0]),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
}
