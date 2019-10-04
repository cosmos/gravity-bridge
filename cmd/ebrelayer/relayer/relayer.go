package relayer

// -----------------------------------------------------
//      Relayer
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

// InitRelayer : Starts an event listener on a specific network, contract, and event
func InitRelayer(cdc *amino.Codec, chainId string, provider string, contractAddress common.Address, eventSig string, validatorName string, passphrase string, validatorAddress sdk.ValAddress, rpcURL string) error {

	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		return err
	}
	fmt.Printf("\nStarted ethereum websocket with provider: %s", provider)

	// We need the contract address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// We will check logs for new events
	logs := make(chan types.Log)

	// Filter by contract and event, write results to logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	fmt.Printf("\nSubscribed to contract events on address: %s\n", contractAddress.Hex())

	// Load Peggy Contract's ABI
	contractABI := contract.LoadABI()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSig {
				fmt.Printf("\n\nNew Lock Transaction:\nTx hash: %v\nBlock number: %v\n",
					vLog.TxHash.Hex(), vLog.BlockNumber)

				// Parse the event data into a new LockEvent using the contract's ABI
				event := events.NewLockEvent(contractABI, clientChainID, contractAddress.Hex(), "LogLock", vLog.Data)

				// Add the event to the record
				events.NewEventWrite(vLog.TxHash.Hex(), event)
				// Parse the event's payload into a struct
				claim, err := txs.ParsePayload(validatorAddress, &event)
				if err != nil {
					return err
				}

				// Initiate the relay
				err = txs.RelayEvent(chainId, cdc, validatorAddress, validatorName, passphrase, &claim, rpcURL)
				if err != nil {
					return err
				}
			}
		}
	}
}
