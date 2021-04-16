package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// All Ethereum claims that we relay from the Peggy contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// Which type of claim this is
	Type() ClaimType
	ValidateBasic() error
	ClaimHash() []byte

	GetOrchestratorAddress() string
}

var (
	_ EthereumClaim = &DepositClaim{}
	_ EthereumClaim = &WithdrawClaim{}
	_ EthereumClaim = &ERC20DeployedClaim{}
	_ EthereumClaim = &LogicCallExecutedClaim{}
)

// GetType returns the type of the claim
func (e *DepositClaim) Type() ClaimType {
	return ClaimType_CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (e *DepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	if err := ValidateEthAddress(e.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.OrchestratorAddress)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg DepositClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("DepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	return val
}

const (
	TypeMsgWithdrawClaim = "withdraw_claim"
)

// Hash implements BridgeDeposit.Hash
func (b *DepositClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", b.TokenContract, string(b.EthereumSender), b.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the claim type
func (e *WithdrawClaim) Type() ClaimType {
	return ClaimType_CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *WithdrawClaim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.OrchestratorAddress)
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (b *WithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.TokenContract, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

func (msg WithdrawClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("WithdrawClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	return val
}

const (
	TypeMsgDepositClaim = "deposit_claim"
)

// EthereumClaim implementation for ERC20DeployedClaim
// ======================================================

// GetType returns the type of the claim
func (e *ERC20DeployedClaim) Type() ClaimType {
	return ClaimType_CLAIM_TYPE_ERC20_DEPLOYED
}

// ValidateBasic performs stateless checks
func (e *ERC20DeployedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.OrchestratorAddress)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg ERC20DeployedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("ERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	return val
}

// Hash implements BridgeDeposit.Hash
func (b *ERC20DeployedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", b.CosmosDenom, b.TokenContract, b.Name, b.Symbol, b.Decimals)
	return tmhash.Sum([]byte(path))
}

// EthereumClaim implementation for LogicCallExecutedClaim
// ======================================================

// GetType returns the type of the claim
func (e *LogicCallExecutedClaim) Type() ClaimType {
	return ClaimType_CLAIM_TYPE_LOGIC_CALL_EXECUTED
}

// ValidateBasic performs stateless checks
func (e *LogicCallExecutedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.OrchestratorAddress)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg LogicCallExecutedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	return val
}

// Hash implements BridgeDeposit.Hash
func (b *LogicCallExecutedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.InvalidationId, b.InvalidationNonce)
	return tmhash.Sum([]byte(path))
}
