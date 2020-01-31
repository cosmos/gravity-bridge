package relayer

// -----------------------------------------------------
//      Ethereum relayer
//
//      Initializes the relayer service, which parses,
//      encodes, and packages named events on an Ethereum
//      Smart Contract for validator's to sign and send
//      to the Cosmos bridge.
// -----------------------------------------------------

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/contract"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitEthereumRelayer Starts an event listener on a specific Ethereum network, contract, and event
func InitEthereumRelayer(
	cdc *codec.Codec,
	chainID string,
	provider string,
	contractAddress common.Address,
	makeClaims bool,
	validatorName string,
	validatorAddress sdk.ValAddress,
	cliContext sdkContext.CLIContext,
	rpcURL string,
	privateKey *ecdsa.PrivateKey,
) error {
	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		return err
	}

	fmt.Printf("\nStarted Ethereum websocket with provider: %s", provider)

	var targetContract txs.ContractRegistry
	var eventName string

	// Load target contract and event name
	switch makeClaims {
	case true:
		targetContract = txs.CosmosBridge
		eventName = events.LogNewProphecyClaim.String()
	case false:
		targetContract = txs.BridgeBank
		eventName = events.LogLock.String()
	}
	// Load our contract's ABI
	contractABI := contract.LoadABI(makeClaims)

	// Load unique event signature from the named event contained within the contract's ABI
	eventSignature := contractABI.Events[eventName].Id().Hex()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get the specific contract's address (CosmosBridge or BridgeBank)
	targetAddress, err := txs.GetAddressFromBridgeRegistry(client, contractAddress, targetContract)
	if err != nil {
		log.Fatal(err)
	}

	// We need the target address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{targetAddress},
	}

	// We will check logs for new events
	logs := make(chan types.Log)

	// Filter by contract and event, write results to logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	fmt.Printf("\nSubscribed to %v contract at address: %s\n", targetContract, targetAddress.Hex())

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSignature {
				fmt.Println("\nWitnessed new event:", eventName)
				fmt.Println("Block number:", vLog.BlockNumber)
				fmt.Println("Tx hash:", vLog.TxHash.Hex())

				switch eventName {
				case events.LogLock.String():
					err := handleLogLockEvent(clientChainID, contractAddress, contractABI, eventName, vLog, chainID, cdc, validatorAddress, validatorName, cliContext, rpcURL)
					if err != nil {
						log.Fatal(err)
					}
				case events.LogNewProphecyClaim.String():
					err := handleLogNewProphecyClaimEvent(contractABI, eventName, vLog, provider, contractAddress, privateKey)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

// handleLogLockEvent unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func handleLogLockEvent(
	clientChainID *big.Int,
	contractAddress common.Address,
	contractABI abi.ABI,
	eventName string,
	log types.Log,
	chainID string,
	cdc *codec.Codec,
	validatorAddress sdk.ValAddress,
	validatorName string,
	cliContext sdkContext.CLIContext,
	rpcURL string,
) error {
	// Unpack the LogLock event using its unique event signature from the contract's ABI
	event := events.UnpackLogLock(clientChainID, contractAddress.Hex(), contractABI, eventName, log.Data)

	// Add the event to the record
	events.NewEventWrite(log.TxHash.Hex(), event)

	// Parse the LogLock event's payload into a struct
	prophecyClaim, err := txs.LogLockToEthBridgeClaim(validatorAddress, &event)
	if err != nil {
		return err
	}

	// Initiate the relay
	err = txs.RelayLockToCosmos(chainID, cdc, validatorAddress, validatorName, cliContext, &prophecyClaim, rpcURL)
	if err != nil {
		return err
	}

	return nil
}

// handleLogNewProphecyClaimEvent unpacks a LogNewProphecyClaim event, converts it to a OracleClaim, and relays a tx to Ethereum
func handleLogNewProphecyClaimEvent(
	contractABI abi.ABI,
	eventName string,
	log types.Log,
	provider string,
	contractAddress common.Address,
	privateKey *ecdsa.PrivateKey,
) error {
	// Unpack the LogNewProphecyClaim event using its unique event signature from the contract's ABI
	event := events.UnpackLogNewProphecyClaim(contractABI, eventName, log.Data)

	// Parse ProphecyClaim's data into an OracleClaim
	oracleClaim, err := txs.ProphecyClaimToSignedOracleClaim(event, privateKey)
	if err != nil {
		return err
	}

	// Initiate the relay
	return txs.RelayOracleClaimToEthereum(provider, contractAddress, events.LogNewProphecyClaim, oracleClaim, privateKey)
}
