package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAddressConfig sets the gravity app's address configuration.
func SetAddressConfig() {
	config := sdk.GetConfig()

	config.SetAddressVerifier(VerifyAddressFormat)
	config.Seal()
}
