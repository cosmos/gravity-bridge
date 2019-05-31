package events

// -------------------------------------------------------------
//      Event
//
//      Utility functions related to contract events, including
//       extration and parsing of an event's fields and values.
// --------------------------------------------------------------

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	log "github.com/golang/glog"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fatih/structs"
)

// ReturnEventFromABI returns abi.Event struct from the ABI
func ReturnEventFromABI(_abi abi.ABI, eventType string) (abi.Event, error) {
	event, ok := _abi.Events[eventType]
	if !ok {
		event, ok = _abi.Events[fmt.Sprintf("_%s", eventType)]
		if !ok {
			return abi.Event{}, errors.Errorf("no event type %v in contract", eventType)
		}
	}
	return event, nil
}

// NewEventFromContractEvent creates a new event after converting eventData to interface{}
func NewEventFromContractEvent(eventType string, contractName string, contractAddress common.Address,
	eventData interface{}, timestamp int64) (*Event, error) {
	event := &Event{}

	payload := NewEventPayload(eventData)

	logPayload, err := extractRawFieldFromEvent(payload)
	if err != nil {
		return event, err
	}
	// convert eventData to map[string]interface{}
	eventPayload, err := extractFieldsFromEvent(payload, eventData, eventType, contractName)
	if err != nil {
		return event, err
	}
	event, err = NewEvent(eventType, contractName, contractAddress, timestamp, eventPayload, logPayload)
	return event, err
}

// NewEvent is a convenience function to create a new Event
func NewEvent(eventType string, contractName string, contractAddress common.Address, timestamp int64,
	eventPayload map[string]interface{}, logPayload *types.Log) (*Event, error) {
	event := &Event{}
	event.eventType = eventType
	event.contractName = contractName
	event.contractAddress = contractAddress
	event.eventPayload = eventPayload
	event.logPayload = logPayload
	event.timestamp = timestamp
	event.eventHash = event.hashEvent()
	return event, nil
}

// Event represents a single smart contract event
// Abi generates event type, handled by filter and validator
type Event struct {

	// eventHash is the hash of event
	eventHash string

	// eventType is the type of event
	eventType string

	// contractAddress of the contract emitting the event
	contractAddress common.Address

	// contractName is the name of the contract
	contractName string

	// time from epoch of the block of the transaction for this event (in seconds)
	timestamp int64

	// event payload that doesn't include the "Raw" field
	eventPayload map[string]interface{}

	// "Raw" types.log field from event
	logPayload *types.Log
}

func extractFieldsFromEvent(payload *EventPayload, eventData interface{}, eventType string, contractName string) (map[string]interface{}, error) {
	eventPayload := make(map[string]interface{}, len(payload.data.Fields()))

	_abi, err := AbiJSON(contractName)
	if err != nil {
		return eventPayload, err
	}

	abiEvent, err := ReturnEventFromABI(_abi, eventType)
	if err != nil {
		return eventPayload, err
	}

	for _, input := range abiEvent.Inputs {
		eventFieldName := strings.Title(input.Name)
		eventField, ok := payload.Value(eventFieldName)
		if !ok {
			return eventPayload, errors.New("can't get event name in event")
		}
		switch input.Type.String() {
		case "address":
			addressVal, ok := eventField.Address()
			if !ok {
				return eventPayload, errors.New("could not convert to common.address type")
			}
			eventPayload[eventFieldName] = addressVal

		case "uint256":
			bigintVal, ok := eventField.BigInt()
			if !ok {
				return eventPayload, errors.New("could not convert to big.int")
			}
			eventPayload[eventFieldName] = bigintVal
		case "string":
			stringVal, ok := eventField.String()
			if !ok {
				return eventPayload, errors.New("could not convert to string")
			}
			eventPayload[eventFieldName] = stringVal
		case "bytes32":
			bytesVal, ok := eventField.Bytes32()
			if !ok {
				return eventPayload, errors.New("Could not convert to bytes32")
			}
			eventPayload[eventFieldName] = bytesVal
		default:
			return eventPayload, errors.Errorf("unsupported type encountered when parsing %v field for %v event %v",
				input.Type.String(), contractName, eventType)
		}
	}

	return eventPayload, nil
}

