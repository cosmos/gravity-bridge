package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs a stateless check of the attestation fields
func (a Attestation) Validate() error {
	if len(a.EventID) == 0 {
		return fmt.Errorf("event id cannot be empty")
	}
	if a.Height == 0 {
		return fmt.Errorf("attestation ethereum height cannot be zero")
	}
	if len(a.Votes) > 0 && a.AttestedPower == 0 {
		return fmt.Errorf("cannot have attested power equal to zero when there are existing validator votes")
	}
	for _, validatorAddr := range a.Votes {
		_, err := sdk.ValAddressFromBech32(validatorAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
