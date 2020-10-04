package types

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type AttestationCertainty uint8

const (
	CertaintyUnknown   AttestationCertainty = 0
	CertaintyRequested AttestationCertainty = 1
	CertaintyObserved  AttestationCertainty = 2
)

var certaintyToNames = map[AttestationCertainty]string{
	CertaintyUnknown:   "unknown",
	CertaintyRequested: "requested",
	CertaintyObserved:  "observed",
}

func (c AttestationCertainty) String() string {
	return certaintyToNames[c]
}

func (c AttestationCertainty) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", c.String())), nil
}

func (c *AttestationCertainty) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	for k, v := range certaintyToNames {
		if s == v {
			*c = k
			return nil
		}
	}
	return sdkerrors.Wrap(ErrUnknown, "certainty")
}

type AttestationProcessStatus uint8

const (
	ProcessStatusUnknown   AttestationProcessStatus = 0
	ProcessStatusInit      AttestationProcessStatus = 1
	ProcessStatusProcessed AttestationProcessStatus = 2 // prevent double processing
	ProcessStatusTimeout   AttestationProcessStatus = 3 // end block process will set this
)

var attestationProcessStatusToNames = map[AttestationProcessStatus]string{
	ProcessStatusUnknown:   "unknown",
	ProcessStatusInit:      "init",
	ProcessStatusProcessed: "processed",
	ProcessStatusTimeout:   "timeout",
}

func (c AttestationProcessStatus) String() string {
	return attestationProcessStatusToNames[c]
}

func (c AttestationProcessStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", c.String())), nil
}

func (c *AttestationProcessStatus) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	for k, v := range attestationProcessStatusToNames {
		if s == v {
			*c = k
			return nil
		}
	}
	return sdkerrors.Wrap(ErrUnknown, "process status")
}

type AttestationProcessResult uint8

const (
	ProcessResultUnknown AttestationProcessResult = 0
	ProcessResultSuccess AttestationProcessResult = 1
	ProcessResultFailure AttestationProcessResult = 2
)

var attestationProcessResultToNames = map[AttestationProcessResult]string{
	ProcessResultUnknown: "unknown",
	ProcessResultSuccess: "success",
	ProcessResultFailure: "failure",
}

func (c AttestationProcessResult) String() string {
	return attestationProcessResultToNames[c]
}

func (c AttestationProcessResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", c.String())), nil
}

func (c *AttestationProcessResult) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	for k, v := range attestationProcessResultToNames {
		if s == v {
			*c = k
			return nil
		}
	}
	return sdkerrors.Wrap(ErrUnknown, "process result")
}

// ClaimType is the cosmos type of an event from the counterpart chain that can be handled
type ClaimType byte

const (
	ClaimTypeUnknown ClaimType = 0
	// oracles
	ClaimTypeEthereumBridgeDeposit ClaimType = 1
	// a withdraw batch was executed on the Ethereum side
	ClaimTypeEthereumBridgeWithdrawalBatch ClaimType = 2
	ClaimTypeEthereumBridgeMultiSigUpdate  ClaimType = 3
	ClaimTypeEthereumBridgeBootstrap       ClaimType = 4

	// signed confirmations on cosmos for Ethereum side
	ClaimTypeOrchestratorSignedMultiSigUpdate ClaimType = 5
	ClaimTypeOrchestratorSignedWithdrawBatch  ClaimType = 6
)

var claimTypeToNames = map[ClaimType]string{
	ClaimTypeEthereumBridgeDeposit:            "bridge_deposit",
	ClaimTypeEthereumBridgeWithdrawalBatch:    "bridge_withdrawal_batch",
	ClaimTypeEthereumBridgeMultiSigUpdate:     "bridge_multisig_update",
	ClaimTypeEthereumBridgeBootstrap:          "bridge_bootstrap",
	ClaimTypeOrchestratorSignedMultiSigUpdate: "orchestrator_signed_multisig_update",
	ClaimTypeOrchestratorSignedWithdrawBatch:  "orchestrator_signed_withdraw_batch",
}

