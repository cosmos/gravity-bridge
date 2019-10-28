package oracle

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the oracle module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the oracle module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the oracle module's types for the given codec.
func (AppModuleBasic) RegisterCodec(_ *codec.Codec) {}

// DefaultGenesis returns default genesis state as raw bytes for the oracle
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the oracle module.
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the oracle module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

// GetTxCmd returns the root tx command for the oracle module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command for the oracle module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command {
	return nil
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the oracle module.
type AppModuleSimulation struct{}

// AppModule implements an application module for the oracle module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
	}
}

// Name returns the oracle module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the oracle module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route returns the message routing key for the oracle module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the oracle module.
func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

// QuerierRoute returns the oracle module's querier route name.
func (AppModule) QuerierRoute() string {
	return ""
}

// NewQuerierHandler returns the oracle module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// InitGenesis performs genesis initialization for the oracle module. It returns
// no validator updates.
func (am AppModule) InitGenesis(_ sdk.Context, _ json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the oracle
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker for the oracle module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the oracle module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
