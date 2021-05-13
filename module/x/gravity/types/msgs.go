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
	_ sdk.Msg = &MsgSignerSetTxSignature{}
	_ sdk.Msg = &MsgSendToEthereum{}
	_ sdk.Msg = &MsgRequestBatchTx{}
	_ sdk.Msg = &MsgBatchTxSignature{}
	_ sdk.Msg = &MsgERC20DeployedEvent{}
	_ sdk.Msg = &MsgContractCallTxSignature{}
	_ sdk.Msg = &MsgContractCallExecutedEvent{}
	_ sdk.Msg = &MsgSendToCosmosEvent{}
	_ sdk.Msg = &MsgBatchExecutedEvent{}
	_ sdk.Msg = &MsgSubmitBadEthereumSignatureEvidence{}
)

// NewMsgDelegateKeys returns a new msgSetOrchestratorAddress
func NewMsgDelegateKeys(val sdk.ValAddress, oper sdk.AccAddress, eth string) *MsgDelegateKeys {
	return &MsgDelegateKeys{
		Validator:    val.String(),
		Orchestrator: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKeys) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKeys) Type() string { return "set_operator_address" }

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKeys) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
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
	acc, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// NewMsgSignerSetTxSignature returns a new msgSignerSetTxSignature
func NewMsgSignerSetTxSignature(
	nonce uint64,
	ethAddress string,
	validator sdk.AccAddress,
	signature string,
) *MsgSignerSetTxSignature {
	return &MsgSignerSetTxSignature{
		Nonce:        nonce,
		Orchestrator: validator.String(),
		EthAddress:   ethAddress,
		Signature:    signature,
	}
}

// Route should return the name of the module
func (msg *MsgSignerSetTxSignature) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSignerSetTxSignature) Type() string { return "signer_set_tx_signature" }

// ValidateBasic performs stateless checks
func (msg *MsgSignerSetTxSignature) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgSignerSetTxSignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSignerSetTxSignature) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToEthereum returns a new msgSendToEthereum
func NewMsgSendToEthereum(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEthereum {
	return &MsgSendToEthereum{
		Sender:    sender.String(),
		EthDest:   destAddress,
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
	if err := ValidateEthAddress(msg.EthDest); err != nil {
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
		Sender: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatchTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatchTx) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatchTx) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatchTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatchTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgBatchTxSignature) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBatchTxSignature) Type() string { return "batch_tx_signature" }

// ValidateBasic performs stateless checks
func (msg MsgBatchTxSignature) ValidateBasic() error {
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
func (msg MsgBatchTxSignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBatchTxSignature) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgContractCallTxSignature) Route() string { return RouterKey }

// Type should return the action
func (msg MsgContractCallTxSignature) Type() string { return "contract_call_tx_signature" }

