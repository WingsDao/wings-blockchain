// Functions to work with genesis data of module.
package poa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dfinance/dnode/x/poa/types"
)

// Initialize genesis for this module.
func (poaKeeper Keeper) InitGenesis(ctx sdk.Context, genesisState types.GenesisState) {
	for _, val := range genesisState.PoAValidators {
		poaKeeper.AddValidator(ctx, val.Address, val.EthAddress)
	}
	poaKeeper.SetParams(ctx, genesisState.Parameters)
}

// Export genesis data for this module.
func (poaKeeper Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		Parameters:    poaKeeper.GetParams(ctx),
		PoAValidators: poaKeeper.GetValidators(ctx),
	}
}
