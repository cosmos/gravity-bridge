package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
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

type Params struct {
	// PeggyID is a random 32 byte value to prevent signature reuse
	PeggyID []byte `json:"peggy_id,omitempty" yaml:"peggy_id"`
	// ContractHash is the code hash of a known good version of the Peggy contract solidity code.
	// It will be used to verify exactly which version of the bridge will be deployed.
	ContractHash []byte `json:"contract_source_hash,omitempty" yaml:"contract_source_hash"`
	// StartThreshold is the percentage of total voting power that must be online and participating in
	// Peggy operations before a bridge can start operating
	StartThreshold uint64 `json:"start_threshold,omitempty" yaml:"start_threshold"`
	// BridgeContractAddress is address of the bridge contract on the Ethereum side
	BridgeContractAddress EthereumAddress `json:"bridge_contract_address,omitempty" yaml:"bridge_contract_address"`
	// BridgeChainID is the unique identifier of the Ethereum chain
	BridgeChainID uint64 `json:"bridge_chain_id,omitempty" yaml:"bridge_chain_id"`
}

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(ParamsStoreKeyPeggyID, &p.PeggyID, validatePeggyID),
		params.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractHash, validateContractHash),
		params.NewParamSetPair(ParamsStoreKeyStartThreshold, &p.StartThreshold, validateStartThreshold),
		params.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeContractAddress, validateBridgeContractAddress),
		params.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainID, validateBridgeChainID),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("PeggyID: %d\n", p.PeggyID))
	sb.WriteString(fmt.Sprintf("ContractHash: %d\n", p.ContractHash))
	sb.WriteString(fmt.Sprintf("StartThreshold: %d\n", p.StartThreshold))
	sb.WriteString(fmt.Sprintf("BridgeContractAddress: %s\n", p.BridgeContractAddress.String()))
	sb.WriteString(fmt.Sprintf("BridgeChainID: %d\n", p.BridgeChainID))
	return sb.String()
}

func validatePeggyID(i interface{}) error {
	_, ok := i.([]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateContractHash(i interface{}) error {
	_, ok := i.([]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateStartThreshold(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validatePeggyID(p.PeggyID); err != nil {
		return sdkerrors.Wrap(err, "peggy id")
	}
	if err := validateContractHash(p.ContractHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateStartThreshold(p.StartThreshold); err != nil {
		return sdkerrors.Wrap(err, "start threshold")
	}
	if err := validateBridgeContractAddress(p.BridgeContractAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainID); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	return nil
}

func validateBridgeChainID(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBridgeContractAddress(i interface{}) error {
	v, ok := i.(EthereumAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.ValidateBasic()
}
