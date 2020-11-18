package types

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

// // ERC20Token unique identifier for an Ethereum erc20 token.
// type ERC20Token struct {
// 	Amount uint64 `json:"amount" yaml:"amount"`
// 	// Symbol is the erc20 human readable token name
// 	Symbol               string          `json:"symbol" yaml:"symbol"`
// 	TokenContractAddress EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
// }

// NewERC20Token returns a new instance of an ERC20
func NewERC20Token(amount uint64, symbol string, tokenContractAddress EthereumAddress) *ERC20Token {
	return &ERC20Token{Amount: sdk.NewInt(int64(amount)), Symbol: symbol, TokenContractAddress: tokenContractAddress.String()}
}

// ValidateBasic permforms stateless validation
func (e *ERC20Token) ValidateBasic() error {
	if err := NewEthereumAddress(string(e.TokenContractAddress)).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO: Validate all the things
	return nil
}

// String converts Token representation into a human readable form containing all data.
// func (e ERC20Token) String() string {
// 	return fmt.Sprintf("%d %s (%s)", e.Amount, e.Symbol, e.TokenContractAddress.String())
// }

// AsVoucherCoin converts the data into a cosmos coin with peggy voucher denom.
func (e *ERC20Token) AsVoucherCoin() sdk.Coin {
	return sdk.NewInt64Coin(NewVoucherDenom(NewEthereumAddress(string(e.TokenContractAddress)), e.Symbol).String(), e.Amount.Int64())
}

// Add adds one ERC20 to another
func (t *ERC20Token) Add(o *ERC20Token) *ERC20Token {
	if t.Symbol != o.Symbol {
		panic("invalid symbol")
	}
	if string(t.TokenContractAddress) != string(o.TokenContractAddress) {
		panic("invalid contract address")
	}
	sum := t.Amount.Add(o.Amount)
	if !sum.IsUint64() {
		panic("invalid amount")
	}
	return NewERC20Token(sum.Uint64(), t.Symbol, NewEthereumAddress(string(t.TokenContractAddress)))
}
