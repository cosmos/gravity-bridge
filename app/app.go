package app

import (
	"os"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/peggy/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	appName = "ethereum-bridge"
)

var (
	// DefaultCLIHome default home directories for ebcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.ebcli")

	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.ebd")

	// ModuleBasics the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)

// EthereumBridgeApp defines the Ethereum-Cosmos peg-zone application
type EthereumBridgeApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// SDK keepers
	// TODO: add governance keeper
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	StakingKeeper staking.Keeper
	SupplyKeeper  supply.Keeper
	CrisisKeeper  crisis.Keeper
	ParamsKeeper  params.Keeper

	// EthBridge keepers
	OracleKeeper oracle.Keeper

	// the module manager
	mm *module.Manager
}

// NewEthereumBridgeApp is a constructor function for EthereumBridgeApp
func NewEthereumBridgeApp(logger log.Logger, db dbm.DB, invCheckPeriod uint,
	baseAppOptions ...func(*bam.BaseApp)) *EthereumBridgeApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, oracle.StoreKey, params.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &EthereumBridgeApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    keys,
		tkeys:   tkeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.BankKeeper = bank.NewBaseKeeper(app.AccountKeeper, bankSubspace, bank.DefaultCodespace, app.ModuleAccountAddrs())
	app.SupplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms)
	app.StakingKeeper = staking.NewKeeper(app.cdc, keys[staking.StoreKey], tkeys[staking.TStoreKey],
		app.SupplyKeeper, stakingSubspace, staking.DefaultCodespace)
	app.CrisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName)
	app.OracleKeeper = oracle.NewKeeper(app.cdc, keys[oracle.StoreKey], app.StakingKeeper, oracle.DefaultCodespace, oracle.DefaultConsensusNeeded)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.AccountKeeper),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		// TODO: Oracle/Ethrelayer module
	)

	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		genaccounts.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
	)

	// NOTE: this doesn't do anything at the moment because peggy does not implement
	// the simulator.
	// TODO: add simulator support
	app.sm = module.NewSimulationManager(app.mm.Modules)
	app.sm.RegisterStoreDecoders()

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// BeginBlocker application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height
func (app *ethereumBridgeApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// Codec returns simapp's codec
func (app *SimApp) Codec() *codec.Codec {
	return app.cdc
}

// GetKey returns the KVStoreKey for the provided store key
func (app *SimApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *SimApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
