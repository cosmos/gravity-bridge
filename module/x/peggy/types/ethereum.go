package types

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethCommon "github.com/ethereum/go-ethereum/common"
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
// const EthereumAddressLength = gethCommon.AddressLength

// EthereumAddress defines a standard ethereum address
type EthereumAddress gethCommon.Address

// EthAddrLessThan migrates the Ethereum address less than function
func EthAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	return EthereumAddress(gethCommon.HexToAddress(address)) //, e.ValidateBasic() // TODO: check and return error
}

func (e EthereumAddress) String() string {
	return gethCommon.Address(e).String()
}

// ValidateEthAddress validates the ethereum address strings
func ValidateEthAddress(a string) error {
	if a == "" {
		return fmt.Errorf("empty")
	}
	if !regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(a) {
		return fmt.Errorf("address(%s) doesn't pass regex", a)
	}
	return nil
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
	if err := ValidateEthAddress(e.Contract); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO: Validate all the things
	return nil
}

// AsVoucherCoin converts the data into a cosmos coin with peggy voucher denom.
func (e *ERC20Token) AsVoucherCoin() sdk.Coin {
	return e.PeggyCoin()
}

// Add adds one ERC20 to another
// TODO: make this return errors instead
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

// IsPeggyCoin returns true if a coin is a peggy representation of an ERC20 token
func IsPeggyCoin(v sdk.Coin) bool {
	spl := strings.Split(v.Denom, "/")
	err := ValidateEthAddress(spl[1])
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
