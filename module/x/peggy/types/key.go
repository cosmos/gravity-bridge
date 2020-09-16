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
	EthAddressKey               = []byte{0x1}
	ValsetRequestKey            = []byte{0x2}
	ValsetConfirmKey            = []byte{0x3}
	OracleClaimKey              = []byte{0x4}
	OracleAttestationKey        = []byte{0x5}
	OutgoingTXPoolKey           = []byte{0x6}
	SequenceKeyPrefix           = []byte{0x7}
	DenomiatorPrefix            = []byte{0x8}
	SecondIndexOutgoingTXFeeKey = []byte{0x9}

	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastCodeId")...)
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

func GetOutgoingTxPoolKey(id uint64) []byte {
	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
}

func GetFeeSecondIndexKey(fee sdk.Coin) []byte {
	assertPeggyVoucher(fee)

	r := make([]byte, 1+VoucherDenomLen+8)
	copy(r[0:], SecondIndexOutgoingTXFeeKey)
	copy(r[len(SecondIndexOutgoingTXFeeKey):], fee.Denom)
	copy(r[len(SecondIndexOutgoingTXFeeKey)+len(fee.Denom):], sdk.Uint64ToBigEndian(fee.Amount.Uint64()))
	return r
}

func GetDenominatorKey(voucherDenominator string) []byte {
	return append(DenomiatorPrefix, []byte(voucherDenominator)...)
}

// GetClaimKey creates a key with claim type and address first as they have a fix length.
// We can use this later for prefix scans to find all claims by type,
// claims by validator or claims for a nonce
func GetClaimKey(claimType ClaimType, nonce Nonce, validator sdk.ValAddress) []byte {
	claimTypeLen := len(claimType)
	key := make([]byte, len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonce))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen:], validator)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen:], nonce)
	return key
}

//func SplitClaimKey(raw []byte) (ClaimType, sdk.ValAddress, Nonce) {
//	return ClaimType(raw[1 : 1+ClaimTypeLen][0]), raw[1+ClaimTypeLen : 1+ClaimTypeLen+sdk.AddrLen], raw[1+ClaimTypeLen+sdk.AddrLen:]
//}

func GetAttestationKey(claimType ClaimType, nonce Nonce) []byte {
	claimTypeLen := len(claimType)
	key := make([]byte, len(OracleAttestationKey)+claimTypeLen+len(nonce))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen:], nonce)
	return key
}
