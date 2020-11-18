package types

import (
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	_ sdk.Msg = &MsgCreateEthereumClaims{}
	_ sdk.Msg = &MsgValsetConfirm{}
	_ sdk.Msg = &MsgValsetRequest{}
	_ sdk.Msg = &MsgSetEthAddress{}
	_ sdk.Msg = &MsgSendToEth{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgConfirmBatch{}
	_ sdk.Msg = &MsgCreateEthereumClaims{}
	_ sdk.Msg = &MsgBridgeSignatureSubmission{}
)

// NewMsgBridgeSignatureSubmission returns a new msgBridgeSignatureSubmission
func NewMsgBridgeSignatureSubmission(signtype SignType, nonce uint64, orch, ethsig string) *MsgBridgeSignatureSubmission {
	return &MsgBridgeSignatureSubmission{
		SignType:          signtype,
		Nonce:             nonce,
		Orchestrator:      orch,
		EthereumSignature: ethsig,
	}
}

// Route should return the name of the module
func (msg *MsgBridgeSignatureSubmission) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgBridgeSignatureSubmission) Type() string { return "bridge_signature_submission" }

// ValidateBasic performs stateless checks
func (msg *MsgBridgeSignatureSubmission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if _, err := hex.DecodeString(msg.EthereumSignature); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.EthereumSignature)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgBridgeSignatureSubmission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgBridgeSignatureSubmission) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgValsetConfirm returns a new msgValsetConfirm
func NewMsgValsetConfirm(nonce uint64, ethAddress EthereumAddress, validator sdk.AccAddress, signature string) *MsgValsetConfirm {
	return &MsgValsetConfirm{
		Nonce:      nonce,
		Validator:  validator.String(),
		EthAddress: ethAddress.String(),
		Signature:  signature,
	}
}

// Route should return the name of the module
func (msg *MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgValsetConfirm) Type() string { return "valset_confirm" }

// ValidateBasic performs stateless checks
func (msg *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := NewEthereumAddress(msg.EthAddress).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgValsetRequest returns a new msgValsetRequest
func NewMsgValsetRequest(requester sdk.AccAddress) *MsgValsetRequest {
	return &MsgValsetRequest{
		Requester: requester.String(),
	}
}

// Route should return the name of the module
func (msg MsgValsetRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetRequest) Type() string { return "valset_request" }

// ValidateBasic performs stateless checks
func (msg MsgValsetRequest) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Requester); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetRequest) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSetEthAddress return a new msgSetEthAddress
// TODO: figure out if we need sdk.ValAddress here
func NewMsgSetEthAddress(address EthereumAddress, validator sdk.AccAddress, signature string) *MsgSetEthAddress {
	return &MsgSetEthAddress{
		Address:   address.String(),
		Validator: validator.String(),
		Signature: signature,
	}
}

// Route should return the name of the module
func (msg MsgSetEthAddress) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetEthAddress) Type() string { return "set_eth_address" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid, and whether the Eth address has signed the validator address
// (proving control of the Eth address)
func (msg MsgSetEthAddress) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := NewEthereumAddress(msg.Address).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	err = ValidateEthereumSignature(crypto.Keccak256([]byte(msg.Validator)), sigBytes, string(msg.Address))
	if err != nil {
		return sdkerrors.Wrapf(err, "digest: %x sig: %x address %s error: %s", crypto.Keccak256([]byte(msg.Validator)), msg.Signature, msg.Address, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetEthAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetEthAddress) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToEth returns a new msgSendToEth
func NewMsgSendToEth(sender sdk.AccAddress, destAddress EthereumAddress, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEth {
	return &MsgSendToEth{
		Sender:    sender.String(),
		EthDest:   destAddress.String(),
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee and amount must be the same type")
	}
	if !IsVoucherDenom(msg.Amount.Denom) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is not a voucher type")
	}
	if !IsVoucherDenom(msg.BridgeFee.Denom) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee is not a voucher type")
	}
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := NewEthereumAddress(msg.EthDest).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO for demo get single allowed demon from the store
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
func NewMsgRequestBatch(requester sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		Requester: requester.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Requester); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester)
	}
	// TODO ensure that Demon matches hardcoded allowed value
	// TODO later make sure that Demon matches a list of tokens already
	// in the bridge to send
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Requester)
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
	if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := NewEthereumAddress(msg.EthSigner).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := NewEthereumAddress(msg.TokenContract).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	err = ValidateEthereumSignature(crypto.Keccak256([]byte(msg.Validator)), sigBytes, string(msg.EthSigner))
	if err != nil {
		return sdkerrors.Wrapf(err, "digest: %x sig: %x address %s error: %s", crypto.Keccak256([]byte(msg.Validator)), msg.Signature, msg.EthSigner, err.Error())
	}

	// TODO get batch from storage
	// TODO generate batch in storage on MsgRequestBatch in the first place
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// EthereumClaim represenets a claim on ethereum state
type EthereumClaim interface {
	GetEventNonce() uint64
	GetType() ClaimType
	ValidateBasic() error
	Details() AttestationDetails
}

