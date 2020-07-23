// +build unit

package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dfinance/dnode/helpers/types"
)

func TestOrders_Genesis_Valid(t *testing.T) {
	//validateGenesis ok
	{
		order := NewMockOrder()
		order.ID = types.NewIDFromUint64(100)
		order2 := NewMockOrder()
		order2.Market.QuoteCurrency.Denom = "btc"

		order3 := NewMockOrder()
		orderT := &order3
		order4 := *orderT

		state := GenesisState{Orders: Orders{order, order2, order3}}
		require.NoError(t, state.Validate())
		require.False(t, state.IsEmpty())

		require.False(t, GenesisState{Orders: Orders{order2}}.Equal(GenesisState{Orders: Orders{order3}}))
		require.True(t, GenesisState{Orders: Orders{order3}}.Equal(GenesisState{Orders: Orders{order4}}))
	}

	// wrong id
	{
		order := NewMockOrder()
		order.ID, _ = types.NewIDFromString("")
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "id")
		require.Contains(t, err, "nil")
	}

	//validateGenesis wrong owner
	{
		order := NewMockOrder()
		order.Owner = sdk.AccAddress{}
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "owner")
		require.Contains(t, err, "empty")
	}

	// wrong owner
	{
		order := NewMockOrder()
		order.Owner = sdk.AccAddress{}
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "owner")
		require.Contains(t, err, "empty")
	}

	// wrong direction
	{
		order := NewMockOrder()
		order.Direction = "wrong"
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "direction")
		require.Contains(t, err, "invalid")
	}

	// wrong price
	{
		order := NewMockOrder()
		order.Price = sdk.NewUint(0)
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "price")
		require.Contains(t, err, "zero")
	}

	// wrong quantity
	{
		order := NewMockOrder()
		order.Quantity = sdk.NewUint(0)
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "quantity")
		require.Contains(t, err, "zero")
	}

	// wrong dates
	{
		order := NewMockOrder()
		order.CreatedAt = time.Unix(2, 0)
		order.UpdatedAt = time.Unix(1, 0)
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "wrong create and update dates")
	}

	// wrong market id
	{
		order := NewMockOrder()
		order.Market.ID, _ = types.NewIDFromString("")
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "market")
		require.Contains(t, err, "id")
		require.Contains(t, err, "nil")
	}

	// wrong market BaseCurrency Denom
	{
		order := NewMockOrder()
		order.Market.BaseCurrency.Denom = ""
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "market")
		require.Contains(t, err, "denom")
		require.Contains(t, err, "base")
		require.Contains(t, err, "empty")
	}

	// wrong market QuoteCurrency Denom
	{
		order := NewMockOrder()
		order.Market.QuoteCurrency.Denom = ""
		err := GenesisState{Orders: Orders{order}}.Validate().Error()
		require.Contains(t, err, "market")
		require.Contains(t, err, "denom")
		require.Contains(t, err, "quote")
		require.Contains(t, err, "empty")
	}
}
