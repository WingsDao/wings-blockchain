package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dfinance/dnode/x/ccstorage"
	"github.com/dfinance/dnode/x/currencies/internal/types"
)

// CreateCurrency redirects CreateCurrency request to the currencies storage.
func (k Keeper) CreateCurrency(ctx sdk.Context, params ccstorage.CurrencyParams) error {
	k.modulePerms.AutoCheck(types.PermCreate)

	return k.ccsKeeper.CreateCurrency(ctx, params)
}
