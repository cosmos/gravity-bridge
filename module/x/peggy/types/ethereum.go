package types

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	PeggyDenomPrefix      = "peggy"
	PeggyDenomSeperator   = "/"
	PeggyDenomPrefixLen   = len(PeggyDenomPrefix)
	PeggyDenomSepLen      = len(PeggyDenomSeperator)
	ETHContractAddressLen = 42
	PeggyDenomLen         = PeggyDenomPrefixLen + PeggyDenomSepLen + ETHContractAddressLen
)

// EthereumAddressLength is the length of an ETH address: 20 bytes
const EthereumAddressLength = gethCommon.AddressLength

var isValidETHAddress = regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString
var emptyAddr [EthereumAddressLength]byte

// EthereumAddress defines a standard ethereum address
type EthereumAddress gethCommon.Address

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	e := EthereumAddress(gethCommon.HexToAddress(address))
	return e //, e.ValidateBasic() // TODO: check and return error
}

func (e EthereumAddress) String() string {
	return gethCommon.Address(e).String()
}

// Bytes return the encoded address string as bytes
func (e EthereumAddress) Bytes() []byte {
	return []byte(e.String())
}

// RawBytes return the unencoded address bytes
func (e EthereumAddress) RawBytes() []byte {
	return e[:]
}

// ValidateBasic validates the address
func (e EthereumAddress) ValidateBasic() error {
	if !isValidETHAddress(e.String()) {
		return ErrInvalid
	}
	return nil
}

// IsEmpty returns if the address is empty
func (e EthereumAddress) IsEmpty() bool {
	return emptyAddr == e
}

// MarshalJSON marshals the ethereum address to JSON
func (e EthereumAddress) MarshalJSON() ([]byte, error) {
	if e.IsEmpty() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("%q", e.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (e *EthereumAddress) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, e[:])
}

// LessThan returns if an address is less than another
func (e EthereumAddress) LessThan(o EthereumAddress) bool {
	return bytes.Compare(e[:], o[:]) == -1
}

/////////////////////////
//     ERC20Token      //
/////////////////////////

// NewERC20Token returns a new instance of an ERC20
func NewERC20Token(amount uint64, contract string) *ERC20Token {
	return &ERC20Token{Amount: sdk.NewIntFromUint64(amount), Contract: contract}
}

// PeggyCoin returns the peggy representation of the ERC20
func (e *ERC20Token) PeggyCoin() sdk.Coin {
	return sdk.NewCoin(fmt.Sprintf("%s/%s", PeggyDenomPrefix, e.Contract), e.Amount)
}

// ValidateBasic permforms stateless validation
func (e *ERC20Token) ValidateBasic() error {
	if err := NewEthereumAddress(string(e.Contract)).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO: Validate all the things
	return nil
}

// String converts Token representation into a human readable form containing all data.
// func (e ERC20Token) String() string {
// 	return fmt.Sprintf("%d %s (%s)", e.Amount, e.Symbol, e.Contract.String())
// }

// AsVoucherCoin converts the data into a cosmos coin with peggy voucher denom.
func (e *ERC20Token) AsVoucherCoin() sdk.Coin {
	return e.PeggyCoin()
}

// Add adds one ERC20 to another
func (e *ERC20Token) Add(o *ERC20Token) *ERC20Token {
	if string(e.Contract) != string(o.Contract) {
		panic("invalid contract address")
	}
	sum := e.Amount.Add(o.Amount)
	if !sum.IsUint64() {
		panic("invalid amount")
	}
	return NewERC20Token(sum.Uint64(), e.Contract)
}

// ERC20FromPeggyCoin returns the ERC20 representation of a given peggy coin
func ERC20FromPeggyCoin(v sdk.Coin) (*ERC20Token, error) {
	if !IsPeggyCoin(v) {
		return nil, fmt.Errorf("%s isn't a valid peggy coin", v.String())
	}
	return &ERC20Token{Contract: strings.Split(v.Denom, "/")[1], Amount: v.Amount}, nil
}

func assertPeggyVoucher(s sdk.Coin) {
	if !IsPeggyCoin(s) {
		panic("invalid denom type")
	}
	if !s.IsValid() {
		panic("invalid amount type")
	}
}

// IsPeggyCoin returns true if a coin is a peggy representation of an ERC20 token
func IsPeggyCoin(v sdk.Coin) bool {
	spl := strings.Split(v.Denom, "/")
	err := NewEthereumAddress(spl[1]).ValidateBasic()
	switch {
	case len(spl) != 2:
		return false
	case spl[0] != "peggy":
		return false
	case err != nil:
		return false
	case len(v.Denom) != PeggyDenomLen:
		return false
	default:
		return true
	}
}
