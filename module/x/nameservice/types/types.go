package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Valset struct {
	Nonce       int64
	Powers      []int64
	EthAdresses []string
}

func (v Valset) GetCheckpoint() []byte {
	// Getting the equivalent of solidity's abi.encodePacked (or abi.encode) does not seem to be straightforward
	// and I am skipping it for now to focus on the overall module structure
	// https://stackoverflow.com/questions/50772811/how-can-i-get-the-same-return-value-as-solidity-abi-encodepacked-in-golang
	return []byte("dothislater")
}

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}
