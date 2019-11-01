package main

// 	Main (ebrelayer) : Implements CLI commands for the Relayer
//		service, such as initialization and event relay.

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"

	app "github.com/cosmos/peggy/app"
	relayer "github.com/cosmos/peggy/cmd/ebrelayer/relayer"
)

var appCodec *amino.Codec

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
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	initCmd.PersistentFlags().String("make-claims", "", "Make oracle claims everytime a prophecy claim is witnessed")

	// Construct Initialization Commands
	initCmd.AddCommand(
		ethereumRelayerCmd(),
		client.LineBreak,
		cosmosRelayerCmd(),
	)

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initCmd,
	)

	executor := cli.PrepareMainCmd(rootCmd, "EBRELAYER", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:          "ebrelayer",
	Short:        "Relayer service which listens for and relays ethereum smart contract events",
	SilenceUsage: true,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization subcommands",
}

//	ethereumRelayerCmd : Initializes a relayer service run by individual
//		validators which streams live events from an Ethereum smart contract.
//		The service automatically signs messages containing the event
//		data and relays them to tendermint for handling by the
//		EthBridge module.
//
func ethereumRelayerCmd() *cobra.Command {
	ethereumRelayerCmd := &cobra.Command{
		Use:   "ethereum [web3Provider] [contractAddress] [validatorFromName] --make-claims [make-claims] --chain-id [chain-id]",
		Short: "Initializes a web socket which streams live events from a smart contract and relays them to the Cosmos network",
		Args:  cobra.ExactArgs(3),
		// NOTE: Preface both parentheses in the event signature with a '\'
		Example: "ebrelayer init ethereum wss://ropsten.infura.io/ws 05d9758cb6b9d9761ecb8b2b48be7873efae15c0 validator --make-claims=false --chain-id=testing",
		RunE:    RunEthereumRelayerCmd,
	}

	return ethereumRelayerCmd
}

//	cosmosRelayerCmd : Initializes a Cosmos relayer service run by individual
//		validators which streams live events from the Cosmos network and then
//		relaying them to an Ethereum smart contract
//
func cosmosRelayerCmd() *cobra.Command {
	cosmosRelayerCmd := &cobra.Command{
		Use:     "cosmos [tendermintNode] [web3Provider] [bridgeContractAddress]",
		Short:   "Initializes a web socket which streams live events from the Cosmos network and relays them to the Ethereum network",
		Args:    cobra.ExactArgs(3),
		Example: "ebrelayer init cosmos tcp://localhost:26657 http://localhost:7545 0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2",
		RunE:    RunCosmosRelayerCmd,
	}

	return cosmosRelayerCmd
}

// RunEthereumRelayerCmd executes the initEthereumRelayerCmd with the provided parameters
func RunEthereumRelayerCmd(cmd *cobra.Command, args []string) error {
	// Parse chain's ID
	chainID := viper.GetString(client.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return fmt.Errorf("Must specify a 'chain-id'")
	}

	// Parse make claims boolean
	var makeClaims bool

	makeClaimsString := viper.GetString("make-claims")
	if strings.TrimSpace(makeClaimsString) == "true" {
		makeClaims = true
	} else {
		makeClaims = false
	}

	// Parse ethereum provider
	ethereumProvider := args[0]
	if !relayer.IsWebsocketURL(ethereumProvider) {
		return fmt.Errorf("Invalid web3-provider: %v", ethereumProvider)
	}

	// Parse the address of the deployed contract
	if !common.IsHexAddress(args[1]) {
		return fmt.Errorf("Invalid contract-address: %v", args[1])
	}
	contractAddress := common.HexToAddress(args[1])

	// Parse the validator's moniker
	validatorFrom := args[2]

	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(validatorFrom, false)
	if err != nil {
		return err
	}
	// Convert the validator's account address into type ValAddress
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Get the validator's passphrase using their moniker
	passphrase, err := keys.GetPassphrase(validatorFrom)
	if err != nil {
		return err
	}

	// Test passphrase is correct
	_, err = authtxb.MakeSignature(nil, validatorName, passphrase, authtxb.StdSignMsg{})
	if err != nil {
		return err
	}

	// Initialize the relayer
	err = relayer.InitEthereumRelayer(
		appCodec,
		chainID,
		ethereumProvider,
		contractAddress,
		makeClaims,
		validatorName,
		passphrase,
		validatorAddress)

	if err != nil {
		return err
	}

	return nil
}

// RunCosmosRelayerCmd executes the initCosmosRelayerCmd with the provided parameters
func RunCosmosRelayerCmd(cmd *cobra.Command, args []string) error {
	// Load config file containing environment variables
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	// Private key for validator's Ethereum address must be set as an environment variable
	privateKey := os.Getenv("ETHEREUM_PRIVATE_KEY")
	if strings.TrimSpace(privateKey) == "" {
		return fmt.Errorf("Error loading validator's private key from .env file")
	}

	// Tendermint node
	tendermintNode := args[0]

	// Ethereum websocket provider
	ethereumProvider := args[1]

	// Deployed contract address
	if !common.IsHexAddress(args[2]) {
		return fmt.Errorf("Invalid contract address: %v", args[2])
	}
	contractAddress := common.HexToAddress(args[2])

	// Initialize the relayer
	err = relayer.InitCosmosRelayer(
		tendermintNode,
		ethereumProvider,
		contractAddress,
		privateKey)

	if err != nil {
		return err
	}

	return nil
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
