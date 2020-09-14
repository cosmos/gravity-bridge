package types

import (
	"fmt"
	"math/big"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Nonce []byte

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
	//ProcessStatusTimeout   AttestationProcessStatus = 3 // who sets it if we return errors?
)

type AttestationProcessResult uint8

const (
	ProcessResultUnknown AttestationProcessResult = 0
	ProcessResultSuccess AttestationProcessResult = 1
	ProcessResultFailure AttestationProcessResult = 2
)

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
	TotalVotesPower    uint64 // can this overflow?
	TotalVotesCount    uint64
	RequiredVotesPower uint64 // todo: revisit if the assumption is true that we can use the values from first claim timestamp
	RequiredVotesCount uint64 // todo: revisit
}

func (t *AttestationTally) addVote(power uint64) {
	t.TotalVotesCount += 1
	t.TotalVotesPower += power
}

func (t AttestationTally) thresholdsReached() bool {
	return t.TotalVotesPower > t.RequiredVotesPower &&
		t.TotalVotesCount > t.RequiredVotesCount
}

func (a *Attestation) AddConfirmation(now time.Time, power uint64) error {
	if a.Status != ProcessStatusInit {
		return sdkerrors.Wrapf(ErrInvalidState, "%d", a.Status) // no status to string impl, yet
	}
	if now.After(a.ConfirmationEndTime) {
		return ErrTimeout
	}
	a.Tally.addVote(power)
	if a.Tally.thresholdsReached() {
		a.Certainty = CertaintyObserved
	}
	return nil
}

// The Fraction type represents a numerator and denominator to enable higher precision thresholds in
// the election rules. For example:
// numerator: 1, denominator: 2 => > 50%
// numerator: 2, denominator: 3 => > 66.666..%
// numerator: 6273, denominator: 10000 => > 62.73%
// Valid range of the fraction is 0.5 to 1.
type Fraction struct {
	// The top number in a fraction.
	Numerator uint32
	// The bottom number
	Denominator uint32
}

// mul multiply
func (f Fraction) Mul(factor uint64) uint64 {
	a := new(big.Int).Mul(big.NewInt(int64(factor)), big.NewInt(int64(f.Numerator)))
	r := new(big.Int).Div(a, big.NewInt(int64(f.Denominator)))
	fmt.Printf("%d *%d / %d = %d \n", f.Numerator, factor, f.Denominator, r.Uint64())
	return r.Uint64() // this is where rounding happens
}
