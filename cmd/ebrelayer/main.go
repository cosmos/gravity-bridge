package main

// -------------------------------------------------------------
//      Main (ebrelayer)
//
//      Implements CLI commands for the Relayer service, such as
//      initalization and event relay.
// -------------------------------------------------------------

import (
	"fmt"
	"os"
	"encoding/hex"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/ethereum/go-ethereum/common"
	// "golang.org/x/crypto"

	app "github.com/swishlabsco/cosmos-ethereum-bridge"
	relayer "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/relayer"
	events "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
)

const (
	storeAcc       = "acc"
	routeEthbridge = "ethbridge"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.ebcli")
var appCodec *amino.Codec
// var keybase *keys.Keybase

func init() {

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	cdc := app.MakeCodec()
	appCodec = cdc

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		getClaimsCmd(),
		getAccountCmd(),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
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

func getClaimsCmd() *cobra.Command {
	getClaimsCmd := &cobra.Command{
		Use:   "check event-id",
		Short: "Prints historical claim information about the event",
		RunE:  RunClaimCmd,
	}

	return getClaimsCmd
}

func getAccountCmd() *cobra.Command {
	getAccountCmd := &cobra.Command{
		Use:   "account unique-account",
		Short: "Does this account exist?",
		RunE:  RunAccountCmd,
	}

	return getAccountCmd
}

func initRelayerCmd() *cobra.Command {
	initRelayerCmd := &cobra.Command{
		Use:   "init chain-id web3-provider contract-address event-signature validator",
		Short: "Initalizes a web socket which streams live events from a smart contract",
		RunE:  RunRelayerCmd,
	}

	return initRelayerCmd
}

// -------------------------------------------------------------------------------------
//  `ebrelayer init "testing" "wss://ropsten.infura.io/ws" "3de4ef81Ba6243A60B0a32d3BCeD4173b6EA02bb"
//	 "LogLock(bytes32,address,bytes,address,uint256,uint256)" "cosmos13mztulrrz3leephsr6dhxker4t68qxew9m9nhn"`
// -------------------------------------------------------------------------------------

func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 5 {
		return fmt.Errorf("Expected 5 arguments, got %s", len(args))
	}

	// Parse chain's ID
	chainId := args[0]
	if chainId == "" {
		return fmt.Errorf("Invalid chain-id: %s", chainId)
	}

	// Parse ethereum provider
	ethereumProvider := args[1]
	if !relayer.IsWebsocketURL(ethereumProvider) {
		return fmt.Errorf("Invalid web3-provider: %s", ethereumProvider)
	}

	// Parse the address of the deployed contract
	bytesContractAddress, err := hex.DecodeString(args[2])
	if err != nil {
		return fmt.Errorf("Invalid contract-address: %s", bytesContractAddress, err)
	}
	contractAddress := common.BytesToAddress(bytesContractAddress)

	// Parse the event signature for the subscription
	eventSig := "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	// eventSig := crypto.Keccak256Hash([]byte(args[3]))
	if eventSig == "" {
		return fmt.Errorf("Invalid event-signature: %s", eventSig)
	}

	// Parse the validator running the relayer service
	validator, valErr := sdk.AccAddressFromBech32(args[4])
	if valErr != nil {
		return fmt.Errorf("Invalid validator: %s", validator)
	}

	// Initialize the relayer
	initErr := relayer.InitRelayer(
		appCodec,
		chainId,
		ethereumProvider,
		contractAddress,
		eventSig,
		validator)

	if initErr != nil {
		fmt.Printf("%s", initErr)
	}

	return nil
}

func RunClaimCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected event-id argument")
	}

	eventId := args[0]

	// TODO: differentiate between an invalid event and an event with 0 claims
	if !events.IsStoredEvent(eventId) {
		return fmt.Errorf("Invalid event-id: %s", eventId)
	}

	events.PrintClaims(eventId)

	return nil
}

func RunAccountCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected account argument")
	}

	account, valErr := sdk.AccAddressFromBech32(args[0])
	if valErr != nil {
		return fmt.Errorf("Invalid account: %s", account)
	}

	return fmt.Errorf("Success!")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
