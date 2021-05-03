package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

type Confirm interface {
	proto.Message

	GetType() string
	GetEthSigner() string
	GetSignature() hexutil.Bytes
	Validate() error
}

var (
	_ Confirm = &ConfirmBatch{}
	_ Confirm = &ConfirmLogicCall{}
	_ Confirm = &ConfirmSignerSet{}
)

// available confirm types
const (
	ConfirmTypeBatch     = "batch"
	ConfirmTypeLogicCall = "logic_call"
	ConfirmTypeSignerSet = "signer_set"
)

// GetType should return the action
func (c ConfirmBatch) GetType() string { return ConfirmTypeBatch }

// Validate performs stateless checks
func (c ConfirmBatch) Validate() error {
	if len(c.BatchID) == 0 {
		return fmt.Errorf("batch id cannot be empty")
	}
	if err := ValidateEthAddress(c.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "ethereum signer address")
	}
	if len(c.Signature) == 0 {
		return fmt.Errorf("ethereum signature cannot be empty")
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (c ConfirmBatch) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationID is a noop to implement confirm interface
func (c ConfirmBatch) GetInvalidationID() tmbytes.HexBytes { return nil }

// GetType should return the action
func (c ConfirmLogicCall) GetType() string { return ConfirmTypeLogicCall }

// Validate performs stateless checks
func (c ConfirmLogicCall) Validate() error {
	if len(c.InvalidationID) == 0 {
		return fmt.Errorf("invalidation id cannot be empty")
	}
	if c.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce cannot be 0")
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "ethereum signer address")
	}
	if len(c.Signature) == 0 {
		return fmt.Errorf("ethereum signature cannot be empty")
	}
	return nil
}

// NewConfirmSignerSet returns a new ConfirmSignerSet
func NewConfirmSignerSet(nonce uint64, ethSigner string, signature hexutil.Bytes) *ConfirmSignerSet {
	return &ConfirmSignerSet{
		Nonce:     nonce,
		EthSigner: ethSigner,
		Signature: signature,
	}
}

// GetType should return the action
func (c ConfirmSignerSet) GetType() string { return ConfirmTypeSignerSet }

// Validate performs stateless checks
func (c ConfirmSignerSet) Validate() (err error) {
	if c.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := ValidateEthAddress(c.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "ethereum signer address")
	}
	if len(c.Signature) == 0 {
		return fmt.Errorf("ethereum signature cannot be empty")
	}
	return nil
}
