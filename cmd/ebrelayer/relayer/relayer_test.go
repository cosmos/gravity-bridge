package relayer

// ------------------------------------------------------------
//    Relayer_Test
//
//    Tests Relayer functionality.
//
// ------------------------------------------------------------

import (
	"testing"
)

const (
	ChainID          = "testing"
	Client           = "wss://ropsten.infura.io/ws"
	ContractAddress  =  "3de4ef81Ba6243A60B0a32d3BCeD4173b6EA02bb"
	// EventSig is hash of "LogLock(bytes32,address,bytes,address,uint256,uint256)"
	EventSig         = "0xe154a56f2d306d5bbe4ac2379cb0cfc906b23685047a2bd2f5f0a0e810888f72"
	Validator        = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"

)

func StartRelayer(t *testing.T) err {
	err = go(InitRelayer(ChainID, Client, ContractAddress, EventSig, Validator))
}

