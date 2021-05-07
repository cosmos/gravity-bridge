package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GravityDenom struct {
	sdk.Coin
}

func (gd *GravityDenom) IsEthereumERC20Token() bool {
	prefix := GravityDenomPrefix + GravityDenomSeparator
	return strings.HasPrefix(gd.Denom, prefix)
}

func (gd *GravityDenom) IsCosmosCoin() bool {
	return !gd.IsEthereumERC20Token()
}

func (gd *GravityDenom) GravityDenomToERC20Contract() common.Address {
	fullPrefix := GravityDenomPrefix + GravityDenomSeparator
	return common.HexToAddress(strings.TrimPrefix(gd.Denom, fullPrefix))
}

func GravityDenomFromContract(contract string) GravityDenom {
	return GravityDenom{sdk.Coin{Denom: contract, Amount: sdk.NewIntFromUint64(0)}}
}

func (gd *GravityDenom) ERC20Token() ERC20Token {
	return ERC20Token{sdk.Coin{Amount: gd.Amount, Denom: gd.Denom}}
}


// ValidateGravityDenom validates that the given denomination is either:
//
//  - A valid base denomination (eg: 'uatom')
//  - A valid gravity bridge token representation (i.e 'gravity/{address}')
func ValidateGravityDenom(denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return err
	}

	denomSplit := strings.SplitN(denom, GravityDenomSeparator, 2)

	switch {
	case strings.TrimSpace(denom) == "",
		len(denomSplit) == 1 && denomSplit[0] == GravityDenomPrefix,
		len(denomSplit) == 2 && (denomSplit[0] != GravityDenomPrefix || strings.TrimSpace(denomSplit[1]) == ""):
		return sdkerrors.Wrapf(fmt.Errorf("invalid gravity denom"), "denomination should be prefixed with the format '%s%s{address}'", GravityDenomPrefix, GravityDenomSeparator)

	case denomSplit[0] == denom && strings.TrimSpace(denom) != "":
		// denom source is from the current chain. Return nil as it has already been validated
		return nil
	}

	// denom source is ethereum. Validate the ethereum hex address
	if err := ValidateEthAddress(denomSplit[1]); err != nil {
		return fmt.Errorf("invalid contract address: %w", err)
	}

	return nil
}