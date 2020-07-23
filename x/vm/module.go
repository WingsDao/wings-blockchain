// Module provides DVM integration with Move script/module execution/deployment.
// VM storage is used for writeSets (VM execution results) operations and used by other modules to preserve module dependant resources.
// DataSource server is a gRPC server for dnode-dvm async communication (DS context is updated during the BeginBlock).
// Module support DVM stdlib update proposal updates corresponding writeSets.
package vm

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dfinance/dnode/x/vm/client"
	"github.com/dfinance/dnode/x/vm/client/rest"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic app module basics object.
type AppModuleBasic struct{}

// Name gets module name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis gets default module genesis state.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates module genesis state.
func (AppModuleBasic) ValidateGenesis(data json.RawMessage) error {
	var state GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &state)

	return state.Validate()
}

// RegisterRESTRoutes registers module REST routes.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	rest.RegisterRoutes(ctx, r)
}

// GetTxCmd returns module root tx command.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetTxCmd(cdc)
}

// GetQueryCmd returns module root query command.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetQueryCmd(cdc)
}

// AppModule is a app module type.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates new AppModule object.
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name gets module name.
func (app AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers module invariants.
func (app AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns module messages route.
func (app AppModule) Route() string {
	return RouterKey
}

// NewHandler returns module messages handler.
func (app AppModule) NewHandler() sdk.Handler {
	return NewHandler(app.keeper)
}

// NewGovHandler retruns governance module proposals handler.
func (app AppModule) NewGovHandler() gov.Handler {
	return NewGovHandler(app.keeper)
}

// QuerierRoute returns module querier route.
func (app AppModule) QuerierRoute() string {
	return RouterKey
}

// NewQuerierHandler creates module querier.
func (app AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(app.keeper)
}

// InitGenesis inits module-genesis state.
func (app AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	app.keeper.InitGenesis(ctx, data)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports module genesis state.
func (app AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	genesisState := app.keeper.ExportGenesis(ctx)

	return ModuleCdc.MustMarshalJSON(genesisState)
}

// BeginBlock performs module actions at a block start.
func (app AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, app.keeper, req)
}

// EndBlock performs module actions at a block end.
func (app AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
