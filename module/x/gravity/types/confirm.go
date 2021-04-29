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
	GetNonce() uint64
	GetSignature() hexutil.Bytes
	Validate() error
}

var (
	_ Confirm = &ConfirmBatch{}
	_ Confirm = &ConfirmLogicCall{}
	_ Confirm = &ConfirmSignerSet{}
)

// GetType should return the action
func (c ConfirmBatch) GetType() string { return "batch" }

// Validate performs stateless checks
func (c ConfirmBatch) Validate() error {
	if c.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := ValidateEthAddress(c.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
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
func (c ConfirmLogicCall) GetType() string { return "logic_Call" }

// Validate performs stateless checks
func (c ConfirmLogicCall) Validate() error {
	if len(c.Signature) == 0 {
		return fmt.Errorf("ethereum signature cannot be empty")
	}
	if len(c.InvalidationID) == 0 {
		return fmt.Errorf("invalidation id is empty")
	}
	if c.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce cannot be 0")
	}
	return nil
}

func (c ConfirmLogicCall) GetNonce() uint64 {
	return 0
}

func (c ConfirmLogicCall) GetTokenContract() string {
	return ""
}

// NewConfirmSignerSet returns a new ConfirmSignerSet
func NewConfirmSignerSet(nonce uint64, signature hexutil.Bytes) *ConfirmSignerSet {
	return &ConfirmSignerSet{
		Nonce:     nonce,
		Signature: signature,
	}
}

// GetType should return the action
func (c ConfirmSignerSet) GetType() string { return "valset" }

// Validate performs stateless checks
func (c ConfirmSignerSet) Validate() (err error) {
	if c.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if len(c.Signature) == 0 {
		return fmt.Errorf("ethereum signature cannot be empty")
	}
	return nil
}

// GetInvalidationNonce is a noop to implement confirm interface
func (c ConfirmSignerSet) GetInvalidationNonce() uint64 { return 0 }

// GetInvalidationID is a noop to implement confirm interface
func (c ConfirmSignerSet) GetInvalidationID() tmbytes.HexBytes { return nil }

func (c ConfirmSignerSet) GetTokenContract() string {
	return ""
}
