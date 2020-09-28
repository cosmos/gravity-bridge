package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// nonce defines an abstract nonce type that is unique within it's context.
type nonce interface {
	// String returns an encoded human readable representation. Used in URLs.
	String() string
	// Bytes returns an encoded raw bytes representation
	Bytes() []byte
	// ValidateBasic returns the result of the syntax check
	ValidateBasic() error
	// GreaterThan than other.
	GreaterThan(o nonce) bool
	IsEmpty() bool
}

type UInt64Nonce uint64

var _ nonce = NewUInt64Nonce(0)

func NewUInt64Nonce(s uint64) UInt64Nonce {
	return UInt64Nonce(s)
}

// UInt64NonceFromBytes create UInt64Nonce from binary big endian representation
func UInt64NonceFromBytes(s []byte) UInt64Nonce {
	return NewUInt64Nonce(binary.BigEndian.Uint64(s))
}

// UInt64NonceFromBytes create UInt64Nonce from human readable string representation
func UInt64NonceFromString(s string) (UInt64Nonce, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return NewUInt64Nonce(0), err
	}
	return NewUInt64Nonce(v), nil
}

func Uint64FromNonce(n nonce) (uint64, error) {
	if v, ok := n.(nonceUint64er); ok {
		return v.Uint64(), nil
	}
	return 0, ErrUnsupported
}

func (n UInt64Nonce) Uint64() uint64 {
	return uint64(n)
}

func (n UInt64Nonce) String() string {
	return strconv.FormatUint(n.Uint64(), 10)
}

func (n UInt64Nonce) Bytes() []byte {
	return sdk.Uint64ToBigEndian(n.Uint64())
}

func (n UInt64Nonce) ValidateBasic() error {
	if n.IsEmpty() {
		return ErrEmpty
	}
	return nil
}

type nonceUint64er interface {
	Uint64() uint64
}

func (n UInt64Nonce) GreaterThan(o nonce) bool {
	if o == nil || reflect.ValueOf(o).IsZero() || o.IsEmpty() {
		return true
	}
	if n.IsEmpty() {
		return false
	}
	if v, ok := o.(nonceUint64er); ok {
		return n.Uint64() > v.Uint64()
	}
	return bytes.Compare(n.Bytes(), o.Bytes()) == 1
}

func (n UInt64Nonce) IsEmpty() bool {
	return n == 0
}

type AttestationCertainty uint8

const (
	CertaintyUnknown   AttestationCertainty = 0
	CertaintyRequested AttestationCertainty = 1
	CertaintyObserved  AttestationCertainty = 2
)

type AttestationProcessStatus uint8

const (
	ProcessStatusUnknown   AttestationProcessStatus = 0
	ProcessStatusInit      AttestationProcessStatus = 1
	ProcessStatusProcessed AttestationProcessStatus = 2 // prevent double processing
	ProcessStatusTimeout   AttestationProcessStatus = 3 // end block process will set this
)

type AttestationProcessResult uint8

const (
	ProcessResultUnknown AttestationProcessResult = 0
	ProcessResultSuccess AttestationProcessResult = 1
	ProcessResultFailure AttestationProcessResult = 2
)

// ClaimType is the cosmos type of an event from the counterpart chain that can be handled
type ClaimType string

const ( // todo: revisit type: length and overlap
	// oracles
	ClaimTypeEthereumBridgeDeposit         ClaimType = "bridge_deposit"
	ClaimTypeEthereumBridgeWithdrawalBatch ClaimType = "bridge_withdrawal_batch"
	ClaimTypeEthereumBridgeMultiSigUpdate  ClaimType = "bridge_multisig_update"
	ClaimTypeEthereumBootstrap             ClaimType = "bridge_bootstrap"

	// signed confirmations to Ethereum
	ClaimTypeOrchestratorSignedMultiSigUpdate ClaimType = "orchestrator_signed_multisig_update"
	ClaimTypeOrchestratorSignedWithdrawBatch  ClaimType = "orchestrator_signed_withdraw_batch"
)

var AllOracleClaimTypes = []ClaimType{ClaimTypeEthereumBridgeDeposit, ClaimTypeEthereumBridgeWithdrawalBatch, ClaimTypeEthereumBridgeMultiSigUpdate, ClaimTypeEthereumBootstrap}
var AllConfirmationClaimTypes = []ClaimType{ClaimTypeOrchestratorSignedMultiSigUpdate, ClaimTypeOrchestratorSignedWithdrawBatch}

