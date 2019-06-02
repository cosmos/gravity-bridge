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
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/ethereum/go-ethereum/common"

	relayer "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/relayer"
	events "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
)

const (
	storeAcc       = "acc"
	routeEthbridge = "ethbridge"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.ebrelayer")

func init() {

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		getClaimsCmd(),
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
//	 "LogLock(bytes32,address,bytes,address,uint256,uint256)" "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"`
// -------------------------------------------------------------------------------------

func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 5 {
		return fmt.Errorf("Expected 5 arguments, got ", len(args))
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
	// TODO: Generate eventSig hash using 'crypto' library instead of hard coding
	// `eventSig := crypto.Keccak256Hash(args[3])`
	if eventSig == "" {
		return fmt.Errorf("Invalid event-signature: %s", eventSig)
	}

	// Parse the validator running the relayer service
	validator := sdk.AccAddress(args[4])
	if validator == nil {
		return fmt.Errorf("Invalid validator: %s", validator)
	}

	// Initialize the relayer
	initErr := relayer.InitRelayer(
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

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
