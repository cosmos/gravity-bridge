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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"

	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"

	"github.com/cosmos/cosmos-sdk/codec"
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
	cdc *codec.Codec,
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
	fmt.Printf("validator: %s\n\n", validator)

	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Deployed contract address and event signature
	b, err := hex.DecodeString(peggyContractAddress)
	if err != nil {
		return fmt.Errorf("Error while decoding contract address")
	}

	contractAddress := common.HexToAddress(peggyContractAddress)
	logLockSig := []byte(eventSignature)
	logLockEvent := sha3.Keccak256Hash(logLockSig)

	fmt.Printf("\n\nContract Address: %s\n Log Lock Signature: %s\n\n",
		b, logLockSig)

	fmt.Printf("%s", logLockEvent)

	// TODO: resolve type casting error between go-ethereum/common and swish/go-ethereum/common
	// Filter currently captures all events from the contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	// Subscribe to the client, filter based on query, write events to logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
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
			fmt.Println("\nBlock Number:", vLog.BlockNumber)

			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logLockEvent.Hex() {

				// Current time is in system time, will be updated to block time
				currentTime := fmt.Println(time.Now().Format(time.RFC850))

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
				fmt.Printf(event.ToString())

				// Add the witness
				errWitness := events.ValidatorMakeClaim(event.payload.Value("_id"))
				if errWitness != nil {
					log.Fatal(errWitness)
				}

				// Parse the event's payload into a golang struct and initiate the relay
				result, txErr := txs.parsePayloadAndRelay(cdc, event.eventPayload)
				if txErr != nil {
					log.Fatal(txErr)
				}
			}
		}
	}
	return fmt.Errorf("Error: Relayer timed out.")
}
