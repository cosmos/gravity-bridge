package common

import (
	gethCommon "github.com/ethereum/go-ethereum/common"
)

type EthereumAddress string

//IsValidEthereumAddress returns true if address is valid
func IsValidEthAddress(address EthereumAddress) bool {
	return gethCommon.IsHexAddress(string(address))
}
