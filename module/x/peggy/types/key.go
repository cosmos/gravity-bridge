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
	OutgoingTXBatchKey          = []byte{0xa}
	// deprecated

	OutgoingTXBatchConfirmKey    = []byte{0xb}
	BatchConfirmKey              = []byte{0xe1}
	SecondIndexNonceByClaimKey   = []byte{0xf}
	LastEventNonceByValidatorKey = []byte{0xf1}
	BridgeObservedSignatureKey   = []byte{0x10}
	LastObservedEventNonceKey    = []byte{0xf2}

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

// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetValsetConfirmKey(nonce UInt64Nonce, validator sdk.AccAddress) []byte {
	return append(ValsetConfirmKey, append(nonce.Bytes(), []byte(validator)...)...)
}

func GetClaimKey(claimType ClaimType, nonce UInt64Nonce, validator sdk.ValAddress, details AttestationDetails) []byte {
	var detailsHash []byte
	if details != nil {
		detailsHash = details.Hash()
	}
	claimTypeLen := len(claimType.Bytes())
	nonceBz := nonce.Bytes()
	key := make([]byte, len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz)+len(detailsHash))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], claimType.Bytes())
	copy(key[len(OracleClaimKey)+claimTypeLen:], validator)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen:], nonceBz)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz):], detailsHash)
	return key
}

func GetLastNonceByClaimTypeSecondIndexKeyPrefix(claimType ClaimType) []byte {
	return append(SecondIndexNonceByClaimKey, claimType.Bytes()...)
}

func GetLastNonceByClaimTypeSecondIndexKey(claimType ClaimType, nonce UInt64Nonce) []byte {
	return append(GetLastNonceByClaimTypeSecondIndexKeyPrefix(claimType), nonce.Bytes()...)
}

func GetAttestationKey(eventNonce UInt64Nonce, details AttestationDetails) []byte {
	return append(eventNonce.Bytes(), details.Hash()...)
}

func GetOutgoingTxPoolKey(id uint64) []byte {
	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
}

func GetOutgoingTxBatchKey(tokenContract EthereumAddress, nonce UInt64Nonce) []byte {
	return append(append(OutgoingTXBatchKey, tokenContract.Bytes()...), nonce.Bytes()...)
}

func GetBatchConfirmKey(tokenContract EthereumAddress, batchNonce UInt64Nonce, validator sdk.AccAddress) []byte {
	a := append(batchNonce.Bytes(), validator...)
	b := append(tokenContract.Bytes(), a...)
	c := append(BatchConfirmKey, b...)
	return c
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

func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

func GetLastObservedEventNonceKey() []byte {
	return LastObservedEventNonceKey
}
