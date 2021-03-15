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
	_ sdk.Msg = &MsgSendToEth{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgSubmitConfirm{}
	_ sdk.Msg = &MsgSubmitClaim{}
)

// NewMsgSetOrchestratorAddress returns a new msgDelegateKeys
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

// NewMsgSendToEth returns a new msgSendToEth
func NewMsgSendToEth(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEth {
	return &MsgSendToEth{
		Sender:    sender.String(),
		EthDest:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
	}
}

// Route should return the name of the module
func (msg MsgSendToEth) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToEth) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
	if msg.Amount.Denom != msg.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
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
func (msg MsgSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(validator sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		Validator: validator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgSubmitConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSubmitConfirm) Type() string { return "confirm_logic" }

// ValidateBasic performs stateless checks
func (msg MsgSubmitConfirm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
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
func (msg MsgSubmitConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSubmitConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

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
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

type Confirm interface {
	// TODO: define
}

var (
	_ EthereumClaim = &MsgSubmitClaim{}
)

// GetType returns the claim type
func (e *MsgSubmitClaim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *MsgSubmitClaim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Signer)
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (b *MsgSubmitClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.TokenContract, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgSubmitClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSubmitClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgSubmitClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.Signer)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgSubmitClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSubmitClaim) Type() string { return "withdraw_claim" }

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
