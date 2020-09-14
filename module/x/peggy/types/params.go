package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// DefaultParamspace defines the default auth module parameter subspace
const DefaultParamspace = ModuleName

// todo: implement oracle constants as params

const AttestationPeriod = 24 * time.Hour // TODO: value????
// voting: threshold >2/3 of validator power AND > 1/2 of validator count?
var (
	AttestationVotesCountThreshold = Fraction{1, 2}
	AttestationVotesPowerThreshold = Fraction{2, 3}
	BridgeContractAddress          = NewEthereumAddress("")
	BridgeContractChainID          = "0" // todo: revisit
)

// TODO: Defaults don't make sense for any of our params. Should we still have them?

// Parameter keys
var (
	KeyPeggyID      = []byte("PeggyID")
	KeyContractHash = []byte("ContractHash")
	KeyStartBlock   = []byte("StartBlock")
)

var _ subspace.ParamSet = &Params{}

type Params struct {
	PeggyID      []byte `json:"peggy_id" yaml:"peggy_id"`
	ContractHash []byte `json:"contract_source_hash" yaml:"contract_source_hash"`
	StartBlock   uint64 `json:"start_block" yaml:"start_block"`
}

// NewParams creates a new Params object
func NewParams(peggyID []byte, contractHash []byte, startBlock uint64) Params {
	return Params{
		PeggyID:      peggyID,
		ContractHash: contractHash,
		StartBlock:   startBlock,
	}
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
		params.NewParamSetPair(KeyPeggyID, &p.PeggyID, validatePeggyID),
		params.NewParamSetPair(KeyContractHash, &p.ContractHash, validateContractHash),
		params.NewParamSetPair(KeyStartBlock, &p.StartBlock, validateStartBlock),
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
	sb.WriteString(fmt.Sprintf("StartBlock: %d\n", p.StartBlock))
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

func validateStartBlock(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validatePeggyID(p.PeggyID); err != nil {
		return err
	}
	if err := validateContractHash(p.ContractHash); err != nil {
		return err
	}
	if err := validateStartBlock(p.StartBlock); err != nil {
		return err
	}

	return nil
}
