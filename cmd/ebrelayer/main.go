package main

// 	Main (ebrelayer) Implements CLI commands for the Relayer
//		service, such as initialization and event relay.

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	app "github.com/cosmos/peggy/app"
	relayer "github.com/cosmos/peggy/cmd/ebrelayer/relayer"
	txs "github.com/cosmos/peggy/cmd/ebrelayer/txs"
	utils "github.com/cosmos/peggy/cmd/ebrelayer/utils"
)

var cdc *codec.Codec

const (
	// FlagRPCURL defines the URL for the tendermint RPC connection
	FlagRPCURL = "rpc-url"
)

func init() {

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	cdc = app.MakeCodec()

	DefaultCLIHome := os.ExpandEnv("$HOME/.ebcli")

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(FlagRPCURL, "", "RPC URL of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		relayerCmd(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "EBRELAYER", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		log.Fatal("failed executing CLI command", err)
	}
}

var rootCmd = &cobra.Command{
	Use:          "ebrelayer",
	Short:        "Relayer service which listens for and relays ethereum smart contract events",
	SilenceUsage: true,
}

//	relayerCmd : Initializes a relayer service subscribed to Ethereum and Tendermint.
//				 On Ethereum, the relayer streams live events from the deployed smart
//				 contracts. On Tendermint, it streams specific message types relevant
//				 to Peggy. The service automatically signed messages containing the
//				 event/message data, repackages it into a new transaction, and relays
//				 it to the opposite chain for processing.
func relayerCmd() *cobra.Command {
	//nolint:lll
	relayerCmd := &cobra.Command{
		Use:     "init [tendermintNode] [web3Provider] [bridgeContractAddress] [validatorFromName] --chain-id [chain-id]",
		Short:   "Initializes a web socket which streams live events from a smart contract and relays them to the Cosmos network",
		Args:    cobra.ExactArgs(4),
		Example: "ebrelayer init tcp://localhost:26657 ws://127.0.0.1:7545/ 0x0823eFE0D0c6bd134a48cBd562fE4460aBE6e92c validator --chain-id=peggy",
		RunE:    RunRelayerCmd,
	}

	return relayerCmd
}

// RunRelayerCmd executes the relayerCmd with the provided parameters
func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())

	// Load the validator's Ethereum private key
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Wrap(err, "invalid [ETHEREUM_PRIVATE_KEY] from .env")
	}

	// Parse chain's ID
	chainID := viper.GetString(flags.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return errors.New("Must specify a 'chain-id'")
	}

	tendermintNode := args[0]

	// Parse ethereum provider
	ethereumProvider := args[1]
	if !relayer.IsWebsocketURL(ethereumProvider) {
		return fmt.Errorf("invalid [web3-provider]: %s", ethereumProvider)
	}

	// Parse the address of the deployed contract
	if !common.IsHexAddress(args[2]) {
		return fmt.Errorf("invalid [bridge-contract-address]: %s", args[2])
	}
	contractAddress := common.HexToAddress(args[2])

	// Parse the validator's moniker
	validatorFrom := args[3]

	// Parse Tendermint RPC URL
	rpcURL := viper.GetString(FlagRPCURL)

	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", rpcURL)
		}
	}

	// Load validator details
	validatorAddress, validatorName, err := utils.LoadValidatorCredentials(validatorFrom, inBuf)
	if err != nil {
		return err
	}

	// Load CLI context
	cliCtx := utils.LoadTendermintCLIContext(cdc, validatorAddress, validatorName, rpcURL, chainID)

	// Load Tx builder
	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(sdkUtils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	// Start an Ethereum websocket
	go func() {
		err := relayer.InitEthereumRelayer(
			cdc,
			ethereumProvider,
			contractAddress,
			validatorName,
			validatorAddress,
			cliCtx,
			txBldr,
			privateKey,
		)
		if err != nil {
			log.Printf("Ethereum relayer failed: %v", err)
		}
	}()
	// Start an Ethereum websocket
	// go relayer.InitEthereumRelayer(
	// 	cdc,
	// 	ethereumProvider,
	// 	contractAddress,
	// 	validatorName,
	// 	validatorAddress,
	// 	cliCtx,
	// 	txBldr,
	// 	privateKey,
	// )

	// Start a Tendermint websocket
	go func() {
		err := relayer.InitCosmosRelayer(
			tendermintNode,
			ethereumProvider,
			contractAddress,
			privateKey,
		)
		if err != nil {
			log.Printf("Cosmos relayer failed: %v", err)
		}
	}()

	// // Start a Tendermint websocket
	// go relayer.InitCosmosRelayer(
	// 	tendermintNode,
	// 	ethereumProvider,
	// 	contractAddress,
	// 	privateKey,
	// )

	return nil
}

func initConfig(cmd *cobra.Command) error {
	return viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
