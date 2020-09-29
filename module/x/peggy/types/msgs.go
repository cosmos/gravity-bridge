package types

import (
	"encoding/hex"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

// ValsetConfirm
// this is the message sent by the validators when they wish to submit their signatures over
// the validator set at a given block height. A validator must first call MsgSetEthAddress to
// set their Ethereum address to be used for signing. Then someone (anyone) must make a ValsetRequest
// the request is essentially a messaging mechanism to determine which block all validators should submit
// signatures over. Finally validators sign the validator set, powers, and Ethereum addresses of the
// entire validator set at the height of a ValsetRequest and submit that signature with this message
// a ValsetConfirm.
//
// If a sufficient number of validators (66% of voting power) (A) have set Ethereum addresses and (B)
// submit ValsetConfirm messages with their signatures it is then possible for anyone to view these
// signatures in the chain store and submit them to Ethereum to update the validator set
// -------------
type MsgValsetConfirm struct {
	Nonce     UInt64Nonce    `json:"nonce"`
	Validator sdk.AccAddress `json:"validator"`
	Signature string         `json:"signature"`
}

func NewMsgValsetConfirm(nonce UInt64Nonce, validator sdk.AccAddress, signature string) MsgValsetConfirm {
	return MsgValsetConfirm{
		Nonce:     nonce,
		Validator: validator,
		Signature: signature,
	}
}

// Route should return the name of the module
func (msg MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetConfirm) Type() string { return "valset_confirm" }

// Stateless checks
func (msg MsgValsetConfirm) ValidateBasic() error {
	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}

// ValsetRequest
// This message starts off the validator set update process by coordinating a block height
// around which signatures over the validators, powers, and ethereum addresses will be made
// and submitted using a ValsetConfirm. Anyone can send this message as it is not authenticated
// in any way. In theory people could spam it and the validators will have to determine which
// block to actually coordinate around by looking over the valset requests and seeing which one
// some other validator has already submitted a ValsetResponse for.
// -------------
type MsgValsetRequest struct {
	Requester sdk.AccAddress `json:"requester"`
}

func NewMsgValsetRequest(requester sdk.AccAddress) MsgValsetRequest {
	return MsgValsetRequest{
		Requester: requester,
	}
}

// Route should return the name of the module
func (msg MsgValsetRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetRequest) Type() string { return "valset_request" }

func (msg MsgValsetRequest) ValidateBasic() error { return nil }

// GetSignBytes encodes the message for signing
func (msg MsgValsetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Requester}
}

// SetEthAddress
// This is used by the validators to set the Ethereum address that represents them on the
// Ethereum side of the bridge. They must sign their Cosmos address using the Ethereum address
// they have submitted.
// Like ValsetResponse this message can in theory be submitted by anyone, but only the current
// validator sets submissions carry any weight.
// -------------
type MsgSetEthAddress struct {
	// the ethereum address
	Address   EthereumAddress `json:"address"`
	Validator sdk.AccAddress  `json:"validator"`
	Signature string          `json:"signature"`
}

func NewMsgSetEthAddress(address EthereumAddress, validator sdk.AccAddress, signature string) MsgSetEthAddress {
	return MsgSetEthAddress{
		Address:   address,
		Validator: validator,
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
	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
	}

	if err := msg.Address.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	if msg.Address.IsEmpty() {
		return sdkerrors.Wrap(ErrEmpty, "ethereum address")
	}
	sigBytes, hexErr := hex.DecodeString(msg.Signature)
	if hexErr != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Could not decode hex string %s", msg.Signature))
	}

	err := utils.ValidateEthSig(crypto.Keccak256(msg.Validator.Bytes()), sigBytes, msg.Address.String())

	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("digest: %x sig: %x address %s error: %s", crypto.Keccak256(msg.Validator.Bytes()), msg.Signature, msg.Address, err.Error()))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetEthAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetEthAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}

// MsgSendToEth
// This is the message that a user calls when they want to bridge an asset
// TODO right now this needs to be locked to a single ERC20
// TODO fixed fee amounts for now, variable fee amounts in the fee field later
// TODO actually remove amounts form the users bank balances
// TODO this message modifies the on chain store by adding itself to a txpool
// it will later be removed when it is included in a batch and successfully submitted
// tokens are removed from the users balance immediately
// -------------
type MsgSendToEth struct {
	// the source address on Cosmos
	Sender sdk.AccAddress `json:"sender"`
	// the destination address on Ethereum
	DestAddress EthereumAddress `json:"dest_address"`
	// the coin to send across the bridge, note the restriction that this is a
	// single coin not a set of coins that is normal in other Cosmos messages
	Amount sdk.Coin `json:"send"`
	// the fee paid for the bridge, distinct from the fee paid to the chain to
	// actually send this message in the first place. So a successful send has
	// two layers of fees for the user
	BridgeFee sdk.Coin `json:"bridge_fee"`
}

