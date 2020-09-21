package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// OutgoingTx is a withdrawal on the bridged contract
type OutgoingTx struct {
	Sender      sdk.AccAddress `json:"sender"`
	DestAddress string         `json:"dest_address"`
	Amount      sdk.Coin       `json:"send"`
	BridgeFee   sdk.Coin       `json:"bridge_fee"`
	//TokenContractAddress string         `json:"bridge_contract_address"` // todo: do we need this?
}

// BridgedDenominator contains bridged token details
type BridgedDenominator struct {
	//ChainID         string
	TokenContractAddress EthereumAddress
	Symbol               string
}

const (
	VoucherDenomPrefix = "peggy"
	DenomSeparator     = "" // todo: only a-z0-9 supported
	//DenomSeparator     = "/"  // todo: not allowed in this versions sdk coin demon
	VoucherDenomLen  = 15 // todo: cut to 15 to match this versions sdk coin demon
	voucherPrefixLen = len(VoucherDenomPrefix + DenomSeparator)

	//hashLen            = 64
	//separatorLen       = len(DenomSeparator)
	//prefixLen          = len(VoucherDenomPrefix)
	//VoucherDenomLen    = hashLen + prefixLen + separatorLen
)

func assertPeggyVoucher(s sdk.Coin) {
	if !IsVoucherDenom(s.Denom) {
		panic("invalid denom type")
	}
	if s.Amount.IsNegative() || !s.Amount.IsUint64() {
		panic("invalid amount type")
	}
}

// VoucherDenom is a unique denominator and identifier for a bridged token.
type VoucherDenom string

func NewVoucherDenom(tokenContractAddr EthereumAddress, symbol string) VoucherDenom {
	denomTrace := fmt.Sprintf("%s/%s/", tokenContractAddr.String(), symbol)
	var hash tmbytes.HexBytes = tmhash.Sum([]byte(denomTrace))
	simpleVoucherDenom := VoucherDenomPrefix + DenomSeparator + hash.String()
	sdkVersionHackDenom := strings.ToLower(simpleVoucherDenom[0:15]) // todo: up to 15 chars (lowercase) allowed in this sdk version only
	return VoucherDenom(sdkVersionHackDenom)
}

// AsVoucherDenom type conversion with `IsVoucherDenom` check.
func AsVoucherDenom(raw string) (VoucherDenom, error) {
	if !IsVoucherDenom(raw) {
		return "", sdkerrors.Wrap(ErrInvalid, "not a voucher denom")
	}
	return VoucherDenom(raw), nil
}
func (d VoucherDenom) Unprefixed() string {
	return string(d[voucherPrefixLen:])
}

func (d VoucherDenom) String() string {
	return string(d)
}

// IsVoucherDenom verifies the given string matches the peggy voucher conditions
func IsVoucherDenom(denom string) bool {
	return len(denom) == VoucherDenomLen && strings.HasPrefix(denom, VoucherDenomPrefix)
}

// IDSet is a collection of DB keys
type IDSet []uint64
