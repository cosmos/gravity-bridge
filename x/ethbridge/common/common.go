package common

import gethCommon "github.com/ethereum/go-ethereum/common"

//IsValidEthereumAddress returns true if address is valid
func IsValidEthAddress(s string) bool {
	return gethCommon.IsHexAddress(s)
}
