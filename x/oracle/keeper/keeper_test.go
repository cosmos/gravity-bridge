package keeper

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

func TestCreateGetProphecy(t *testing.T) {
	ctx, _, keeper, validatorAddresses, _ := CreateTestKeepers(t, false, 0.7, []int64{10})

	validator1Pow3 := validatorAddresses[0]

	//Test normal Creation
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test bad Creation with blank id
	status, err = keeper.ProcessClaim(ctx, "", validator1Pow3, types.TestString)
	require.Error(t, err)

	//Test bad Creation with blank claim
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, "")
	require.Error(t, err)

	//Test retrieval
	prophecy, err := keeper.GetProphecy(ctx, types.TestID)
	require.NoError(t, err)
	require.Equal(t, prophecy.ID, types.TestID)
	require.Equal(t, prophecy.Status.StatusText, types.PendingStatusText)
	require.Equal(t, prophecy.ClaimValidators[types.TestString][0], validator1Pow3)
	require.Equal(t, prophecy.ValidatorClaims[validator1Pow3.String()], types.TestString)
}

func TestBadConsensusForOracle(t *testing.T) {
	_, _, _, _, err := CreateTestKeepers(t, false, 0, []int64{10})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "consensus proportion of validator staking power must be > 0 and <= 1"))

	_, _, _, _, err = CreateTestKeepers(t, false, 1.2, []int64{10})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "consensus proportion of validator staking power must be > 0 and <= 1"))
}

func TestBadMsgs(t *testing.T) {
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.6, []int64{3, 3})
	validator1Pow3 := validatorAddresses[0]

	//Test empty claim
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, "")
	require.Error(t, err)
	require.Equal(t, status.FinalClaim, "")
	require.True(t, strings.Contains(err.Error(), "Claim cannot be empty string"))

	//Test normal Creation
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test duplicate message
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Already processed message from validator for this id"))

	//Test second but non duplicate message
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.AlternateTestString)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Already processed message from validator for this id"))

}

func TestSuccessfulProphecy(t *testing.T) {
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.6, []int64{3, 3, 4})
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test second claim completes and finalizes to success
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.TestString)

	//Test third claim not possible
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator3Pow4, types.TestString)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Prophecy already finalized"))
}

func TestSuccessfulProphecyWithDisagreement(t *testing.T) {
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.6, []int64{3, 3, 4})
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test second disagreeing claim processed fine
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow3, types.AlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test third claim agrees and finalizes to success
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator3Pow4, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.TestString)
}

func TestFailedProphecy(t *testing.T) {
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.6, []int64{3, 3, 4})
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test second disagreeing claim processed fine
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow3, types.AlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)
	require.Equal(t, status.FinalClaim, "")

	//Test third disagreeing claim processed fine and prophecy fails
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator3Pow4, types.AnotherAlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.FailedStatusText)
	require.Equal(t, status.FinalClaim, "")
}

func TestPower(t *testing.T) {
	//Testing with 2 validators but one has high enough power to overrule
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.7, []int64{3, 7})
	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test first claim
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test second disagreeing claim processed fine and finalized to its bytes
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow7, types.AlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.AlternateTestString)

	//Test alternate power setup with validators of 5/4/3/9 and total power 22 and 12/21 required
	ctx, _, keeper, validatorAddresses, err = CreateTestKeepers(t, false, 0.571, []int64{5, 4, 3, 9})
	validator1Pow5 := validatorAddresses[0]
	validator2Pow4 := validatorAddresses[1]
	validator3Pow3 := validatorAddresses[2]
	validator4Pow9 := validatorAddresses[3]

	//Test claim by v1
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator1Pow5, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test claim by v2
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow4, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test alternate claim by v4
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator4Pow9, types.AlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test finalclaim by v3
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator3Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.TestString)
}

func TestMultipleProphecies(t *testing.T) {
	//Test multiple prophecies running in parallel work fine as expected
	ctx, _, keeper, validatorAddresses, err := CreateTestKeepers(t, false, 0.7, []int64{3, 7})
	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test claim on first id with first validator
	status, err := keeper.ProcessClaim(ctx, types.TestID, validator1Pow3, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.PendingStatusText)

	//Test claim on second id with second validator
	status, err = keeper.ProcessClaim(ctx, types.AlternateTestID, validator2Pow7, types.AlternateTestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.AlternateTestString)

	//Test claim on first id with second validator
	status, err = keeper.ProcessClaim(ctx, types.TestID, validator2Pow7, types.TestString)
	require.NoError(t, err)
	require.Equal(t, status.StatusText, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, types.TestString)

	//Test claim on second id with first validator
	status, err = keeper.ProcessClaim(ctx, types.AlternateTestID, validator1Pow3, types.AlternateTestString)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Prophecy already finalized"))
}

func TestNonValidator(t *testing.T) {
	//TODO: anything from User that is not actually a validator fails
	err := errors.New("not yet implemented")
	require.NoError(t, err)
}