// AbiJSON returns parsed abi of this particular contract.
func AbiJSON(contractName string) (abi.ABI, error) {
	// contractType, ok := NameToContractTypes.GetFromContractName(contractName)
	// if !ok {
	// 	return abi.ABI{}, errors.New("contract name does not exist")
	// }
	// contractSpecs, ok := ContractTypeToSpecs.Get(contractType)
	// if !ok {
	// 	return abi.ABI{}, errors.New("invalid contract type")
	// }

	// TODO: PASS ABI
	_abi, err := abi.JSON(strings.NewReader("{ ABI }")) //contractSpecs.AbiStr())
	if err != nil {
		return abi.ABI{}, errors.New("cannot parse abi string")
	}
	return _abi, nil
}

func extractRawFieldFromEvent(payload *EventPayload) (*types.Log, error) {
	rawPayload, ok := payload.Value("Raw")
	if !ok {
		return &types.Log{}, errors.New("can't get raw value for event")
	}
	logPayload, ok := rawPayload.Log()
	if !ok {
		return &types.Log{}, errors.New("can't get log field of raw value for event")
	}
	return logPayload, nil
}

// hashEvent returns a hash for event using contractAddress, eventType, log index, and transaction hash
func (e *Event) hashEvent() string {
	logIndex := int(e.logPayload.Index)
	txHash := e.logPayload.TxHash.Hex()
	eventBytes, err := rlp.EncodeToBytes([]interface{}{
		e.contractAddress.Hex(),
		e.eventType, // nolint: gas, gosec
		strconv.Itoa(logIndex),
		txHash,
	})
	if err != nil {
		log.Errorf("Error encoding to bytes: err: %v", err)
		return ""
	}
	h := crypto.Keccak256Hash(eventBytes)
	return h.Hex()
}

// Hash returns the hash of the Event
func (e *Event) Hash() string {
	return e.eventHash
}

// EventType returns the eventType for the Event
func (e *Event) EventType() string {
	return e.eventType
}

// ContractAddress returns the contractAddress for the Event
func (e *Event) ContractAddress() common.Address {
	return e.contractAddress
}

// Timestamp returns the timestamp for the Event
func (e *Event) Timestamp() int64 {
	return e.timestamp
}

// SetTimestamp returns the timestamp for the Event
func (e *Event) SetTimestamp(ts int64) {
	e.timestamp = ts
}

// EventPayload returns the event payload for the Event
func (e *Event) EventPayload() map[string]interface{} {
	return e.eventPayload
}

// ContractName returns the contract name
func (e *Event) ContractName() string {
	return e.contractName
}

// LogPayload returns the log payload from the block
func (e *Event) LogPayload() *types.Log {
	// make a copy so fields are immutable
	logPayloadCopy := &types.Log{
		Address:     e.logPayload.Address,
		Topics:      e.logPayload.Topics,
		Data:        e.logPayload.Data,
		BlockNumber: e.logPayload.BlockNumber,
		TxHash:      e.logPayload.TxHash,
		TxIndex:     e.logPayload.TxIndex,
		BlockHash:   e.logPayload.BlockHash,
		Index:       e.logPayload.Index,
		Removed:     e.logPayload.Removed}
	return logPayloadCopy
}

// LogTopics returns the list of topics provided by the contract
func (e *Event) LogTopics() []common.Hash {
	return e.logPayload.Topics
}

// LogData is data provided by the contract, ABI encoded
func (e *Event) LogData() []byte {
	return e.logPayload.Data
}

// BlockNumber is the block number for this event
func (e *Event) BlockNumber() uint64 {
	return e.logPayload.BlockNumber
}

// TxHash is the hash of the transaction
func (e *Event) TxHash() common.Hash {
	return e.logPayload.TxHash
}

// TxIndex is the index of the transaction in the block
func (e *Event) TxIndex() uint {
	return e.logPayload.TxIndex
}

// BlockHash gets the block hash from the Event Payload
func (e *Event) BlockHash() common.Hash {
	return e.logPayload.BlockHash
}

// LogIndex is the log index position in the block
func (e *Event) LogIndex() uint {
	return e.logPayload.Index
}

// LogRemoved is true if log was reverted due to chain reorganization.
func (e *Event) LogRemoved() bool {
	return e.logPayload.Removed
}

