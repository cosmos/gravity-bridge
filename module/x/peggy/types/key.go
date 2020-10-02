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
	EthAddressKey                    = []byte{0x1}
	ValsetRequestKey                 = []byte{0x2}
	ValsetConfirmKey                 = []byte{0x3}
	OracleClaimKey                   = []byte{0x4}
	OracleAttestationKey             = []byte{0x5}
	OutgoingTXPoolKey                = []byte{0x6}
	SequenceKeyPrefix                = []byte{0x7}
	DenomiatorPrefix                 = []byte{0x8}
	SecondIndexOutgoingTXFeeKey      = []byte{0x9}
	OutgoingTXBatchKey               = []byte{0xa}
	OutgoingTXBatchConfirmKey        = []byte{0xb}
	SecondIndexLastValsetApprovedKey = []byte{0xc}
	SecondIndexLastValsetObservedKey = []byte{0xd}

	// sequence keys
	KeyLastTXPoolID        = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)
)

func GetEthAddressKey(validator sdk.ValAddress) []byte {
	return append(EthAddressKey, []byte(validator)...)
}

func GetValsetRequestKey(nonce UInt64Nonce) []byte {
	return append(ValsetRequestKey, nonce.Bytes()...)
}

func GetValsetConfirmKey(nonce UInt64Nonce, validator sdk.AccAddress) []byte {
	return append(ValsetConfirmKey, append(nonce.Bytes(), []byte(validator)...)...)
}

func GetSecondIndexLastValsetApprovedKey(nonce UInt64Nonce) []byte {
	return append(SecondIndexLastValsetApprovedKey, nonce.Bytes()...)
}

func GetSecondIndexLastValsetObservedKey(nonce UInt64Nonce) []byte {
	return append(SecondIndexLastValsetObservedKey, nonce.Bytes()...)
}

func GetClaimKey(claimType ClaimType, nonce UInt64Nonce, validator sdk.ValAddress, details AttestationDetails) []byte {
	var detailsHash []byte
	if details != nil {
		detailsHash = details.Hash()
	}
	claimTypeLen := len(claimType)
	nonceBz := nonce.Bytes()
	key := make([]byte, len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz)+len(detailsHash))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen:], validator)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen:], nonceBz)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz):], detailsHash)
	return key
}

func GetAttestationKey(claimType ClaimType, nonce UInt64Nonce) []byte {
	claimTypeLen := len(claimType)
	nonceBz := nonce.Bytes()
	key := make([]byte, len(OracleAttestationKey)+claimTypeLen+len(nonceBz))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen:], nonceBz)
	return key
}

func GetOutgoingTxPoolKey(id uint64) []byte {
	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
}

func GetOutgoingTxBatchKey(batchID uint64) []byte {
	return append(OutgoingTXBatchKey, sdk.Uint64ToBigEndian(batchID)...)
}

func GetOutgoingTXBatchConfirmKey(batchID uint64, validator sdk.ValAddress) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, batchID)
	return append(OutgoingTXBatchConfirmKey, append(bz, []byte(validator)...)...)
}

func GetFeeSecondIndexKey(fee sdk.Coin) []byte {
	assertPeggyVoucher(fee)

	r := make([]byte, 1+VoucherDenomLen+8)
	copy(r[0:], SecondIndexOutgoingTXFeeKey)
	voucherDenom, _ := AsVoucherDenom(fee.Denom)
	copy(r[len(SecondIndexOutgoingTXFeeKey):], voucherDenom.Unprefixed())
	copy(r[len(SecondIndexOutgoingTXFeeKey)+len(voucherDenom.Unprefixed()):], sdk.Uint64ToBigEndian(fee.Amount.Uint64()))
	return r
}

func GetDenominatorKey(voucherDenominator string) []byte {
	return append(DenomiatorPrefix, []byte(voucherDenominator)...)
}

func DecodeUin64(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}
