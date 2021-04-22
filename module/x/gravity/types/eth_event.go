package types

import (
	"crypto/sha256"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// EthereumEvent represents a event on ethereum state
type EthereumEvent interface {
	proto.Message
	// All Ethereum event that we relay from the bridge contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what event goes to what nonce means someone is lying.
	GetNonce() uint64
	// The block height that the evented event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all events. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetEthereumHeight() uint64
	GetType() string
	Validate() error
	Hash() tmbytes.HexBytes
}

var (
	_ EthereumEvent = &DepositEvent{}
	_ EthereumEvent = &WithdrawEvent{}
	_ EthereumEvent = &CosmosERC20DeployedEvent{}
	_ EthereumEvent = &LogicCallExecutedEvent{}
)

// GetType returns the type of the event
func (e DepositEvent) GetType() string {
	return "deposit"
}

// Validate performs stateless checks
func (e DepositEvent) Validate() error {
	if e.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if !e.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	if err := ValidateEthAddress(e.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "ethereum sender")
	}
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	if e.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

const ()

// Hash implements BridgeDeposit.Hash
func (e DepositEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%s/%s/", e.TokenContract, e.EthereumSender, e.CosmosReceiver)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

// GetType returns the event type
func (e WithdrawEvent) GetType() string {
	return "withdraw"
}

// Validate performs stateless checks
func (e WithdrawEvent) Validate() error {
	if e.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if e.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (e WithdrawEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%d/", e.TokenContract, e.Nonce)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

// EthereumEvent implementation for CosmosERC20DeployedEvent
// ======================================================

// GetType returns the type of the event
func (e CosmosERC20DeployedEvent) GetType() string {
	return "cosmos_erc20_deployed"
}

// Validate performs stateless checks
func (e CosmosERC20DeployedEvent) Validate() error {
	if e.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if err := sdk.ValidateDenom(e.CosmosDenom); err != nil {
		return err
	}
	if strings.TrimSpace(e.Name) == "" {
		return fmt.Errorf("token name cannot be blank")
	}
	if strings.TrimSpace(e.Symbol) == "" {
		return fmt.Errorf("token symbol cannot be blank")
	}
	if e.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

// Hash implements BridgeDeposit.Hash
func (e CosmosERC20DeployedEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", e.CosmosDenom, e.TokenContract, e.Name, e.Symbol, e.Decimals)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

// EthereumEvent implementation for LogicCallExecutedEvent
// ======================================================

// GetType returns the type of the event
func (e LogicCallExecutedEvent) GetType() string {
	return "logic_call_executed"
}

// Validate performs stateless checks
func (e LogicCallExecutedEvent) Validate() error {
	if e.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if len(e.InvalidationID) == 0 {
		return fmt.Errorf("invalidation id cannot be empty")
	}
	if e.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce cannot be 0")
	}
	if e.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

// Hash implements BridgeDeposit.Hash
func (e LogicCallExecutedEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%d/", e.InvalidationID, e.InvalidationNonce)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}
