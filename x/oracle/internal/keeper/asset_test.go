// +build unit

package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	dnTypes "github.com/dfinance/dnode/helpers/types"
	"github.com/dfinance/dnode/x/oracle/internal/types"
)

// Check SetAsset method with various sets of arguments.
func TestOracleKeeper_SetAsset(t *testing.T) {
	t.Parallel()

	input := NewTestInput(t)
	keeper := input.keeper
	ctx := input.ctx

	asset := types.NewAsset(input.stdAssetCode, []types.Oracle{}, true)

	// set asset
	{
		err := keeper.SetAsset(ctx, input.stdNominee, asset)
		require.Nil(t, err)
	}

	// set asset with wrong nominee
	{
		err := keeper.SetAsset(ctx, "wrongNominee", asset)
		require.Error(t, err)
	}

	// wrong asset code, doesn't exist
	{
		assetT := &asset
		asset2 := *assetT
		asset2.AssetCode = dnTypes.AssetCode("btc_eth")
		err := keeper.SetAsset(ctx, input.stdNominee, asset2)
		require.Error(t, err)
	}
}

// Check AddAsset method with various sets of arguments.
func TestOracleKeeper_AddAsset(t *testing.T) {
	t.Parallel()

	input := NewTestInput(t)
	keeper := input.keeper
	ctx := input.ctx

	asset := types.NewAsset("btc_usdt", []types.Oracle{}, true)

	// add asset
	{
		err := keeper.AddAsset(ctx, input.stdNominee, asset)
		require.Nil(t, err)

		_, ok := keeper.GetAsset(ctx, asset.AssetCode)
		require.True(t, ok)
	}

	// add asset with wrong nominee
	{
		err := keeper.AddAsset(ctx, "wrongNominee", asset)
		require.Error(t, err)
	}

	// double add
	{
		asset2 := types.NewAsset(input.stdAssetCode, []types.Oracle{}, true)
		err := keeper.AddAsset(ctx, input.stdNominee, asset2)
		require.Error(t, err)
	}
}

// Check GetAsset method returning values.
func TestOracleKeeper_GetAsset(t *testing.T) {
	t.Parallel()

	input := NewTestInput(t)
	keeper := input.keeper
	ctx := input.ctx

	asset := types.NewAsset(input.stdAssetCode, []types.Oracle{}, true)

	// get asset
	{
		err := keeper.SetAsset(ctx, input.stdNominee, asset)
		require.Nil(t, err)

		a, ok := keeper.GetAsset(ctx, input.stdAssetCode)
		require.Equal(t, true, ok)
		require.Equal(t, a.AssetCode, input.stdAssetCode)

		_, ok = keeper.GetAsset(ctx, "btc_eth")
		require.Equal(t, false, ok)
	}
}
