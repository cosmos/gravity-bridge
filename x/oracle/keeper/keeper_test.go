package keeper

import (
	"strings"
	"testing"

	"github.com/trinhtan/peggy/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestCreateGetProphecy(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	validator1Pow3 := validatorAddresses[0]

	//Test normal Creation
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test bad Creation with blank id
	oracleClaim = types.NewClaim("", validator1Pow3, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)

	//Test bad Creation with blank claim
	oracleClaim = types.NewClaim(TestID, validator1Pow3, "")
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)

	//Test retrieval
	prophecy, found := keeper.GetProphecy(ctx, TestID)
	require.True(t, found)
	require.Equal(t, prophecy.ID, TestID)
	require.Equal(t, prophecy.Status.Text, types.PendingStatusText)
	require.Equal(t, prophecy.ClaimValidators[TestString][0], validator1Pow3)
	require.Equal(t, prophecy.ValidatorClaims[validator1Pow3.String()], TestString)
}

func TestBadConsensusForOracle(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_, _, _, _, _, _ = CreateTestKeepers(t, 0, []int64{10}, "")
	_, _, _, _, _, _ = CreateTestKeepers(t, 1.2, []int64{10}, "")
}

func TestBadMsgs(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.6, []int64{3, 3}, "")

	validator1Pow3 := validatorAddresses[0]

	//Test empty claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3, "")
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.Equal(t, status.FinalClaim, "")
	require.True(t, strings.Contains(err.Error(), "claim cannot be empty string"))

	//Test normal Creation
	oracleClaim = types.NewClaim(TestID, validator1Pow3, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test duplicate message
	oracleClaim = types.NewClaim(TestID, validator1Pow3, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	//Test second but non duplicate message
	oracleClaim = types.NewClaim(TestID, validator1Pow3, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestSuccessfulProphecy(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test second claim completes and finalizes to success
	oracleClaim = types.NewClaim(TestID, validator2Pow3, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, TestString)

	//Test third claim not possible
	oracleClaim = types.NewClaim(TestID, validator3Pow4, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
}

func TestSuccessfulProphecyWithDisagreement(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test second disagreeing claim processed fine
	oracleClaim = types.NewClaim(TestID, validator2Pow3, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test third claim agrees and finalizes to success
	oracleClaim = types.NewClaim(TestID, validator3Pow4, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, TestString)
}

func TestFailedProphecy(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.6, []int64{3, 3, 4}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	validator3Pow4 := validatorAddresses[2]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test second disagreeing claim processed fine
	oracleClaim = types.NewClaim(TestID, validator2Pow3, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)
	require.Equal(t, status.FinalClaim, "")

	//Test third disagreeing claim processed fine and prophecy fails
	oracleClaim = types.NewClaim(TestID, validator3Pow4, AnotherAlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.FailedStatusText)
	require.Equal(t, status.FinalClaim, "")
}

func TestPowerOverrule(t *testing.T) {
	//Testing with 2 validators but one has high enough power to overrule
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test first claim
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test second disagreeing claim processed fine and finalized to its bytes
	oracleClaim = types.NewClaim(TestID, validator2Pow7, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, AlternateTestString)
}
func TestPowerAternate(t *testing.T) {
	//Test alternate power setup with validators of 5/4/3/9 and total power 22 and 12/21 required
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.571, []int64{5, 4, 3, 9}, "")

	validator1Pow5 := validatorAddresses[0]
	validator2Pow4 := validatorAddresses[1]
	validator3Pow3 := validatorAddresses[2]
	validator4Pow9 := validatorAddresses[3]

	//Test claim by v1
	oracleClaim := types.NewClaim(TestID, validator1Pow5, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test claim by v2
	oracleClaim = types.NewClaim(TestID, validator2Pow4, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test alternate claim by v4
	oracleClaim = types.NewClaim(TestID, validator4Pow9, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test finalclaim by v3
	oracleClaim = types.NewClaim(TestID, validator3Pow3, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, TestString)
}

func TestMultipleProphecies(t *testing.T) {
	//Test multiple prophecies running in parallel work fine as expected
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow7 := validatorAddresses[1]

	//Test claim on first id with first validator
	oracleClaim := types.NewClaim(TestID, validator1Pow3, TestString)
	status, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.PendingStatusText)

	//Test claim on second id with second validator
	oracleClaim = types.NewClaim(AlternateTestID, validator2Pow7, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, AlternateTestString)

	//Test claim on first id with second validator
	oracleClaim = types.NewClaim(TestID, validator2Pow7, TestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, types.SuccessStatusText)
	require.Equal(t, status.FinalClaim, TestString)

	//Test claim on second id with first validator
	oracleClaim = types.NewClaim(AlternateTestID, validator1Pow3, AlternateTestString)
	status, err = keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
}

func TestNonValidator(t *testing.T) {
	//Test multiple prophecies running in parallel work fine as expected
	ctx, keeper, _, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 7}, "")

	_, testValidatorAddresses := CreateTestAddrs(10)
	inActiveValidatorAddress := testValidatorAddresses[9]

	//Test claim on first id with first validator
	oracleClaim := types.NewClaim(TestID, inActiveValidatorAddress, TestString)
	_, err := keeper.ProcessClaim(ctx, oracleClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "claim must be made by actively bonded validator"))
}
