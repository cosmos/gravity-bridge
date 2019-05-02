package types

import (
	"testing"
)

const (
	TestID           = "oracleID"
	TestByteString   = "{value: 5}"
	TestMinimumPower = 5
)

func CreateTestProphecy(t *testing.T) Prophecy {
	claim := NewClaim(TestID, []byte(TestByteString))
	claims := []Claim{claim}
	newProphecy := NewProphecy(TestID, PendingStatus, TestMinimumPower, claims)
	return newProphecy
}
