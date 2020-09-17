package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

type OutgoingTx struct {
	Sender      sdk.AccAddress `json:"sender"`
	DestAddress string         `json:"dest_address"`
	Amount      sdk.Coin       `json:"send"`
	BridgeFee   sdk.Coin       `json:"bridge_fee"`
	//BridgeContractAddress string         `json:"bridge_contract_address"` // todo: do we need this?
}

type BridgedDenominator struct {
	//ChainID         string
	BridgeContractAddress string
	TokenID               string
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

type VoucherDenom string

func NewVoucherDenom(contractAddr, token string) VoucherDenom {
	denomTrace := fmt.Sprintf("%s/%s/", contractAddr, token)
	var hash tmbytes.HexBytes = tmhash.Sum([]byte(denomTrace))
	simpleVoucherDenom := VoucherDenomPrefix + DenomSeparator + hash.String()
	sdkVersionHackDenom := strings.ToLower(simpleVoucherDenom[0:15]) // todo: up to 15 chars (lowercase) allowed in this sdk version only
	return VoucherDenom(sdkVersionHackDenom)
}

func AsVoucherDenom(s string) (VoucherDenom, error) {
	if !IsVoucherDenom(s) {
		return "", sdkerrors.Wrap(ErrInvalid, "not a voucher denom")
	}
	return VoucherDenom(s), nil
}
func (d VoucherDenom) Unprefixed() string {
	return string(d[voucherPrefixLen:])
}

func IsVoucherDenom(denom string) bool {
	return len(denom) == VoucherDenomLen && strings.HasPrefix(denom, VoucherDenomPrefix)
}

// IDSet is a collection of DB keys
type IDSet []uint64