var (
	_ EthereumClaim = &EthereumBridgeDepositClaim{}
	_ EthereumClaim = &EthereumBridgeWithdrawalBatchClaim{}
)

// NoUniqueClaimDetails is a NIL object to
var NoUniqueClaimDetails AttestationDetails = nil

// GetType returns the type of the claim
func (e *EthereumBridgeDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_ETHEREUM_BRIDGE_DEPOSIT
}

// GetEventNonce returns the event nonce for the claim
func (e *EthereumBridgeDepositClaim) GetEventNonce() uint64 {
	return e.Nonce
}

// ValidateBasic performs stateless checks
func (e *EthereumBridgeDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	if err := NewEthereumAddress(e.EthereumSender).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := e.Erc20Token.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if e.Nonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// Details returns the attestation details fromt he bridge deposit claim
func (e *EthereumBridgeDepositClaim) Details() AttestationDetails {
	return &BridgeDeposit{
		Erc20Token:     e.Erc20Token,
		EthereumSender: e.EthereumSender,
		CosmosReceiver: e.CosmosReceiver,
	}
}

// GetType returns the claim type
func (e *EthereumBridgeWithdrawalBatchClaim) GetType() ClaimType {
	return CLAIM_TYPE_ETHEREUM_BRIDGE_WITHDRAWAL_BATCH
}

// ValidateBasic performs stateless checks
func (e *EthereumBridgeWithdrawalBatchClaim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	return nil
}

// Details returns the attestation details from the claim
func (e *EthereumBridgeWithdrawalBatchClaim) Details() AttestationDetails {
	return NoUniqueClaimDetails
}

const (
	// TypeMsgCreateEthereumClaims is the claim type
	TypeMsgCreateEthereumClaims = "create_eth_claims"
)

// NewMsgCreateEthereumClaims returns a new msgCreateEthereumClaims
func NewMsgCreateEthereumClaims(ethereumChainID uint64, bridgeContractAddress EthereumAddress, orchestrator sdk.AccAddress, claims []EthereumClaim) *MsgCreateEthereumClaims {
	var packedClaims []*codectypes.Any
	for _, c := range claims {
		pc, err := PackEthereumClaim(c)
		if err != nil {
			panic(err)
		}
		packedClaims = append(packedClaims, pc)
	}
	return &MsgCreateEthereumClaims{EthereumChainId: ethereumChainID, BridgeContractAddress: bridgeContractAddress.String(), Orchestrator: orchestrator.String(), Claims: packedClaims}
}

// Route returns the route for the msg
func (m MsgCreateEthereumClaims) Route() string {
	return RouterKey
}

// Type returns the type of msg
func (m MsgCreateEthereumClaims) Type() string {
	return TypeMsgCreateEthereumClaims
}

// ValidateBasic performs stateless checks on the message
func (m MsgCreateEthereumClaims) ValidateBasic() error {
	// todo: validate ethereum chain id
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if err := NewEthereumAddress(m.BridgeContractAddress).ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	for i := range m.Claims {
		claim, err := UnpackEthereumClaim(m.Claims[i])
		if err != nil {
			return sdkerrors.Wrapf(err, "claim %d failed to unpack any", i)
		}
		if err := claim.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "claim %d failed ValidateBasic()", i)
		}
	}
	return nil
}

// GetSignBytes returns the bytes to sign over
// TODO: deperecate GetSignBytes methods
func (m MsgCreateEthereumClaims) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the signers of the message
func (m MsgCreateEthereumClaims) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}
