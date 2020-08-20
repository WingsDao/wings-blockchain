// +build unit

package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	dnTypes "github.com/dfinance/dnode/helpers/types"
	"github.com/dfinance/dnode/x/oracle/internal/types"
)

// Check Params for the set/get and specific methods for the asset, nominee, price parameters.
func TestOracleKeeper_Params(t *testing.T) {
	t.Parallel()
	input := NewTestInput(t)
	keeper := input.keeper
	ctx := input.ctx

	assetCode := dnTypes.AssetCode("btc_xfi")

	assetsMock := []types.Asset{
		types.Asset{AssetCode: assetCode, Oracles: []types.Oracle{}, Active: true},
	}

	nomineesMock := []string{"nominee1", "nominee2"}

	postPriceMock := types.PostPriceParams{ReceivedAtDiffInS: 100}

	paramsMock := types.Params{
		Assets:    assetsMock,
		Nominees:  nomineesMock,
		PostPrice: postPriceMock,
	}

	keeper.SetParams(ctx, paramsMock)

	// check GetAssetParams
	{
		assets := keeper.GetAssetParams(ctx)
		require.Equal(t, assets[0].AssetCode, assetCode)
		require.Equal(t, assets[0].Oracles, types.Oracles(types.Oracles(nil)))
	}

	// check GetNomineeParams
	{
		nominee := keeper.GetNomineeParams(ctx)
		require.Equal(t, nominee[0], nomineesMock[0])
		require.Equal(t, nominee[1], nomineesMock[1])
	}

	// check GetPostPriceParams
	{
		priceParam := keeper.GetPostPriceParams(ctx)
		require.Equal(t, priceParam, postPriceMock)
	}

	// check GetAllParams
	{
		params := keeper.GetParams(ctx)
		require.Equal(t, params.Assets[0].AssetCode, assetCode)
		require.Equal(t, params.Assets[0].Oracles, types.Oracles(types.Oracles(nil)))
		require.Equal(t, params.Nominees, nomineesMock)
		require.Equal(t, params.PostPrice, postPriceMock)
	}
}
