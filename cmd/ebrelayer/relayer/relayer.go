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
	"encoding/hex"
	"fmt"
	"log"
	// "time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	// "golang.org/x/crypto/sha3"
	// "golang.org/x/crypto"

	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	// "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/txs"

	// "github.com/cosmos/cosmos-sdk/codec"
)

type WitnessedLogLock struct {
	EthereumSender  string         `json:"ethereum_sender"`
	CosmosRecipient sdk.AccAddress `json:"cosmos_receiver"`
	Validator       sdk.AccAddress `json:"validator"`
	Amount          sdk.Coins      `json:"amount"`
	Nonce           int            `json:"nonce"`
}

// -------------------------------------------------------------------------
// Starts an event listener on a specific network, contract, and event
// -------------------------------------------------------------------------

func InitRelayer(
	// cdc *codec.Codec,
	chainId string,
	provider string,
	peggyContractAddress string,
	eventSignature string,
	validator sdk.AccAddress) error {

	// Console log for testing purposes...
	fmt.Printf("\n\ninitRelayer() received params:\n")
	fmt.Printf("chainId: %s\n", chainId)
	fmt.Printf("provider: %s\n", provider)
	fmt.Printf("peggyContractAddress: %s\n", peggyContractAddress)
	fmt.Printf("eventSignature: %s\n", eventSignature)
	fmt.Printf("validator: %s\n", validator)

	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Deployed contract address and event signature
	bytesContractAddress, err := hex.DecodeString(peggyContractAddress)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	contractAddress := common.BytesToAddress(bytesContractAddress)
	logLockSig := []byte(eventSignature)
	//logLockEvent := sha3.Keccak256Hash(logLockSig) //crypto
	logLockEvent := "3e43256e124a7860d7fd775c424fb6eb9e1988b31b374011335beb396e201d90"

	fmt.Printf("\nContract Address: %s\n", contractAddress.Hex())
	fmt.Printf("LogLockEvent Signature: %s\n", logLockSig)
	fmt.Printf("LogLockEvent: %s\n", logLockEvent)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	// Subscribe to the client, filter based on query, write events to logs
	fmt.Printf("\nStarting event listener on address: %s...\n", contractAddress.Hex())
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nEvent listener started on address: %s!\n", contractAddress.Hex())

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			fmt.Println("\nEvent wittnessed at block number: %d", vLog.BlockNumber)

			// Check if the event is a 'LogLock' event
			fmt.Printf("\nEvent topic HEX: %s", vLog.Topics[0].Hex())
			fmt.Printf("\nEvent topic: %s", vLog.Topics[0])
			fmt.Printf("\nLock event: %s", logLockEvent)
			// if vLog.Topics[0].Hex() == logLockEvent.Hex() {
			if vLog.Topics[0].Hex() == logLockEvent {

				// Current time is in system time, will be updated to block time
				//TODO:
				var currentTime int64 = 1111;
				// currentTime, errPrint := fmt.Printf(time.Now().Format(time.RFC850))
				// if errPrint != nil {
				// 	log.Fatal(errPrint)
				// }

				// Parse contract event data into package
				event, eventErr := events.NewEventFromContractEvent(
					"LogLock",
					"PeggyContract",
					contractAddress,
					vLog,
					currentTime,
				)
				if eventErr != nil {
					log.Fatal(eventErr)
				}

				// Print event data to console
				// event.EventPayload().Keys()

				eventId, exists := event.EventPayload()["_id"].(string)
				if exists != true {
					fmt.Printf("event _id does not exist in payload")
				}

				// Add the witnessing validator to the event
				events.ValidatorMakeClaim(eventId, validator)

				// Parse the event's payload into a golang struct and initiate the relay
				// result, txErr := txs.parsePayloadAndRelay(cdc, validator, event)
				// if txErr != nil {
				// 	log.Fatal(txErr)
				// }
			}
		}
	}
	return fmt.Errorf("Error: Relayer timed out.")
}