// LogPayloadToString is a string representation of some fields of log
func (e *Event) LogPayloadToString() string {
	log := e.logPayload
	return fmt.Sprintf(
		"log: addr: %v, blknum: %v, txhash: %v, txidx: %v, blkhash: %v, idx: %v, rem: %v",
		log.Address.Hex(),
		log.BlockNumber,
		log.TxHash.Hex(),
		log.TxIndex,
		log.BlockHash.Hex(),
		log.Index,
		log.Removed,
	)
}

// EventPayload represents the data from a contract event
type EventPayload struct {

	// data is a Struct from the structs package. Just makes it easier
	// to handle access for any kind of event struct.
	data *structs.Struct
}

// NewEventPayload creates a new event payload
func NewEventPayload(eventData interface{}) *EventPayload {
	payload := &EventPayload{
		data: structs.New(eventData),
	}
	return payload
}

// Keys retrieves all the available key names in the event payload
func (p *EventPayload) Keys() []string {
	keyFields := p.data.Fields()
	keys := make([]string, len(keyFields))
	for ind, field := range keyFields {
		keys[ind] = field.Name()
	}
	return keys
}

// Value returns the EventPayloadValue of the given key
func (p *EventPayload) Value(key string) (*EventPayloadValue, bool) {
	field, ok := p.data.FieldOk(key)
	if !ok {
		return nil, ok
	}
	return &EventPayloadValue{value: field}, ok
}

// ToString returns a string representation for the payload
func (p *EventPayload) ToString() string {
	strs := []string{}
	for _, key := range p.Keys() {
		var str string
		val, _ := p.Value(key)
		if v, ok := val.Address(); ok {
			str = fmt.Sprintf("%v: %v", key, v.Hex())
		} else if v, ok := val.Log(); ok {
			str = fmt.Sprintf(
				"%v: addr: %v, blknum: %v, ind: %v, rem: %v",
				key,
				v.Address.Hex(),
				v.BlockNumber,
				v.Index,
				v.Removed,
			)
		} else if v, ok := val.BigInt(); ok {
			str = fmt.Sprintf("%v: %v", key, v)
		} else if v, ok := val.String(); ok {
			str = fmt.Sprintf("%v: %v", key, v)
		} else if v, ok := val.Int64(); ok {
			str = fmt.Sprintf("%v: %v", key, v)
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, "\n")
}

// EventPayloadValue represents a single value for a key in the payload
type EventPayloadValue struct {
	value *structs.Field
}

// Kind returns the value's basic type as described with reflect.Kind
func (v *EventPayloadValue) Kind() reflect.Kind {
	return v.value.Kind()
}

// Val returns the value as an unknown type interface{}
func (v *EventPayloadValue) Val() interface{} {
	return v.value.Value()
}

// String returns the value as a string
// Returns bool as false if unable to assert value as type string
func (v *EventPayloadValue) String() (string, bool) {
	val, ok := v.value.Value().(string)
	return val, ok
}

// Int64 returns the value as a int64.
// Returns bool as false if unable to assert value as type int64
func (v *EventPayloadValue) Int64() (int64, bool) {
	val, ok := v.BigInt()
	if !ok {
		return 0, ok
	}
	return val.Int64(), ok
}

// BigInt returns the value as a big.Int
// Returns bool as false if unable to assert value as type big.Int
func (v *EventPayloadValue) BigInt() (*big.Int, bool) {
	val, ok := v.value.Value().(*big.Int)
	return val, ok
}

// Bytes32 returns the value as a bytes32 object
// Returns bool as false if unable to assert value
func (v *EventPayloadValue) Bytes32() ([32]byte, bool) {
	val, ok := v.value.Value().([32]byte)
	return val, ok
}

// Address returns the value as common.Address
// Returns bool as false if unable to assert value as type common.Address
func (v *EventPayloadValue) Address() (common.Address, bool) {
	val, ok := v.value.Value().(common.Address)
	return val, ok
}

// Log returns the value as types.Log
// Returns bool as false if unable to assert value as type types.Log
func (v *EventPayloadValue) Log() (*types.Log, bool) {
	val, ok := v.value.Value().(types.Log)
	return &val, ok
}
