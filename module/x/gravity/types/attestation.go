package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs a stateless check of the attestation fields
func (a Attestation) Validate() error {
	if a.Height == 0 {
		return fmt.Errorf("attestation ethereum height cannot be zero")
	}
	if len(a.Votes) == 0 {
		return fmt.Errorf("cannot have a attestation with no votes attached")
	}
	eve, err := UnpackEvent(a.Event)
	if err != nil {
		return err
	}
	if err := eve.Validate(); err != nil {
		return err
	}
	for _, validatorAddr := range a.Votes {
		_, err := sdk.ValAddressFromBech32(validatorAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
