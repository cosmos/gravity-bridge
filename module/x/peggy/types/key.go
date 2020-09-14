package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "peggy"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	EthAddressKey        = []byte{0x1}
	ValsetRequestKey     = []byte{0x2}
	ValsetConfirmKey     = []byte{0x3}
	OracleClaimKey       = []byte{0x4}
	OracleAttestationKey = []byte{0x5}
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

// GetClaimKey creates a key with claim type and address first as they have a fix length.
// We can use this later for prefix scans to find all claims by type,
// claims by validator or claims for a nonce
func GetClaimKey(claimType ClaimType, nonce Nonce, validator sdk.ValAddress) []byte {
	key := make([]byte, len(OracleClaimKey)+ClaimTypeLen+sdk.AddrLen+len(nonce))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+ClaimTypeLen:], validator)
	copy(key[len(OracleClaimKey)+ClaimTypeLen+sdk.AddrLen:], nonce)
	return key
}
func SplitClaimKey(raw []byte) (ClaimType, sdk.ValAddress, Nonce) {
	return ClaimType(raw[1 : 1+ClaimTypeLen][0]), raw[1+ClaimTypeLen : 1+ClaimTypeLen+sdk.AddrLen], raw[1+ClaimTypeLen+sdk.AddrLen:]
}

func GetAttestationKey(claimType ClaimType, nonce Nonce) []byte {
	key := make([]byte, len(OracleAttestationKey)+ClaimTypeLen+len(nonce))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+ClaimTypeLen:], nonce)
	return key
}
