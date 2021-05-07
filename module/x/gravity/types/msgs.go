package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	_ sdk.Msg = &MsgDelegateKeys{}
	_ sdk.Msg = &MsgSendToEthereum{}
	_ sdk.Msg = &MsgRequestBatchTx{}
	_ sdk.Msg = &MsgSubmitEthereumEvent{}
	_ sdk.Msg = &MsgSubmitEthereumSignature{}
	_ sdk.Msg = &MsgCancelSendToEthereum{}

	//_ sdk.Msg = &MsgUpdateSignerSetTx{}
	//_ sdk.Msg = &MsgConfirmBatch{}
	//_ sdk.Msg = &MsgERC20DeployedClaim{}
	//_ sdk.Msg = &MsgConfirmLogicCall{}
	//_ sdk.Msg = &MsgLogicCallExecutedClaim{}
	//_ sdk.Msg = &MsgDepositClaim{}
	//_ sdk.Msg = &MsgWithdrawClaim{}
)

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSetOrchestratorAddress(val sdk.ValAddress, oper sdk.AccAddress, eth string) *MsgDelegateKeys {
	return &MsgDelegateKeys{
		ValidatorAddress:    val.String(),
		OrchestratorAddress: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKeys) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKeys) Type() string { return "delegate_keys"}

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKeys) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress)
	}
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgDelegateKeys) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgDelegateKeys) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}


// Route should return the name of the module
func (msg *MsgSubmitEthereumEvent) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEthereumEvent) Type() string { return "submit_ethereum_event" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEthereumEvent) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event := EthereumEvent(msg.Event)

	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEthereumEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEthereumEvent) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToEthereum returns a new MsgSendToEthereum
func NewMsgSendToEthereum(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEthereum {
	return &MsgSendToEthereum{
		Sender:    sender.String(),
		EthRecipient:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
	}
}

// Route should return the name of the module
func (msg MsgSendToEthereum) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToEthereum) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEthereum) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
	if msg.Amount.Denom != msg.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins,
			fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := ValidateEthAddress(msg.EthRecipient); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToEthereum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEthereum) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatchTx returns a new msgRequestBatch
func NewMsgRequestBatchTx(orchestrator sdk.AccAddress) *MsgRequestBatchTx {
	return &MsgRequestBatchTx{
		Signer: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatchTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatchTx) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatchTx) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatchTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatchTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConfirmBatch) Type() string { return "confirm_batch" }

// ValidateBasic performs stateless checks
func (msg MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgConfirmLogicCall) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConfirmLogicCall) Type() string { return "confirm_logic" }

// ValidateBasic performs stateless checks
func (msg MsgConfirmLogicCall) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	_, err = hex.DecodeString(msg.InvalidationId)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.InvalidationId)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmLogicCall) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmLogicCall) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// All Ethereum claims that we relay from the Gravity contract and into the module
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
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ EthereumClaim = &MsgDepositClaim{}
	_ EthereumClaim = &MsgWithdrawClaim{}
	_ EthereumClaim = &MsgERC20DeployedClaim{}
	_ EthereumClaim = &MsgLogicCallExecutedClaim{}
)

// GetType returns the type of the claim
func (msg *MsgDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (msg *MsgDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver)
	}
	if err := ValidateEthAddress(msg.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDepositClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDepositClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgDepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return val
}

// GetSigners defines whose signature is required
func (msg MsgDepositClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgDepositClaim) Type() string { return "deposit_claim" }

// Route should return the name of the module
func (msg MsgDepositClaim) Route() string { return RouterKey }

const (
	TypeMsgWithdrawClaim = "withdraw_claim"
)

// Hash implements BridgeDeposit.Hash
func (msg *MsgDepositClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", msg.TokenContract, string(msg.EthereumSender), msg.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the claim type
func (msg *MsgWithdrawClaim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *MsgWithdrawClaim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (msg *MsgWithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", msg.TokenContract, msg.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdrawClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgWithdrawClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgWithdrawClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgWithdrawClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWithdrawClaim) Type() string { return "withdraw_claim" }

const (
	TypeMsgDepositClaim = "deposit_claim"
)

// EthereumClaim implementation for MsgERC20DeployedClaim
// ======================================================

// GetType returns the type of the claim
func (e *MsgERC20DeployedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ERC20_DEPLOYED
}

// ValidateBasic performs stateless checks
func (e *MsgERC20DeployedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgERC20DeployedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgERC20DeployedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgERC20DeployedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgERC20DeployedClaim) Type() string { return "ERC20_deployed_claim" }

// Route should return the name of the module
func (msg MsgERC20DeployedClaim) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgERC20DeployedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", b.CosmosDenom, b.TokenContract, b.Name, b.Symbol, b.Decimals)
	return tmhash.Sum([]byte(path))
}

// EthereumClaim implementation for MsgLogicCallExecutedClaim
// ======================================================

// GetType returns the type of the claim
func (e *MsgLogicCallExecutedClaim) GetType() ClaimType {
	return CLAIM_TYPE_LOGIC_CALL_EXECUTED
}

// ValidateBasic performs stateless checks
func (e *MsgLogicCallExecutedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLogicCallExecutedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLogicCallExecutedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgLogicCallExecutedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgLogicCallExecutedClaim) Type() string { return "Logic_Call_Executed_Claim" }

// Route should return the name of the module
func (msg MsgLogicCallExecutedClaim) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgLogicCallExecutedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.InvalidationId, b.InvalidationNonce)
	return tmhash.Sum([]byte(path))
}

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgCancelSendToEth(val sdk.ValAddress, id uint64) *MsgCancelSendToEth {
	return &MsgCancelSendToEth{
		TransactionId: id,
	}
}

// Route should return the name of the module
func (msg *MsgCancelSendToEth) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCancelSendToEth) Type() string { return "cancel_send_to_eth" }

// ValidateBasic performs stateless checks
func (msg *MsgCancelSendToEth) ValidateBasic() (err error) {
	_, err = sdk.ValAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCancelSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}