// AllOracleClaimTypes types that are observed and submitted by the current orchestrator set
var AllOracleClaimTypes = []ClaimType{ClaimTypeEthereumBridgeDeposit, ClaimTypeEthereumBridgeWithdrawalBatch, ClaimTypeEthereumBridgeMultiSigUpdate, ClaimTypeEthereumBridgeBootstrap}

// AllSignerApprovalClaimTypes types that are signed with by the bridge multisig set
var AllSignerApprovalClaimTypes = []ClaimType{ClaimTypeOrchestratorSignedMultiSigUpdate, ClaimTypeOrchestratorSignedWithdrawBatch}

func ClaimTypeFromName(s string) (ClaimType, bool) {
	for _, v := range append(AllOracleClaimTypes, AllSignerApprovalClaimTypes...) {
		name, ok := claimTypeToNames[v]
		if ok && name == s {
			return v, true
		}
	}
	return ClaimTypeUnknown, false
}
func ToClaimTypeNames(s ...ClaimType) []string {
	r := make([]string, len(s))
	for i := range s {
		r[i] = s[i].String()
	}
	return r
}

func IsSignerApprovalClaimType(s ClaimType) bool {
	for _, v := range AllSignerApprovalClaimTypes {
		if v == s {
			return true
		}
	}
	return false
}

func IsOracleObservationClaimType(s ClaimType) bool {
	for _, v := range AllOracleClaimTypes {
		if v == s {
			return true
		}
	}
	return false
}

func (c ClaimType) String() string {
	return claimTypeToNames[c]
}

func (c ClaimType) Bytes() []byte {
	return []byte{byte(c)}
}

func (e ClaimType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", e.String())), nil
}

func (e *ClaimType) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	c, exists := ClaimTypeFromName(s)
	if !exists {
		return sdkerrors.Wrap(ErrUnknown, "claim type")
	}
	*e = c
	return nil
}

// Attestation is an aggregate of `claims` that eventually becomes `observed` by all orchestrators
type Attestation struct {
	ClaimType           ClaimType                `json:"claim_type"`
	Nonce               UInt64Nonce              `json:"nonce"`
	Certainty           AttestationCertainty     `json:"certainty"`
	Status              AttestationProcessStatus `json:"status"`
	ProcessResult       AttestationProcessResult `json:"process_result"`
	Tally               AttestationTally         `json:"tally"`
	SubmitTime          time.Time                `json:"submit_time"`
	ConfirmationEndTime time.Time                `json:"confirmation_end_time"` // votes collected <= end time. should be < unbonding period
	// ExpiryTime time.Time // todo: do we want to keep Attestations forever persisted or can we delete them?
	Details AttestationDetails `json:"details,omitempty"`
}

type AttestationTally struct {
	TotalVotesPower    sdk.Uint `json:"total_votes_power"`
	TotalVotesCount    uint64   `json:"total_votes_count"`
	RequiredVotesPower sdk.Uint `json:"required_votes_power"` // todo: revisit if the assumption is true that we can use the values from first claim timestamp
	RequiredVotesCount uint64   `json:"required_votes_count"` // todo: revisit as above
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
	Nonce          UInt64Nonce     `json:"nonce"` // redundant information but required for a unique hash. Two deposits should not have the same hash.
	ERC20Token     ERC20Token      `json:"erc_20_token"`
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
	AllowedValidatorSet []EthereumAddress `json:"allowed_validator_set"`
	ValidatorPowers     []uint64          `json:"validator_powers"`
	PeggyID             string            `json:"peggy_id" yaml:"peggy_id"`
	StartThreshold      uint64            `json:"start_threshold,omitempty" yaml:"start_threshold"`
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
	hasher.Write([]uint8(b.PeggyID))
	hasher.Write(sdk.Uint64ToBigEndian(b.StartThreshold))
	return hasher.Sum(nil)
}
