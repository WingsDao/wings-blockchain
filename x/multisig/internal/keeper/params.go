package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dfinance/dnode/x/multisig/internal/types"
)

// GetIntervalToExecute returns intervalToExecute param.
func (k Keeper) GetIntervalToExecute(ctx sdk.Context) (res int64) {
	k.modulePerms.AutoCheck(types.PermRead)

	k.paramStore.Get(ctx, types.ParamStoreKeyIntervalToExecute, &res)

	return
}

// SetIntervalToExecute updates intervalToExecute param.
func (k Keeper) SetIntervalToExecute(ctx sdk.Context, value int64) {
	k.modulePerms.AutoCheck(types.PermWrite)

	k.paramStore.Set(ctx, types.ParamStoreKeyIntervalToExecute, value)
}
