package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestID                  = "oracleID"
	TestByteString          = "{value: 5}"
	AlternateTestByteString = "{value: 7}"
)

func CreateTestProphecy(validator sdk.AccAddress) Prophecy {
	claim := CreateTestClaimForValidator(validator)
	claims := []Claim{claim}
	newProphecy := NewProphecy(TestID, PendingStatus, claims)
	return newProphecy
}

func CreateTestClaimForValidator(validator sdk.AccAddress) Claim {
	claim := NewClaim(TestID, validator, []byte(TestByteString))
	return claim
}

func CreateAlternateTestClaimForValidator(validator sdk.AccAddress) Claim {
	claim := NewClaim(TestID, validator, []byte(AlternateTestByteString))
	return claim
}
