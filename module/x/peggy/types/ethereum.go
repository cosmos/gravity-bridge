package types

import (
	"fmt"
	"reflect"

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
	Amount               int64           `json:"amount" yaml:"amount"`
	Symbol               string          `json:"symbol" yaml:"symbol"`
	TokenContractAddress EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
}
