package events

// -----------------------------------------------------
//    Event
//
// -----------------------------------------------------

import (
	// "encoding/hex"
	// "fmt"
	// "log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"

	// "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/contract"
)

var ApplicationBinaryInterface abi.ABI

// LockEvent represents a single smart contract event
type LockEvent struct {
	Id      [32]byte
	From    common.Address
	To      []byte
	Token   common.Address
	Value   *big.Int
	Nonce   *big.Int
}

func NewEvent(eventName string, eventData []byte) string {

	// Load Peggy smart contract abi
	// if eventName == "LogLock" {
		// ApplicationBinaryInterface := contract.LoadABI()
	// } else {
	// 	return "Only LogLock events are currently supported."
	// }

	// Parse the event's attributes as Ethereum network variables
	// event := LockEvent{}
	// err := ApplicationBinaryInterface.Unpack(&event, eventName, eventData)
	// if err != nil {
	//     log.Fatal("Unpacking: ", err)
	// }

	// // Convert the variables into a printable format
	// id := hex.EncodeToString(event.Id[:])
	// sender := event.From.Hex()
	// recipient := string(event.To[:])
	// token := event.Token.Hex()
	// value := event.Value
	// nonce := event.Nonce

	// // Print the event's information
	// fmt.Println("\n Event data:")
	// fmt.Println("Event ID: ", id)
	// fmt.Println("Token : ", token)
	// fmt.Println("Sender : ", sender)
	// fmt.Println("Recipient : ", recipient)
	// fmt.Println("Value : ", value)
	// fmt.Println("Nonce : ", nonce)

	// return id
	return eventName
}


