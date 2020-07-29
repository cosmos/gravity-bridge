package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "nameservice"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	EthAddressKey    = []byte{0x1}
	ValsetRequestKey = []byte{0x2}
	ValsetConfirmKey = []byte{0x3}
)

func GetEthAddressKey(validator sdk.AccAddress) []byte {
	return append(EthAddressKey, []byte(validator)...)
}

func GetValsetRequestKey(nonce int64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))

	return append(ValsetRequestKey, nonceBytes...)
}

func GetValsetConfirmKey(nonce int64, validator sdk.AccAddress) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))

	return append(ValsetConfirmKey, append(nonceBytes, []byte(validator)...)...)
}
