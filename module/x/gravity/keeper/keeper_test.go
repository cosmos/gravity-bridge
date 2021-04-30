package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/app"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Gravity

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	gravityApp := app.Setup(checkTx)

	suite.ctx = gravityApp.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
	suite.app = gravityApp

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, gravityApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, gravityApp.GravityKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestBridgeIDCRUD() {
	id := []byte("id")
	suite.app.GravityKeeper.SetBridgeID(suite.ctx, id)
	returnedID := suite.app.GravityKeeper.GetBridgeID(suite.ctx)
	suite.Require().Equal(id, returnedID)
}
