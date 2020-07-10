// Currencies module issues and withdraws currencies.
// Module is integrated with currencies storage module for CurrencyInfo and Balance resources.
// Issue is a multisig operation.
package currencies

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	codec "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dfinance/dnode/x/ccstorage"
	"github.com/dfinance/dnode/x/core/msmodule"
	"github.com/dfinance/dnode/x/currencies/client"
	"github.com/dfinance/dnode/x/currencies/client/rest"
	"github.com/dfinance/dnode/x/currencies/internal/keeper"
)

var (
	_ msmodule.AppMsModule  = AppModule{}
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
func (AppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

// ValidateGenesis validates module genesis state.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error { return nil }

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
	ccKeeper  keeper.Keeper
	ccsKeeper ccstorage.Keeper
}

// NewAppMsModule creates new AppMsModule object.
func NewAppMsModule(ccKeeper keeper.Keeper, ccsKeeper ccstorage.Keeper) msmodule.AppMsModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		ccKeeper:       ccKeeper,
		ccsKeeper:      ccsKeeper,
	}
}

// Name gets module name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers module invariants.
func (app AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, app.ccKeeper)
}

// Route returns module messages route.
func (app AppModule) Route() string {
	return RouterKey
}

// NewHandler returns module messages handler.
func (app AppModule) NewHandler() sdk.Handler {
	return NewHandler(app.ccKeeper)
}

// NewMsHandler returns module multisig messages handler.
func (app AppModule) NewMsHandler() msmodule.MsHandler {
	return NewMsHandler(app.ccKeeper)
}

// QuerierRoute returns module querier route.
func (app AppModule) QuerierRoute() string {
	return RouterKey
}

// NewQuerierHandler creates module querier.
func (app AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(app.ccKeeper)
}

// InitGenesis inits module-genesis state.
func (app AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports module genesis state.
func (app AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage { return nil }

// BeginBlock performs module actions at a block start.
func (app AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, app.ccsKeeper)
}

// EndBlock performs module actions at a block end.
func (app AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