func NewMsgSendToEth(sender sdk.AccAddress, destAddress EthereumAddress, send sdk.Coin, bridgeFee sdk.Coin) MsgSendToEth {
	return MsgSendToEth{
		Sender:      sender,
		DestAddress: destAddress,
		Amount:      send,
		BridgeFee:   bridgeFee,
	}
}

// Route should return the name of the module
func (msg MsgSendToEth) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToEth) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEth) ValidateBasic() error {
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
	// TODO validate eth address
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
	return []sdk.AccAddress{msg.Sender}
}

// MsgRequestBatch
// this is a message anyone can send that requests a batch of transactions to send across
// the bridge be created for whatever block height this message is included in. This acts as
// a coordination point, the handler for this message looks at the AddToOutgoingPool tx's in the store
// and generates a batch, also available in the store tied to this message. The validators then
// grab this batch, sign it, submit the signatures with a MsgConfirmBatch before a relayer can
// finally submit the batch
// -------------
type MsgRequestBatch struct {
	Requester sdk.AccAddress `json:"requester"`
	Denom     VoucherDenom   `json:"denom"`
}

func NewMsgRequestBatch(requester sdk.AccAddress) MsgRequestBatch {
	return MsgRequestBatch{
		Requester: requester,
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

func (msg MsgRequestBatch) ValidateBasic() error {
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
	return []sdk.AccAddress{msg.Requester}
}

// MsgConfirmBatch
// When validators observe a MsgRequestBatch they form a batch by ordering transactions currently
// in the txqueue in order of highest to lowest fee, cutting off when the batch either reaches a
// hardcoded maximum size (to be decided, probably around 100) or when transactions stop being
// profitable (TODO determine this without nondeterminism)
// This message includes the batch as well as an Ethereum signature over this batch by the validator
// -------------
type MsgConfirmBatch struct {
	Nonce        uint64         `json:"nonce"`
	Orchestrator sdk.AccAddress `json:"validator"`
	Signature    string         `json:"signature"`
}

func NewMsgConfirmBatch(nonce uint64, orchestrator sdk.AccAddress, signature string) MsgConfirmBatch {
	return MsgConfirmBatch{
		Nonce:        nonce,
		Orchestrator: orchestrator,
		Signature:    signature,
	}
}

// Route should return the name of the module
func (msg MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConfirmBatch) Type() string { return "confirm_batch" }

func (msg MsgConfirmBatch) ValidateBasic() error {
	// TODO validate signature
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
	return []sdk.AccAddress{msg.Orchestrator}
}

type EthereumClaim interface {
	GetNonce() UInt64Nonce
	GetType() ClaimType
	ValidateBasic() error
	Details() AttestationDetails
}

var (
	_ EthereumClaim = EthereumBridgeDepositClaim{}
	_ EthereumClaim = EthereumBridgeWithdrawalBatchClaim{}
	_ EthereumClaim = EthereumBridgeMultiSigUpdateClaim{}
	_ EthereumClaim = EthereumBridgeBootstrappedClaim{}
)

// NoUniqueClaimDetails is a NIL object to
var NoUniqueClaimDetails AttestationDetails = nil

// EthereumBridgeDepositClaim claims that a token was deposited on the bridge contract.
type EthereumBridgeDepositClaim struct {
	Nonce          UInt64Nonce `json:"nonce" yaml:"nonce"`
	ERC20Token     ERC20Token
	EthereumSender EthereumAddress `json:"ethereum_sender" yaml:"ethereum_sender"`
	CosmosReceiver sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
}

func (e EthereumBridgeDepositClaim) GetType() ClaimType {
	return ClaimTypeEthereumBridgeDeposit
}

func (e EthereumBridgeDepositClaim) GetNonce() UInt64Nonce {
	return e.Nonce
}

func (e EthereumBridgeDepositClaim) ValidateBasic() error {
	// todo: validate all fields
	if err := e.Nonce.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "nonce")
	}
	return nil
}

// EthereumBridgeDepositClaim
// When more than 66% of the active validator set has
// claimed to have seen the deposit enter the ethereum blockchain coins are issued
// to the Cosmos address in question
// -------------
func (e EthereumBridgeDepositClaim) Details() AttestationDetails {
	return BridgeDeposit{
		Nonce:          e.Nonce,
		ERC20Token:     e.ERC20Token,
		EthereumSender: e.EthereumSender,
		CosmosReceiver: e.CosmosReceiver,
	}
}

// EthereumBridgeWithdrawalBatchClaim claims that a batch of withdrawal operations on the bridge contract was executed.
type EthereumBridgeWithdrawalBatchClaim struct {
	Nonce UInt64Nonce `json:"nonce" yaml:"nonce"`
}

func (e EthereumBridgeWithdrawalBatchClaim) GetType() ClaimType {
	return ClaimTypeEthereumBridgeWithdrawalBatch
}

func (e EthereumBridgeWithdrawalBatchClaim) GetNonce() UInt64Nonce {
	return e.Nonce
}

