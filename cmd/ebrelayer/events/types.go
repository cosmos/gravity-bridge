package events

// Event : enum containing supported contract events
type Event int

const (
	// Unsupported : unsupported Cosmos or Ethereum event
	Unsupported Event = iota
	// MsgBurn : Cosmos event 'CosmosMsg' type MsgBurn
	MsgBurn
	// MsgLock :  Cosmos event 'CosmosMsg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogNewProphecyClaim : Ethereum event 'NewProphecyClaimEvent'
	LogNewProphecyClaim
)

// String : returns the event type as a string
func (d Event) String() string {
	return [...]string{"unsupported", "burn", "lock", "LogLock", "LogNewProphecyClaim"}[d]
}

// CosmosMsgAttributeKey : enum containing supported attribute keys
type CosmosMsgAttributeKey int

const (
	// UnsupportedAttributeKey : unsupported attribute key
	UnsupportedAttributeKey CosmosMsgAttributeKey = iota
	// CosmosSender : sender's address on Cosmos network
	CosmosSender
	// EthereumReceiver : receiver's address on Ethereum network
	EthereumReceiver
	// Coin : coin type
	Coin
	// TokenContractAddress : coin's corresponding contract address deployed on the Ethereum network
	TokenContractAddress
)

// String : returns the event type as a string
func (d CosmosMsgAttributeKey) String() string {
	return [...]string{"unsupported", "cosmos_sender", "ethereum_receiver", "amount", "token_contract_address"}[d]
}
