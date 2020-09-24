package types

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthereumAddress defines a standard ethereum address
type EthereumAddress gethCommon.Address

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	return EthereumAddress(gethCommon.HexToAddress(address))
}

// Route should return the name of the module
func (ethAddr EthereumAddress) String() string {
	return gethCommon.Address(ethAddr).String()
}

// MarshalJSON marshals the etherum address to JSON
func (ethAddr EthereumAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ethAddr.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (ethAddr *EthereumAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, ethAddr[:])
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