// ValidateBasic performs stateless checks
func (msg MsgContractCallTxSignature) ValidateBasic() error {
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
func (msg MsgContractCallTxSignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgContractCallTxSignature) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// EthereumEvent represents a claim on ethereum state
type EthereumEvent interface {
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
	// the delegate address of the claimer, for MsgSendToCosmosEvent and MsgBatchExecutedEvent
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// Which type of claim this is
	GetType() EventType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ EthereumEvent = &MsgSendToCosmosEvent{}
	_ EthereumEvent = &MsgBatchExecutedEvent{}
	_ EthereumEvent = &MsgERC20DeployedEvent{}
	_ EthereumEvent = &MsgContractCallExecutedEvent{}
)

// GetType returns the type of the claim
func (msg *MsgSendToCosmosEvent) GetType() EventType {
	return EVENT_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (msg *MsgSendToCosmosEvent) ValidateBasic() error {
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
func (msg MsgSendToCosmosEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSendToCosmosEvent) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgSendToCosmosEvent failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return val
}

// GetSigners defines whose signature is required
func (msg MsgSendToCosmosEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgSendToCosmosEvent) Type() string { return "deposit_claim" }

// Route should return the name of the module
func (msg MsgSendToCosmosEvent) Route() string { return RouterKey }

const (
	TypeMsgBatchExecutedEvent = "withdraw_claim"
)

// Hash implements BridgeDeposit.Hash
func (msg *MsgSendToCosmosEvent) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", msg.TokenContract, string(msg.EthereumSender), msg.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the claim type
func (msg *MsgBatchExecutedEvent) GetType() EventType {
	return EVENT_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *MsgBatchExecutedEvent) ValidateBasic() error {
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
func (msg *MsgBatchExecutedEvent) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", msg.TokenContract, msg.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgBatchExecutedEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBatchExecutedEvent) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgBatchExecutedEvent failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgBatchExecutedEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgBatchExecutedEvent) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBatchExecutedEvent) Type() string { return "withdraw_claim" }

const (
	TypeMsgSendToCosmosEvent = "deposit_claim"
)

// EthereumEvent implementation for MsgERC20DeployedEvent
// ======================================================

// GetType returns the type of the claim
func (e *MsgERC20DeployedEvent) GetType() EventType {
	return EVENT_TYPE_COSMOS_ERC20_DEPLOYED
}

// ValidateBasic performs stateless checks
func (e *MsgERC20DeployedEvent) ValidateBasic() error {
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
func (msg MsgERC20DeployedEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgERC20DeployedEvent) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedEvent failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgERC20DeployedEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgERC20DeployedEvent) Type() string { return "ERC20_deployed_claim" }

// Route should return the name of the module
func (msg MsgERC20DeployedEvent) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgERC20DeployedEvent) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", b.CosmosDenom, b.TokenContract, b.Name, b.Symbol, b.Decimals)
	return tmhash.Sum([]byte(path))
}

// EthereumEvent implementation for MsgContractCallExecutedEvent
// ======================================================

// GetType returns the type of the claim
func (e *MsgContractCallExecutedEvent) GetType() EventType {
	return EVENT_TYPE_CONTRACT_CALL_EXECUTED
}

// ValidateBasic performs stateless checks
func (e *MsgContractCallExecutedEvent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgContractCallExecutedEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgContractCallExecutedEvent) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedEvent failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgContractCallExecutedEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgContractCallExecutedEvent) Type() string { return "Logic_Call_Executed_Claim" }

// Route should return the name of the module
func (msg MsgContractCallExecutedEvent) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgContractCallExecutedEvent) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.InvalidationId, b.InvalidationNonce)
	return tmhash.Sum([]byte(path))
}

// EthereumEvent implementation for MsgSignerSetUpdatedEvent
// ======================================================

// GetType returns the type of the claim
func (e *MsgSignerSetUpdatedEvent) GetType() EventType {
	return EVENT_TYPE_SIGNER_SET_UPDATED
}

// ValidateBasic performs stateless checks
func (e *MsgSignerSetUpdatedEvent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSignerSetUpdatedEvent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSignerSetUpdatedEvent) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedEvent failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgSignerSetUpdatedEvent) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgSignerSetUpdatedEvent) Type() string { return "SignerSetTx_Updated_Claim" }

// Route should return the name of the module
func (msg MsgSignerSetUpdatedEvent) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgSignerSetUpdatedEvent) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%s/", b.SignerSetNonce, b.EventNonce, b.BlockHeight, b.Members)
	return tmhash.Sum([]byte(path))
}

// NewMsgCancelSendToEthereum returns a new msgSetOrchestratorAddress
func NewMsgCancelSendToEthereum(val sdk.ValAddress, id uint64) *MsgCancelSendToEthereum {
	return &MsgCancelSendToEthereum{
		TransactionId: id,
	}
}

// Route should return the name of the module
func (msg *MsgCancelSendToEthereum) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCancelSendToEthereum) Type() string { return "cancel_send_to_eth" }

// ValidateBasic performs stateless checks
func (msg *MsgCancelSendToEthereum) ValidateBasic() (err error) {
	_, err = sdk.ValAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelSendToEthereum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCancelSendToEthereum) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// MsgSubmitBadEthereumSignatureEvidence
// ======================================================

// ValidateBasic performs stateless checks
func (e *MsgSubmitBadEthereumSignatureEvidence) ValidateBasic() error {
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSubmitBadEthereumSignatureEvidence) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSubmitBadEthereumSignatureEvidence) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

// Type should return the action
func (msg MsgSubmitBadEthereumSignatureEvidence) Type() string {
	return "Submit_Bad_Signature_Evidence"
}

// Route should return the name of the module
func (msg MsgSubmitBadEthereumSignatureEvidence) Route() string { return RouterKey }
