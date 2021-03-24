package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	peggyparams "github.com/cosmos/gravity-bridge/module/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for gravity.
func MakeEncodingConfig() peggyparams.EncodingConfig {
	encodingConfig := peggyparams.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
