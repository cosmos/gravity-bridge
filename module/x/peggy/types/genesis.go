package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default auth module parameter subspace
const (
	// todo: implement oracle constants as params
	DefaultParamspace = ModuleName
	AttestationPeriod = 24 * time.Hour // TODO: value????
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeyPeggyID stores the peggy id
	ParamsStoreKeyPeggyID = []byte("PeggyID")

	// ParamsStoreKeyContractHash stores the contract hash
	ParamsStoreKeyContractHash = []byte("ContractHash")

	// ParamsStoreKeyStartThreshold stores the start threshold
	ParamsStoreKeyStartThreshold = []byte("StartThreshold")

	// ParamsStoreKeyBridgeContractAddress stores the contract address
	ParamsStoreKeyBridgeContractAddress = []byte("BridgeContractAddress")

	// ParamsStoreKeyBridgeContractChainID stores the bridge chain id
	ParamsStoreKeyBridgeContractChainID = []byte("BridgeChainID")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// TODO: set some better defaults here
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// // DefaultParams returns a copy of the default params
// func DefaultParams() *Params {
// 	return &Params{
// 		PeggyId: "defaultpeggyid",
// 	}
// }

// // ValidateBasic checks that the parameters have valid values.
// func (p Params) ValidateBasic() error {
// 	if err := validatePeggyID(p.PeggyId); err != nil {
// 		return sdkerrors.Wrap(err, "peggy id")
// 	}
// 	if err := validateContractHash(p.ContractSourceHash); err != nil {
// 		return sdkerrors.Wrap(err, "contract hash")
// 	}
// 	if err := validateStartThreshold(p.StartThreshold); err != nil {
// 		return sdkerrors.Wrap(err, "start threshold")
// 	}
// 	if err := validateBridgeContractAddress(p.EthereumAddress); err != nil {
// 		return sdkerrors.Wrap(err, "bridge contract address")
// 	}
// 	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
// 		return sdkerrors.Wrap(err, "bridge chain id")
// 	}
// 	return nil
// }

// // ParamKeyTable for auth module
// func ParamKeyTable() paramtypes.KeyTable {
// 	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
// }

// // ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// // pairs of auth module's parameters.
// func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
// 	return paramtypes.ParamSetPairs{
// 		paramtypes.NewParamSetPair(ParamsStoreKeyPeggyID, &p.PeggyId, validatePeggyID),
// 		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractHash),
// 		paramtypes.NewParamSetPair(ParamsStoreKeyStartThreshold, &p.StartThreshold, validateStartThreshold),
// 		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.EthereumAddress, validateBridgeContractAddress),
// 		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
// 	}
// }

// // Equal returns a boolean determining if two Params types are identical.
// func (p Params) Equal(p2 Params) bool {
// 	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
// 	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
// 	return bytes.Equal(bz1, bz2)
// }

// func validatePeggyID(i interface{}) error {
// 	v, ok := i.(string)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}
// 	if _, err := strToFixByteArray(v); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func validateContractHash(i interface{}) error {
// 	if _, ok := i.(string); !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}
// 	return nil
// }

// func validateStartThreshold(i interface{}) error {
// 	if _, ok := i.(uint64); !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}
// 	return nil
// }

// func validateBridgeChainID(i interface{}) error {
// 	if _, ok := i.(uint64); !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}
// 	return nil
// }

// func validateBridgeContractAddress(i interface{}) error {
// 	v, ok := i.(string)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}
// 	if err := ValidateEthAddress(v); err != nil {
// 		// TODO: ensure that empty addresses are valid in params
// 		if !strings.Contains(err.Error(), "empty") {
// 			return err
// 		}
// 	}
// 	return nil
// }

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
