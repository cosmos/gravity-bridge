package main

// -------------------------------------------------------------
//      Main (ebrelayer)
//
//      Implements CLI commands for the Relayer service, such as
//      initalization and event relay.
// -------------------------------------------------------------

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	// "golang.org/x/crypto"

	app "github.com/swishlabsco/cosmos-ethereum-bridge"
	events "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	relayer "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/relayer"
)

const (
	storeAcc       = "acc"
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

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		getCheckCmd(),
		getPrintCmd(),
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

func initRelayerCmd() *cobra.Command {
	initRelayerCmd := &cobra.Command{
		Use:   "init chain-id web3-provider contract-address event-signature validatorFromName",
		Short: "Initalizes a web socket which streams live events from a smart contract",
		RunE:  RunRelayerCmd,
	}

	return initRelayerCmd
}

func getPrintCmd() *cobra.Command {
	getPrintCmd := &cobra.Command{
		Use:   "print",
		Short: "Print all stored events",
		RunE:  RunPrintCmd,
	}

	return getPrintCmd
}

func getCheckCmd() *cobra.Command {
	getCheckCmd := &cobra.Command{
		Use:   "check tx-hash",
		Short: "Check a transaction hash to see if it is available for claim submission",
		RunE:  RunCheckCmd,
	}

	return getCheckCmd
}

// -------------------------------------------------------------------------------------
//  `ebrelayer init "testing" "wss://ropsten.infura.io/ws" "3de4ef81Ba6243A60B0a32d3BCeD4173b6EA02bb"
//	 "LogLock(bytes32,address,bytes,address,uint256,uint256)" "cosmos13mztulrrz3leephsr6dhxker4t68qxew9m9nhn"`
// -------------------------------------------------------------------------------------

func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 5 {
		return fmt.Errorf("Expected 5 arguments, got %v", len(args))
	}

	// Parse chain's ID
	chainId := args[0]
	if chainId == "" {
		return fmt.Errorf("Invalid chain-id: %v", chainId)
	}

	// Parse ethereum provider
	ethereumProvider := args[1]
	if !relayer.IsWebsocketURL(ethereumProvider) {
		return fmt.Errorf("Invalid web3-provider: %v", ethereumProvider)
	}

	// Parse the address of the deployed contract
	bytesContractAddress, err := hex.DecodeString(args[2])
	if err != nil {
		return fmt.Errorf("Invalid contract-address: %v", bytesContractAddress)
	}
	contractAddress := common.BytesToAddress(bytesContractAddress)

	// Parse the event signature for the subscription
	eventSig := "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	// eventSig := crypto.Keccak256Hash([]byte(args[3]))
	if eventSig == "" {
		return fmt.Errorf("Invalid event-signature: %v", eventSig)
	}

	// Parse the validator running the relayer service
	validatorFrom := args[4]

	// Initialize the relayer
	initErr := relayer.InitRelayer(
		appCodec,
		chainId,
		ethereumProvider,
		contractAddress,
		eventSig,
		validatorFrom)

	if initErr != nil {
		fmt.Printf("%v", initErr)
		return initErr
	}

	return nil
}

func RunCheckCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected tx-hash")
	}
	txHash := args[0]
	events.PrintEventByTx(txHash)

	return nil
}

func RunPrintCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Expected 0 args")
	}
	events.PrintEvents()
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
