package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default auth module parameter subspace
const DefaultParamspace = ModuleName

// todo: implement oracle constants as params
const AttestationPeriod = 24 * time.Hour // TODO: value????

// voting: threshold >2/3 of validator power AND > 1/2 of validator count?
var (
	// AttestationVotesCountThreshold threshold of vote counts to succeed
	AttestationVotesCountThreshold = sdk.NewUint(50)
	// AttestationVotesCountThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)
)

// Parameter keys
var (
	ParamsStoreKeyPeggyID               = []byte("PeggyID")
	ParamsStoreKeyContractHash          = []byte("ContractHash")
	ParamsStoreKeyStartThreshold        = []byte("StartThreshold")
	ParamsStoreKeyBridgeContractAddress = []byte("BridgeContractAddress")
	ParamsStoreKeyBridgeContractChainID = []byte("BridgeChainID")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyPeggyID, &p.PeggyId, validatePeggyID),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyStartThreshold, &p.StartThreshold, validateStartThreshold),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.EthereumAddress, validateBridgeContractAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func validatePeggyID(i interface{}) error {
	if _, ok := i.([]byte); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateContractHash(i interface{}) error {
	if _, ok := i.(string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateStartThreshold(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validatePeggyID(p.PeggyId); err != nil {
		return sdkerrors.Wrap(err, "peggy id")
	}
	if err := validateContractHash(p.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateStartThreshold(p.StartThreshold); err != nil {
		return sdkerrors.Wrap(err, "start threshold")
	}
	if err := validateBridgeContractAddress(p.EthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	return nil
}

func validateBridgeChainID(i interface{}) error {

	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBridgeContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := ValidateEthAddress(v); err != nil {
		// Empty addresses are valid
		if !strings.Contains(err.Error(), "empty") {
			return err
		}
	}
	return nil
}
