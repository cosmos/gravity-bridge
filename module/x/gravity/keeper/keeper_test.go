package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/app"
	"github.com/cosmos/gravity-bridge/module/app/params"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Gravity

	queryClient types.QueryClient
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *app.Gravity {
	gravityApp, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		gravityApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return gravityApp
}

// MakeEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	app.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	app.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func setup(withGenesis bool, invCheckPeriod uint) (*app.Gravity, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := MakeEncodingConfig()
	gravityApp := app.NewGravityApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{})
	if withGenesis {
		return gravityApp, app.NewDefaultGenesisState()
	}
	return gravityApp, app.GenesisState{}
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	gravityApp := Setup(checkTx)

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
