package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Nonce []byte

func NonceFromUint64(s uint64) Nonce {
	return sdk.Uint64ToBigEndian(s)
}

func (n Nonce) String() string {
	return string(n)
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

// Attestation is an aggregate of `claims` that eventually becomes `observed` by all orchestrators
type Attestation struct {
	ClaimType           ClaimType
	Nonce               Nonce // or bytes or int?
	Certainty           AttestationCertainty
	Status              AttestationProcessStatus
	ProcessResult       AttestationProcessResult
	Tally               AttestationTally
	SubmitTime          time.Time
	ConfirmationEndTime time.Time // votes collected <= end time. should be < unbonding period
	// ExpiryTime time.Time // todo: do we want to keep Attestations forever persisted or can we delete them?
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

// ID is the unique identifiert used in DB
func (a *Attestation) ID() []byte {
	return GetAttestationKey(a.ClaimType, a.Nonce)
}
