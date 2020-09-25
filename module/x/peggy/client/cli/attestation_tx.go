package cli

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func GetObservedCmd(cdc *codec.Codec) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "observed",
		Short:                      "submit observed ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand(flags.PostCommands(
		CmdSendETHBootstrapRequest(cdc),
		CmdSendETHDepositRequest(cdc),
		CmdSendETHWithdrawalRequest(cdc),
		CmdSendETHMultiSigRequest(cdc),
	)...)

	return testingTxCmd
}

func GetApprovedCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "approved",
		Short:                      "approve an operation",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand(flags.PostCommands(
		CmdValsetConfirm(storeKey, cdc),
	)...)

	return testingTxCmd
}

func CmdSendETHBootstrapRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap [eth chain id] [eth contract address] [nonce] [allowed_validators] [validator_powers] [peggy_id] [start_threshold]",
		Short: "Submit a claim that the bridge contract bootstrap was completed on the Ethereum side",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			ethChainID := args[0]
			ethContractAddress := args[1]
			nonce, err := parseNonce(args[2])
			if err != nil {
				return err
			}
			var validators []types.EthereumAddress
			var powers []uint64
			for _, v := range strings.Split(args[3], ",") {
				validators = append(validators, types.NewEthereumAddress(v))
			}

			for _, v := range strings.Split(args[4], ",") {
				p, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return sdkerrors.Wrap(err, "power")
				}
				powers = append(powers, p)
			}

			startThreshold, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return sdkerrors.Wrap(err, "start threshold")
			}

			// Make the message
			msg := types.MsgCreateEthereumClaims{
				EthereumChainID:       ethChainID,
				BridgeContractAddress: types.NewEthereumAddress(ethContractAddress),
				Orchestrator:          cosmosAddr,
				Claims: []types.EthereumClaim{
					types.EthereumBridgeBootstrappedClaim{
						Nonce:               nonce,
						AllowedValidatorSet: validators,
						ValidatorPowers:     powers,
						PeggyID:             []byte(args[5]), // simplest solution without decoding
						StartThreshold:      startThreshold,
					},
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
func CmdSendETHDepositRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [eth chain id] [eth contract address] [nonce] [cosmos receiver] [amount] [eth erc20 symbol] [eth erc20 contract addr] [eth sender address]",
		Short: "Submit a claim that a deposit was made on the Ethereum side",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			ethChainID := args[0]
			ethContractAddress := args[1]
			nonce, err := parseNonce(args[2])
			if err != nil {
				return err
			}
			receiverAddr, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return sdkerrors.Wrap(err, "cosmos receiver")
			}
			amount, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}
			tokenSymbol := args[5]

			// Make the message
			tokenContractAddr := types.NewEthereumAddress(args[6])
			ethSenderAddr := types.NewEthereumAddress(args[7])
			msg := types.MsgCreateEthereumClaims{
				EthereumChainID:       ethChainID,
				BridgeContractAddress: types.NewEthereumAddress(ethContractAddress),
				Orchestrator:          cosmosAddr,
				Claims: []types.EthereumClaim{
					types.EthereumBridgeDepositClaim{
						Nonce:          nonce,
						ERC20Token:     types.NewERC20Token(uint64(amount), tokenSymbol, tokenContractAddr),
						EthereumSender: ethSenderAddr,
						CosmosReceiver: receiverAddr,
					},
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdSendETHWithdrawalRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdrawal [eth chain id] [eth contract address] [nonce]",
		Short: "Submit a claim that a withdrawal was executed on the Ethereum side",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			ethChainID := args[0]
			ethContractAddress := args[1]
			nonce, err := parseNonce(args[2])
			if err != nil {
				return err
			}
			msg := types.MsgCreateEthereumClaims{
				EthereumChainID:       ethChainID,
				BridgeContractAddress: types.NewEthereumAddress(ethContractAddress),
				Orchestrator:          cosmosAddr,
				Claims: []types.EthereumClaim{
					types.EthereumBridgeWithdrawalBatchClaim{
						Nonce: nonce,
					},
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdSendETHMultiSigRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "multisig-update [eth chain id] [eth contract address] [nonce]",
		Short: "Submit a claim that the 'multisig set' update was executed on the Ethereum side",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			ethChainID := args[0]
			ethContractAddress := types.NewEthereumAddress(args[1])
			nonce, err := parseNonce(args[2])
			if err != nil {
				return err
			}
			msg := types.MsgCreateEthereumClaims{
				EthereumChainID:       ethChainID,
				BridgeContractAddress: ethContractAddress,
				Orchestrator:          cosmosAddr,
				Claims: []types.EthereumClaim{
					types.EthereumBridgeMultiSigUpdateClaim{
						Nonce: nonce,
					},
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdValsetConfirm(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [eth private key]",
		Short: "Sign a `multisig set` update for given nonce with the Ethereum key and submit to cosmos side",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Make Eth Signature over valset
			privKeyString := args[1][2:]
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				return err
			}

			nonce := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, nonce), nil)
			if err != nil {
				fmt.Printf("could not get valset")
				return nil
			}

			var valset types.Valset
			cdc.MustUnmarshalJSON(res, &valset)
			checkpoint := valset.GetCheckpoint()

			signature, err := ethCrypto.Sign(checkpoint, privateKey)
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()
			// Make the message
			msg := types.NewMsgValsetConfirm(valset.Nonce, cosmosAddr, hex.EncodeToString(signature))

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// todo: we mix nonces as int64 and base64 bytes at the moment
func parseNonce(nonceArg string) (types.Nonce, error) {
	if len(nonceArg) != base64.StdEncoding.EncodedLen(8) {
		// not a byte nonce byte representation
		v, err := strconv.ParseUint(nonceArg, 10, 64)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "nonce")
		}
		return types.NonceFromUint64(v), nil
	}
	return base64.StdEncoding.DecodeString(nonceArg)
}