func IsClaimType(s string) bool {
	for _, v := range append(AllOracleClaimTypes, AllConfirmationClaimTypes...) {
		if string(v) == s {
			return true
		}
	}
	return false
}

func (c ClaimType) String() string {
	return string(c)
}

func (c ClaimType) Bytes() []byte {
	return []byte(c)
}

// Attestation is an aggregate of `claims` that eventually becomes `observed` by all orchestrators
type Attestation struct {
	ClaimType           ClaimType
	Nonce               UInt64Nonce
	Certainty           AttestationCertainty
	Status              AttestationProcessStatus
	ProcessResult       AttestationProcessResult
	Tally               AttestationTally
	SubmitTime          time.Time
	ConfirmationEndTime time.Time // votes collected <= end time. should be < unbonding period
	// ExpiryTime time.Time // todo: do we want to keep Attestations forever persisted or can we delete them?
	Details AttestationDetails
}

type AttestationTally struct {
	TotalVotesPower    sdk.Uint
	TotalVotesCount    uint64
	RequiredVotesPower sdk.Uint // todo: revisit if the assumption is true that we can use the values from first claim timestamp
	RequiredVotesCount uint64   // todo: revisit as above
}

func (t *AttestationTally) addVote(power uint64) {
	t.TotalVotesCount += 1
	t.TotalVotesPower = t.TotalVotesPower.AddUint64(power)
}

// ThresholdsReached returns true when votes power > 66% of the validators AND total votes > 50% of the validator count
func (t AttestationTally) ThresholdsReached() bool {
	return t.TotalVotesPower.GT(t.RequiredVotesPower) &&
		t.TotalVotesCount > t.RequiredVotesCount
}

func (a *Attestation) AddVote(now time.Time, power uint64) error {
	if a.Status != ProcessStatusInit {
		return sdkerrors.Wrapf(ErrInvalid, "%d", a.Status) // no status to string impl, yet
	}
	if now.After(a.ConfirmationEndTime) {
		return ErrTimeout
	}
	a.Tally.addVote(power)
	if a.Tally.ThresholdsReached() {
		a.Certainty = CertaintyObserved
	}
	return nil
}

// ID is the unique identifier used in DB
func (a *Attestation) ID() []byte {
	return GetAttestationKey(a.ClaimType, a.Nonce)
}

// AttestationDetails is the payload of an attestation.
type AttestationDetails interface {
	// Hash creates hash of the object that is supposed to be unique during the live time of the block chain.
	// purpose of the hash is to very that orchestrators submit the same payload data and not only the nonce.
	Hash() []byte
}

var (
	_ AttestationDetails = BridgeDeposit{}
	_ AttestationDetails = SignedCheckpoint{}
	_ AttestationDetails = BridgeBootstrap{}
)

// BridgeDeposit is an attestation detail that adds vouchers to an account when executed
type BridgeDeposit struct {
	Nonce          UInt64Nonce // redundant information but required for a unique hash. Two deposits should not have the same hash.
	ERC20Token     ERC20Token
	EthereumSender EthereumAddress `json:"ethereum_sender" yaml:"ethereum_sender"`
	CosmosReceiver sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
}

func (b BridgeDeposit) Hash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/", b.Nonce.String(), b.EthereumSender.String(), b.ERC20Token.String(), b.CosmosReceiver.String())
	return tmhash.Sum([]byte(path))
}

// SignedCheckpoint is an attestation detail that approves an update for a checkpoint
type SignedCheckpoint struct {
	Checkpoint []byte // is a hash already
}

func (s SignedCheckpoint) Hash() []byte {
	return s.Checkpoint
}

type BridgeBootstrap struct {
	AllowedValidatorSet []EthereumAddress
	ValidatorPowers     []uint64
	PeggyID             []byte `json:"peggy_id,omitempty" yaml:"peggy_id"`
	StartThreshold      uint64 `json:"start_threshold,omitempty" yaml:"start_threshold"`
}

func (b BridgeBootstrap) Hash() []byte {
	hasher := tmhash.New()
	for i := range b.AllowedValidatorSet {
		_, err := hasher.Write(b.AllowedValidatorSet[i].RawBytes())
		if err != nil {
			panic(err) // can not happen in used sha256 impl
		}
	}
	for i := range b.ValidatorPowers {
		_, err := hasher.Write(sdk.Uint64ToBigEndian(b.ValidatorPowers[i]))
		if err != nil {
			panic(err) // can not happen in used sha256 impl
		}
	}
	hasher.Write(b.PeggyID)
	hasher.Write(sdk.Uint64ToBigEndian(b.StartThreshold))
	return hasher.Sum(nil)
}
