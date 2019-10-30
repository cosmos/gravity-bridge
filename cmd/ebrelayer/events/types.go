package events

// Event : enum containing supported contract events
type Event int

const (
	// MsgBurn : Cosmos event 'CosmosMsg' type MsgBurn
	MsgBurn Event = iota
	// MsgLock :  Cosmos event 'CosmosMsg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogNewProphecyClaim : Ethereum event 'NewProphecyClaimEvent'
	LogNewProphecyClaim
	// Unsupported : unsupported Cosmos or Ethereum event
	Unsupported
)

// String : returns the event type as a string
func (d Event) String() string {
	return [...]string{"burn", "lock", "LogLock", "LogNewProphecyClaim", "unsupported"}[d]
}
