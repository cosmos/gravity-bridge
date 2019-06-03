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

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/contract"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/txs"
)

// -------------------------------------------------------------------------
// Starts an event listener on a specific network, contract, and event
// -------------------------------------------------------------------------

func InitRelayer(cdc *amino.Codec, chainId string, provider string,
	contractAddress common.Address, eventSig string,
	validatorFrom string) error {

	validatorAddress, validatorName, err := sdkContext.GetFromFields(validatorFrom)
	if err != nil {
		fmt.Printf("failed to get from fields: %v", err)
		return err
	}

	passphrase, err := keys.GetPassphrase(validatorFrom)
	if err != nil {
		return err
	}

	//Test passhprase is correct
	_, err = authtxb.MakeSignature(nil, validatorName, passphrase, authtxb.StdSignMsg{})
	if err != nil {
		fmt.Printf("passphrase error: %v", err)
		return err
	}

	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		fmt.Errorf("%s", err)
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
		fmt.Errorf("%s", err)
	} else {
		fmt.Printf("\nSubscribed to contract events on address: %s\n", contractAddress.Hex())
	}

	// Load Peggy Contract's ABI
	contractABI := contract.LoadABI()

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSig {
				fmt.Printf("\n\nNew Lock Transaction:\nTx hash: %v\nBlock number: %v",
					vLog.TxHash.Hex(), vLog.BlockNumber)

				// Parse the event data into a new LockEvent using the contract's ABI
				event := events.NewLockEvent(contractABI, "LogLock", vLog.Data)

				// Add the event to the record
				successfulStore := events.NewEventWrite(vLog.TxHash.Hex(), event)
				if successfulStore != true {
					fmt.Errorf("Error: event not stored")
				}

				// Parse the event's payload into a struct
				claim, claimErr := txs.ParsePayload(validatorAddress, &event)
				if claimErr != nil {
					fmt.Errorf("Error: %s", claimErr)
				}

				// Initiate the relay
				relayErr := txs.RelayEvent(chainId, cdc, validatorAddress, validatorName, passphrase, &claim)
				if relayErr != nil {
					fmt.Errorf("Error: %s", relayErr)
				}
			}
		}
	}
	return fmt.Errorf("Error: Relayer timed out.")
}
