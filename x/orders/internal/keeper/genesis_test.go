// +build unit

package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dfinance/dnode/helpers/perms"
	marketsClient "github.com/dfinance/dnode/x/markets/client"
	"github.com/dfinance/dnode/x/orders/internal/types"
)

func TestOrdersKeeper_Genesis_Init(t *testing.T) {
	input := NewTestInput(
		t,
		perms.Permissions{
			marketsClient.PermCreate,
			marketsClient.PermRead,
		},
	)

	keeper := input.keeper
	ctx := input.ctx
	ctx = ctx.WithBlockTime(time.Now().Add(time.Hour))
	cdc := input.cdc

	// default genesis
	{
		keeper.InitGenesis(ctx, cdc.MustMarshalJSON(types.DefaultGenesisState()))
		orders, err := keeper.GetList(input.ctx)
		require.Nil(t, err)
		require.Len(t, orders, 0)
	}

	// import and export
	{

		m, err := keeper.marketKeeper.Add(ctx, "btc", "dfi")
		require.Nil(t, err)

		exM, err := keeper.marketKeeper.GetExtended(ctx, m.ID)
		require.Nil(t, err)

		order := NewBtcDfiMockOrder(types.Ask)
		order.ID = keeper.nextID(ctx)
		order.Market = exM
		keeper.setID(ctx, order.ID)

		order2 := NewBtcDfiMockOrder(types.Bid)
		order2.ID = keeper.nextID(ctx)
		order2.Market = exM
		keeper.setID(ctx, order2.ID)

		lastId := keeper.getLastOrderID(ctx)

		state := types.GenesisState{
			Orders:      types.Orders{order, order2},
			LastOrderId: &lastId,
		}

		keeper.InitGenesis(ctx, cdc.MustMarshalJSON(state))
		orders, err := keeper.GetList(ctx)
		require.Nil(t, err)
		require.Len(t, orders, len(state.Orders))

		var exportedState types.GenesisState
		cdc.MustUnmarshalJSON(keeper.ExportGenesis(ctx), &exportedState)

		require.False(t, exportedState.IsEmpty())
		require.True(t, exportedState.Equal(state))
	}
}
