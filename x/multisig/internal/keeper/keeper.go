// Multi signature module keeper stores call objects, calls queue with submitting, confirming and revoking.
package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dfinance/dnode/helpers/perms"
	"github.com/dfinance/dnode/x/core/msmodule"
	"github.com/dfinance/dnode/x/multisig/internal/types"
	"github.com/dfinance/dnode/x/poa"
)

// Module keeper object.
type Keeper struct {
	cdc         *codec.Codec
	storeKey    sdk.StoreKey
	paramStore  params.Subspace
	router      msmodule.MsRouter
	poaKeeper   poa.Keeper
	modulePerms perms.ModulePermissions
}

// GetLogger gets logger with keeper context.
func (k Keeper) GetLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetRouteHandler returns multi signature router handler for specific path.
func (k Keeper) GetRouteHandler(route string) msmodule.MsHandler {
	k.modulePerms.AutoCheck(types.PermReader)

	return k.router.GetRoute(route)
}

// CheckAddressIsPoaValidator checks if {address} is a registered POA validator.
func (k Keeper) CheckAddressIsPoaValidator(ctx sdk.Context, address sdk.AccAddress) error {
	k.modulePerms.AutoCheck(types.PermPoaReader)

	if !k.poaKeeper.HasValidator(ctx, address) {
		return sdkErrors.Wrap(types.ErrPoaNotValidator, address.String())
	}

	return nil
}

// GetPoaMinConfirmationsCount return POA module minimum confirmations count to approve call.
func (k Keeper) GetPoaMinConfirmationsCount(ctx sdk.Context) uint16 {
	k.modulePerms.AutoCheck(types.PermPoaReader)

	return k.poaKeeper.GetEnoughConfirmations(ctx)
}

// Create new currency keeper.
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramStore params.Subspace,
	router msmodule.MsRouter,
	poaKeeper poa.Keeper,
	permsRequesters ...perms.RequestModulePermissions,
) Keeper {
	k := Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramStore: paramStore.WithKeyTable(types.ParamKeyTable()),
		router:     router,
		poaKeeper:  poaKeeper,
		modulePerms: types.NewModulePerms(),
	}
	for _, requester := range permsRequesters {
		k.modulePerms.AutoAddRequester(requester)
	}

	return k
}
