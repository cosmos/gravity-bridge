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
	"fmt"
	"log"

	amino "github.com/tendermint/go-amino"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/contract"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitEthereumRelayer : Starts an event listener on a specific Ethereum network, contract, and event
func InitEthereumRelayer(cdc *amino.Codec, chainID string, provider string, contractAddress common.Address, cosmosSupport bool, validatorName string, passphrase string, validatorAddress sdk.ValAddress) error {
	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		return err
	}
	fmt.Printf("\nStarted ethereum websocket with provider: %s", provider)

	// Load our contract's ABI
	contractABI := contract.LoadABI(cosmosSupport)

	var targetContract txs.ContractRegistry
	var eventName string
	var eventSignature string

	switch cosmosSupport {
	case true:
		targetContract = txs.CosmosBridge
		eventName = "LogNewProphecyClaim"
		eventSignature = contractABI.Events["LogNewProphecyClaim"].Id().String()
	case false:
		targetContract = txs.BridgeBank
		eventName = "LogLock"
		eventSignature = contractABI.Events["LogLock"].Id().Hex()
	}

	// Get the specific contract's address (Valset, Oracle, CosmosBridge, or BridgeBank)
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
			// TODO: Remove this log
			fmt.Printf("\n\nTx hash: %v\nTopics: %v",
				vLog.TxHash.Hex(), vLog.Topics[0].Hex())
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSignature {
				fmt.Printf("\n\nNew \"%v\":\nTx hash: %v\nBlock number: %v",
					eventName, vLog.TxHash.Hex(), vLog.BlockNumber)

				switch eventName {
				case events.LogLock.String():
					event := events.UnpackLogLock(contractABI, eventName, vLog.Data)

					// Add the event to the record
					events.NewEventWrite(vLog.TxHash.Hex(), event)

					// Parse the LogLock event's payload into a struct
					claim, err := txs.ParseLogLockPayload(validatorAddress, &event)
					if err != nil {
						return err
					}

					// Initiate the relay
					err = txs.RelayLockToCosmos(chainID, cdc, validatorAddress, validatorName, passphrase, &claim)
					if err != nil {
						return err
					}
				case events.LogNewProphecyClaim.String():
					event := events.UnpackLogNewProphecyClaim(contractABI, eventName, vLog.Data)

					err = txs.RelayOracleClaimToEthereum(provider, contractAddress, events.LogNewProphecyClaim, event)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
