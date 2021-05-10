package types

import tmbytes "github.com/tendermint/tendermint/libs/bytes"

var (
	_ EthereumEvent = &SendToCosmosEvent{}
	_ EthereumEvent = &BatchExecutedEvent{}
	_ EthereumEvent = &ContractCallExecutedEvent{}
	_ EthereumEvent = &ERC20DeployedEvent{}
)

//////////////
// GetNonce //
//////////////

func (stce *SendToCosmosEvent) GetNonce() uint64 {
	return stce.EventNonce
}

func (bee *BatchExecutedEvent) GetNonce() uint64 {
	return bee.EventNonce
}

func (ccee *ContractCallExecutedEvent) GetNonce() uint64 {
	return ccee.EventNonce
}

func (e20de *ERC20DeployedEvent) GetNonce() uint64 {
	return e20de.EventNonce
}

//////////
// Hash //
//////////

func (stce *SendToCosmosEvent) Hash() tmbytes.HexBytes {
	panic("NOT IMPLEMENTED")
}

func (bee *BatchExecutedEvent) Hash() tmbytes.HexBytes {
	panic("NOT IMPLEMENTED")
}

func (ccee *ContractCallExecutedEvent) Hash() tmbytes.HexBytes {
	panic("NOT IMPLEMENTED")
}

func (e20de *ERC20DeployedEvent) Hash() tmbytes.HexBytes {
	panic("NOT IMPLEMENTED")
}

//////////////
// Validate //
//////////////


func (stce *SendToCosmosEvent) Validate() error {
	panic("NOT IMPLEMENTED")
}

func (bee *BatchExecutedEvent) Validate() error {
	panic("NOT IMPLEMENTED")
}

func (ccee *ContractCallExecutedEvent) Validate() error {
	panic("NOT IMPLEMENTED")
}

func (e20de *ERC20DeployedEvent) Validate() error {
	panic("NOT IMPLEMENTED")
}
