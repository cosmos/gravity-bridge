package stats

// ------------------------------------------------------------
//    Stats
//
//    Reports errors and tracks error/success statistics,
//		printing information and total time elapsed once an
//		error is thrown.
// ------------------------------------------------------------

import (
	"time"

	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/log"
)

// Stats
type Stats struct {
	numError      int
	numSuccessful int
	started       bool
	startedAt     time.Time
}

func NewStats() Stats {
	return Stats{0, 0, false, time.Now()}
}

// Add success to stats
func (s *Stats) AddSuccess() { s.numSuccessful++ }

// Add  error to stats
func (s *Stats) AddError() { s.numError++ }

// Prints the current stats
func (s *Stats) Print() {
	if !s.started {
		s.startedAt = time.Now()
		s.started = true
	}

	total := s.numSuccessful + s.numError
	secsPassed := time.Now().Sub(s.startedAt).Seconds()

	log.Log.Infof("========== Stats: Total %v, Error: %v, %%Successful: %v, TPS: %v\n", total, s.numError, float64(s.numSuccessful)/float64(total), float64(s.numSuccessful)/secsPassed)
}
