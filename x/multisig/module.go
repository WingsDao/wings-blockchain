// Module provides multi signature messages routing and handling with submit, confirm and revoke functions.
// Once message call is submitted it should be confirmed by 2/3 of POA validators, that changes call status to approved.
// An approved call executes message handler (via multisig router).
// Call has a TTL level in blocks, which changes its state to rejected if call isn't approved withing defined period.
package multisig

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dfinance/dnode/x/multisig/client"
	"github.com/dfinance/dnode/x/multisig/client/rest"
	"github.com/dfinance/dnode/x/multisig/internal/keeper"
	"github.com/dfinance/dnode/x/poa"
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
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	state := GenesisState{}
	ModuleCdc.MustUnmarshalJSON(bz, &state)

	return state.Validate(-1)
}

// RegisterRESTRoutes registers module REST routes.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	rest.RegisterRoutes(ctx, r)
}

// GetTxCmd returns module root tx command.
func (AppModuleBasic) GetTxCmd(cdc *amino.Codec) *cobra.Command {
	return client.GetTxCmd(cdc)
}

// GetQueryCmd returns module root query command.
func (AppModuleBasic) GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	return client.GetQueryCmd(cdc)
}

// AppModule is a app module type.
type AppModule struct {
	AppModuleBasic
	msKeeper  keeper.Keeper
	poaKeeper poa.Keeper
}

// NewAppModule creates new AppModule object.
func NewAppModule(msKeeper keeper.Keeper, poaKeeper poa.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		msKeeper:       msKeeper,
		poaKeeper:      poaKeeper,
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
	return NewHandler(app.msKeeper, app.poaKeeper)
}

// QuerierRoute returns module querier route.
func (app AppModule) QuerierRoute() string {
	return RouterKey
}

// NewQuerierHandler creates module querier.
func (app AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(app.msKeeper)
}

// InitGenesis inits module-genesis state.
func (app AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	app.msKeeper.InitGenesis(ctx, data)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports module genesis state.
func (app AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return app.msKeeper.ExportGenesis(ctx)
}

// BeginBlock performs module actions at a block start.
func (app AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock performs module actions at a block end.
func (app AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, app.msKeeper, app.poaKeeper)
}
