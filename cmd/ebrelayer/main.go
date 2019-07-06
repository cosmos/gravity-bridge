package main

// 	Main (ebrelayer) : Implements CLI commands for the Relayer
//		service, such as initalization and event relay.

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	app "github.com/cosmos/peggy"
	relayer "github.com/cosmos/peggy/cmd/ebrelayer/relayer"
)

const (
	routeEthbridge = "ethbridge"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.ebcli")
var appCodec *amino.Codec

func init() {

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	cdc := app.MakeCodec()
	appCodec = cdc

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "EBRELAYER", defaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:          "ebrelayer",
	Short:        "Relayer service which listens for and relays ethereum smart contract events",
	SilenceUsage: true,
}

//	initRelayerCmd : Initializes a relayer service run by individual
//		validators which streams live events from a smart contract.
//		The service automatically signs messages containing the event
//		data and relays them to tendermint for handling by the
//		EthBridge module.
//
func initRelayerCmd() *cobra.Command {
	initRelayerCmd := &cobra.Command{
		Use:   "init [web3Provider] [contractAddress] [eventSignature] [validatorFromName] --chain-id [chain-id]",
		Short: "Initalizes a web socket which streams live events from a smart contract",
		Args:  cobra.ExactArgs(4),
		// NOTE: Preface both parentheses in the event signature with a '\'
		Example: "ebrelayer init wss://ropsten.infura.io/ws ec6df30846baab06fce9b1721608853193913c19 LogLock(bytes32,address,bytes,address,uint256,uint256) validator --chain-id=testing",
		RunE:    RunRelayerCmd,
	}

	return initRelayerCmd
}

//	RunRelayerCmd : RunRelayerCmd executes the initRelayerCmd with the provided parameters
func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	// Parse chain's ID
	chainID := viper.GetString(client.FlagChainID)
	if chainID == "" {
		return fmt.Errorf("Must specify a 'chain-id'")
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

	// Convert event signature to []bytes and apply the Keccak256Hash
	eventSigHash := crypto.Keccak256Hash([]byte(args[2]))

	// Get the hex event signature from the hash.
	// Expected value: "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	eventSig := eventSigHash.Hex()

	// Parse the validator's moniker
	validatorFrom := args[3]

	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(validatorFrom)
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
	err = relayer.InitRelayer(
		appCodec,
		chainID,
		ethereumProvider,
		contractAddress,
		eventSig,
		validatorName,
		passphrase,
		validatorAddress)

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
