package types

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var isValidETHAddress = regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString
var emptyAddr [gethCommon.AddressLength]byte

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

func (e EthereumAddress) ValidateBasic() error {
	if !isValidETHAddress(e.String()) {
		return ErrInvalid
	}
	return nil
}

func (e EthereumAddress) IsEmpty() bool {
	return emptyAddr == e
}

// MarshalJSON marshals the etherum address to JSON
func (e EthereumAddress) MarshalJSON() ([]byte, error) {
	if e.IsEmpty() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("%q", e.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (e *EthereumAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, e[:])
}

func (e EthereumAddress) LessThan(o EthereumAddress) bool {
	return bytes.Compare(e[:], o[:]) == -1
}

// ERC20Token unique identifier for an Ethereum erc20 token.
type ERC20Token struct {
	Amount uint64 `json:"amount" yaml:"amount"`
	// Symbol is the erc20 human readable token name
	Symbol               string          `json:"symbol" yaml:"symbol"`
	TokenContractAddress EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
}

func NewERC20Token(amount uint64, symbol string, tokenContractAddress EthereumAddress) ERC20Token {
	return ERC20Token{Amount: amount, Symbol: symbol, TokenContractAddress: tokenContractAddress}
}

// String converts Token representation into a human readable form containing all data.
func (e ERC20Token) String() string {
	return fmt.Sprintf("%d %s (%s)", e.Amount, e.Symbol, e.TokenContractAddress.String())
}

// AsVoucherCoin converts the data into a cosmos coin with peggy voucher denom.
func (e ERC20Token) AsVoucherCoin() sdk.Coin {
	return sdk.NewInt64Coin(NewVoucherDenom(e.TokenContractAddress, e.Symbol).String(), int64(e.Amount))
}

func (t ERC20Token) Add(o ERC20Token) ERC20Token {
	if t.Symbol != o.Symbol {
		panic("invalid symbol")
	}
	if t.TokenContractAddress != o.TokenContractAddress {
		panic("invalid contract address")
	}
	sum := sdk.NewInt(int64(t.Amount)).AddRaw(int64(o.Amount))
	if !sum.IsUint64() {
		panic("invalid amount")
	}
	return NewERC20Token(sum.Uint64(), t.Symbol, t.TokenContractAddress)
}
