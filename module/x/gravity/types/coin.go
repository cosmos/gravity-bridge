package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func IsEthereumERC20Token(denom string) bool {
	prefix := GravityDenomPrefix + GravityDenomSeparator
	return strings.HasPrefix(denom, prefix)
}

func IsCosmosCoin(denom string) bool {
	return !IsEthereumERC20Token(denom)
}

func GravityDenomToERC20Contract(denom string) string {
	fullPrefix := GravityDenomPrefix + GravityDenomSeparator
	return strings.TrimPrefix(denom, fullPrefix)
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
		return sdkerrors.Wrapf(ErrInvalidGravityDenom, "denomination should be prefixed with the format '%s%s{address}'", GravityDenomPrefix, GravityDenomSeparator)

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