func (e EthereumBridgeWithdrawalBatchClaim) ValidateBasic() error {
	if err := e.Nonce.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "nonce")
	}
	return nil
}

func (e EthereumBridgeWithdrawalBatchClaim) Details() AttestationDetails {
	return NoUniqueClaimDetails
}

// EthereumBridgeMultiSigUpdateClaim claims that the multisig set was updated on the bridge contract.
type EthereumBridgeMultiSigUpdateClaim struct {
	Nonce UInt64Nonce `json:"nonce" yaml:"nonce"`
}

func (e EthereumBridgeMultiSigUpdateClaim) GetType() ClaimType {
	return ClaimTypeEthereumBridgeMultiSigUpdate
}

func (e EthereumBridgeMultiSigUpdateClaim) GetNonce() UInt64Nonce {
	return e.Nonce
}

func (e EthereumBridgeMultiSigUpdateClaim) ValidateBasic() error {
	if err := e.Nonce.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "nonce")
	}
	return nil
}

func (e EthereumBridgeMultiSigUpdateClaim) Details() AttestationDetails {
	return NoUniqueClaimDetails
}

const (
	TypeMsgCreateEthereumClaims = "create_eth_claims"
)

var (
	_ sdk.Msg = &MsgCreateEthereumClaims{}
)

// EthereumBridgeBootstrappedClaim orchestrators confirm that the contract is setup on the Ethereum side and the init data.
type EthereumBridgeBootstrappedClaim struct {
	Nonce UInt64Nonce `json:"nonce" yaml:"nonce"`
	// AllowedValidatorSet addresses to participate
	AllowedValidatorSet []EthereumAddress
	// ValidatorPowers the validator's power values
	ValidatorPowers []uint64
	// PeggyID is a random 32 byte value to prevent signature reuse
	PeggyID []byte `json:"peggy_id,omitempty" yaml:"peggy_id"`
	// StartThreshold is the percentage of total voting power that must be online and participating in
	// Peggy operations before a bridge can start operating
	StartThreshold uint64 `json:"start_threshold,omitempty" yaml:"start_threshold"`
}

func (e EthereumBridgeBootstrappedClaim) GetNonce() UInt64Nonce {
	return e.Nonce
}

func (e EthereumBridgeBootstrappedClaim) GetType() ClaimType {
	return ClaimTypeEthereumBootstrap
}

func (e EthereumBridgeBootstrappedClaim) ValidateBasic() error {
	for i := range e.AllowedValidatorSet {
		if e.AllowedValidatorSet[i].IsEmpty() {
			return sdkerrors.Wrap(ErrEmpty, "ethereum address")
		}
	}
	for i := range e.ValidatorPowers {
		if e.ValidatorPowers[i] == 0 {
			return sdkerrors.Wrap(ErrEmpty, "power")
		}
	}
	if len(e.AllowedValidatorSet) != len(e.ValidatorPowers) {
		return sdkerrors.Wrap(ErrInvalid, "validator and power element count does not match")
	}
	// todo: implement me proper
	return nil
}

func (e EthereumBridgeBootstrappedClaim) Details() AttestationDetails {
	return BridgeBootstrap{
		AllowedValidatorSet: e.AllowedValidatorSet,
		ValidatorPowers:     e.ValidatorPowers,
		PeggyID:             e.PeggyID,
		StartThreshold:      e.StartThreshold,
	}
}

// MsgCreateEthereumClaims
// this message essentially acts as the oracle between Ethereum and Cosmos, when an orchestrator sees
// that a batch/ deposit/ multisig set update has been submitted on to the Ethereum blockchain they
// will submit this message which acts as their oracle attestation. When more than 66% of the active
// validator set has claimed to have seen the transaction enter the ethereum blockchain it is "observed"
// and state transitions and operations are triggered on the cosmos side.
type MsgCreateEthereumClaims struct {
	EthereumChainID       string // todo: revisit type. can be int/ string/ ?
	BridgeContractAddress EthereumAddress
	Orchestrator          sdk.AccAddress
	Claims                []EthereumClaim
}

func NewMsgCreateEthereumClaims(ethereumChainID string, bridgeContractAddress EthereumAddress, orchestrator sdk.AccAddress, claims []EthereumClaim) *MsgCreateEthereumClaims {
	return &MsgCreateEthereumClaims{EthereumChainID: ethereumChainID, BridgeContractAddress: bridgeContractAddress, Orchestrator: orchestrator, Claims: claims}
}

func (m MsgCreateEthereumClaims) Route() string {
	return RouterKey
}

func (m MsgCreateEthereumClaims) Type() string {
	return TypeMsgCreateEthereumClaims
}

func (m MsgCreateEthereumClaims) ValidateBasic() error {
	// todo: validate all fields
	if err := sdk.VerifyAddressFormat(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(err, "orchestrator")
	}
	for i := range m.Claims {
		if err := m.Claims[i].ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "claim %d", i)
		}
	}
	return nil
}

func (m MsgCreateEthereumClaims) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateEthereumClaims) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Orchestrator}
}
