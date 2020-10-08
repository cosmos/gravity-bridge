package types

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nonce defines an abstract nonce type that is unique within it's context.
type nonce interface {
	// String returns an encoded human readable representation. Used in URLs.
	String() string
	// Bytes returns an encoded raw bytes representation
	Bytes() []byte
	// ValidateBasic returns the result of the syntax check
	ValidateBasic() error
	// GreaterThan than other.
	GreaterThan(o nonce) bool
	IsEmpty() bool
}

const UInt64NonceByteLen = 8

type UInt64Nonce uint64

var _ nonce = NewUInt64Nonce(0)

func NewUInt64Nonce(s uint64) UInt64Nonce {
	return UInt64Nonce(s)
}

// UInt64NonceFromBytes create UInt64Nonce from binary big endian representation
func UInt64NonceFromBytes(s []byte) UInt64Nonce {
	return NewUInt64Nonce(binary.BigEndian.Uint64(s))
}

// UInt64NonceFromBytes create UInt64Nonce from human readable string representation
func UInt64NonceFromString(s string) (UInt64Nonce, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return NewUInt64Nonce(0), err
	}
	return NewUInt64Nonce(v), nil
}

func Uint64FromNonce(n nonce) (uint64, error) {
	if v, ok := n.(nonceUint64er); ok {
		return v.Uint64(), nil
	}
	return 0, ErrUnsupported
}

func (n UInt64Nonce) Uint64() uint64 {
	return uint64(n)
}

func (n UInt64Nonce) String() string {
	return strconv.FormatUint(n.Uint64(), 10)
}

func (n UInt64Nonce) Bytes() []byte {
	return sdk.Uint64ToBigEndian(n.Uint64())
}

func (n UInt64Nonce) ValidateBasic() error {
	if n.IsEmpty() {
		return ErrEmpty
	}
	return nil
}

type nonceUint64er interface {
	Uint64() uint64
}

func (n UInt64Nonce) GreaterThan(o nonce) bool {
	if o == nil || reflect.ValueOf(o).IsZero() || o.IsEmpty() {
		return true
	}
	if n.IsEmpty() {
		return false
	}
	if v, ok := o.(nonceUint64er); ok {
		return n.Uint64() > v.Uint64()
	}
	return bytes.Compare(n.Bytes(), o.Bytes()) == 1
}

func (n UInt64Nonce) IsEmpty() bool {
	return n == 0
}
