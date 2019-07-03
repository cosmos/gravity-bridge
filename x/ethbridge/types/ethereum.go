package types

import (
	gethCommon "github.com/ethereum/go-ethereum/common"
)

type EthereumAddress gethCommon.Address

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	return EthereumAddress(gethCommon.HexToAddress(address))
}

// Route should return the name of the module
func (ethereumAddress EthereumAddress) String() string {
	return gethCommon.Address(ethereumAddress).String()
}
