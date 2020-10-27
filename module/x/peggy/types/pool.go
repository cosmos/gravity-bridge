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
	Sender      sdk.AccAddress  `json:"sender"`
	DestAddress EthereumAddress `json:"dest_address"`
	Amount      sdk.Coin        `json:"send"`
	BridgeFee   sdk.Coin        `json:"bridge_fee"`
}

// BridgedDenominator track and identify the ported ERC20 tokens into Peggy.
// An ERC20 token on the Ethereum side can be uniquely identified by the ERC20 contract address and the human readable symbol
// used for it in the contract
// In Peggy this is represented as "vouchers" that get minted and burned when interacting with the Ethereum side. These "vouchers"
// are identified by a prefixed string representation. See VoucherDenom type.
type BridgedDenominator struct {
	// TokenContractAddress is the ERC20 contract address
	TokenContractAddress EthereumAddress `json:"token_contract_address"`
	// Symbol is the human readable ERC20 token name.
	Symbol string `json:"symbol"`
	// CosmosVoucherDenom is used as denom in sdk.Coin
	CosmosVoucherDenom VoucherDenom `json:"cosmos_voucher_denom"`
}

func NewBridgedDenominator(tokenContractAddress EthereumAddress, erc20Symbol string) BridgedDenominator {
	v := NewVoucherDenom(tokenContractAddress, erc20Symbol)
	return BridgedDenominator{TokenContractAddress: tokenContractAddress, Symbol: erc20Symbol, CosmosVoucherDenom: v}
}

// ToERC20Token converts the given voucher amount to the matching ERC20Token object of same type
func (b BridgedDenominator) ToERC20Token(s sdk.Coin) ERC20Token {
	if b.CosmosVoucherDenom.String() != s.Denom {
		panic("invalid denom")
	}
	return b.ToUint64ERC20Token(s.Amount.Uint64())
}

// ToUint64ERC20Token creates a erc20 token instance for given amount
func (b BridgedDenominator) ToUint64ERC20Token(amount uint64) ERC20Token {
	return NewERC20Token(amount, b.Symbol, b.TokenContractAddress)
}

// ToVoucherCoin creates a new Peggy voucher coin instance with given amount
func (b BridgedDenominator) ToVoucherCoin(amount uint64) sdk.Coin {
	return sdk.NewInt64Coin(b.CosmosVoucherDenom.String(), int64(amount))
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

// VoucherDenom is a unique denominator and identifier for an ERC20 token in the cosmos world
type VoucherDenom string

// NewVoucherDenom builds a Peggy voucher denominator from the ERC20 contract address and the human readable ERC20 symbol.
func NewVoucherDenom(tokenContractAddr EthereumAddress, erc20Symbol string) VoucherDenom {
	denomTrace := fmt.Sprintf("%s/%s/", tokenContractAddr.String(), erc20Symbol)
	var hash tmbytes.HexBytes = tmhash.Sum([]byte(denomTrace))
	simpleVoucherDenom := VoucherDenomPrefix + DenomSeparator + hash.String()
	// todo: up to 15 chars (lowercase) allowed in this sdk version only
	// THIS NEEDS TO BE CHANGED BEFORE PRODUCTION to not truncate the address.
	// The truncation weakens the collision resistance of the address.
	sdkVersionHackDenom := strings.ToLower(simpleVoucherDenom[0:15])
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

// IDSet is a collection of DB keys in a second index.
type IDSet []uint64
