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
	"io/ioutil"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	// "golang.org/x/crypto/sha3"
	// "golang.org/x/crypto"

	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/txs"
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
	logLockEvent := "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	// TODO: Generate logLockEvent hash using 'crypto' library instead of hard coding,
	// 			 i.e. `logLockEvent := crypto.Keccak256Hash(logLockSig)`

	fmt.Printf("\nContract Address: %s\n", contractAddress.Hex())
	fmt.Printf("LogLockEvent Signature: %s\n", logLockSig)
	fmt.Printf("LogLockEvent: %s\n", logLockEvent)

	// We need the contract address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	// Subscribe to the web socket, filter by contract and event, write results to logs
	fmt.Printf("\nStarting event listener on address: %s.", contractAddress.Hex())
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Event listener started on address: %s.\n", contractAddress.Hex())

	// Open the file containing Peggy Contract's ABI
	rawContractAbi, errorMsg := ioutil.ReadFile("cmd/ebrelayer/contract/PeggyABI.json")
	if errorMsg != nil {
		log.Fatal(errorMsg)
	}

	// Convert the raw abi into a usable format
	contractAbi, err := abi.JSON(strings.NewReader(string(rawContractAbi)))

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			fmt.Println("\nNew event:")
			fmt.Println("BlockHash: ", vLog.BlockHash.Hex())
			fmt.Println("BlockNumber: ", vLog.BlockNumber)
			fmt.Println("TxHash: ", vLog.TxHash.Hex())

			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logLockEvent {

				// Parse the event's attributes as Ethereum network variables
				event := events.LockEvent{}
				err := contractAbi.Unpack(&event, "LogLock", vLog.Data)
				if err != nil {
					log.Fatal("Unpacking: ", err)
				}

				id := hex.EncodeToString(event.Id[:])
				sender := event.From.Hex()
				recipient := string(event.To[:])
				token := event.Token.Hex()
				value := event.Value
				nonce := event.Nonce

				fmt.Println("\nLogLock data:")
				fmt.Println("Event ID: ", id)
				fmt.Println("Token : ", token)
				fmt.Println("Sender : ", sender)
				fmt.Println("Recipient : ", recipient)
				fmt.Println("Value : ", value)
				fmt.Println("Nonce : ", nonce)

				// Add the witnessing validator to the event
				claimCount := events.ValidatorMakeClaim(id, validator)
				fmt.Println("Total claims on this event: ", claimCount)

				// Parse the event's payload into a struct and initiate the relay
				txErr := txs.ParsePayloadAndRelay(validator, &event) //cdc,
				if txErr != "" {
					fmt.Printf("txErr")
				}
			}
		}
	}
	return fmt.Errorf("Error: Relayer timed out.")
}
