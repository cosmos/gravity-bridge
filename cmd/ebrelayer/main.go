package main

// 	Main (ebrelayer) : Implements CLI commands for the Relayer
//		service, such as initialization and event relay.

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	app "github.com/cosmos/peggy/app"
	relayer "github.com/cosmos/peggy/cmd/ebrelayer/relayer"
	txs "github.com/cosmos/peggy/cmd/ebrelayer/txs"
	utils "github.com/cosmos/peggy/cmd/ebrelayer/utils"
)

var appCodec *amino.Codec

const FlagRPCURL = "rpc-url"

func init() {

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	appCodec = app.MakeCodec()

	DefaultCLIHome := os.ExpandEnv("$HOME/.ebcli")

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
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
	relayerCmd := &cobra.Command{
		Use:     "init [tendermintNode] [web3Provider] [bridgeContractAddress] [validatorFromName] --chain-id [chain-id]",
		Short:   "Initializes a web socket which streams live events from a smart contract and relays them to the Cosmos network",
		Args:    cobra.ExactArgs(4),
		Example: "ebrelayer init tcp://localhost:26657 ws://127.0.0.1:7545/ 0x0823eFE0D0c6bd134a48cBd562fE4460aBE6e92c validator --chain-id=peggy",
		Run:     RunRelayerCmd,
	}

	return relayerCmd
}

// RunRelayerCmd executes the relayerCmd with the provided parameters
func RunRelayerCmd(cmd *cobra.Command, args []string) {
	// Load the validator's Ethereum private key
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		log.Fatal("invalid [ETHEREUM_PRIVATE_KEY] from .env")
	}

	// Parse chain's ID
	chainID := viper.GetString(client.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		log.Fatal("Must specify a 'chain-id'")
	}

	tendermintNode := args[0]

	// Parse ethereum provider
	ethereumProvider := args[1]
	if !relayer.IsWebsocketURL(ethereumProvider) {
		log.Fatalf("invalid [web3-provider]: %s", ethereumProvider)
	}

	// Parse the address of the deployed contract
	if !common.IsHexAddress(args[2]) {
		log.Fatalf("invalid [bridge-contract-address]: %s", args[2])
	}
	contractAddress := common.HexToAddress(args[2])

	// Parse the validator's moniker
	validatorFrom := args[3]

	// Parse Tendermint RPC URL
	rpcURL := viper.GetString(FlagRPCURL)

	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			log.Fatalf("invalid RPC URL: %v", rpcURL)
		}
	}

	// Load validator details
	validatorAddress, moniker, passphrase := utils.LoadValidatorCredentials(validatorFrom)

	// Load CLI context
	cliCtx := utils.LoadTendermintCLIContext(appCodec, validatorAddress, moniker, rpcURL, chainID)

	// Load Tx builder
	txBldr := authtypes.NewTxBuilderFromCLI().
		WithTxEncoder(sdkUtils.GetTxEncoder(appCodec)).
		WithChainID(chainID)

	// Start an Ethereum websocket
	go relayer.InitEthereumRelayer(
		appCodec,
		ethereumProvider,
		contractAddress,
		moniker,
		passphrase,
		validatorAddress,
		cliCtx,
		txBldr,
		privateKey,
	)

	// Start a Tendermint websocket
	go relayer.InitCosmosRelayer(
		tendermintNode,
		ethereumProvider,
		contractAddress,
		privateKey,
	)
}

func initConfig(cmd *cobra.Command) error {
	return viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
