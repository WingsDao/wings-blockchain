package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	orderTypes "github.com/dfinance/dnode/x/order"
	"github.com/dfinance/dnode/x/orderbook/internal/keeper"
	"github.com/dfinance/dnode/x/orderbook/internal/types"
)

// EndBlocker iterates over Order module orders, processes them and returns back to the Order module.
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	iterator := k.GetOrderIterator(ctx)
	defer iterator.Close()

	matcherPool := keeper.NewMatcherPool(k.GetLogger(ctx))
	for ; iterator.Valid(); iterator.Next() {
		order := orderTypes.Order{}
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &order)

		if err := matcherPool.AddOrder(order); err != nil {
			panic(err)
		}
	}

	for _, result := range matcherPool.Process() {
		k.ProcessOrderFills(ctx, result.OrderFills)
	}

	return []abci.ValidatorUpdate{}
}