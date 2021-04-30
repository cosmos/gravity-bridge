package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs a stateless check of the attestation fields
func (a *Attestation) Validate() error {
	if len(a.EventID) == 0 {
		return fmt.Errorf("event id cannot be empty")
	}
	if a.Height == 0 {
		return fmt.Errorf("attestation ethereum height cannot be zero")
	}
	if a.Event == nil {
		return fmt.Errorf("attestation event cannot be nil")
	}
	if _, err := UnpackEvent(a.Event); err != nil {
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

func (a *Attestation) GetEthereumEvent() EthereumEvent {
	event, err := UnpackEvent(a.Event)
	if err != nil {
		panic("try validating attestation before getting event")
	}
	return event
}
