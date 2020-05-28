package orders

import (
	"github.com/dfinance/dnode/x/orders/internal/keeper"
	"github.com/dfinance/dnode/x/orders/internal/types"
)

type (
	Keeper         = keeper.Keeper
	Order          = types.Order
	Orders         = types.Orders
	OrderFill      = types.OrderFill
	OrderFills     = types.OrderFills
	Direction      = types.Direction
	MsgPostOrder   = types.MsgPostOrder
	MsgRevokeOrder = types.MsgRevokeOrder
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	BidDirection = types.Bid
	AskDirection = types.Ask
)

var (
	// variable aliases
	ModuleCdc = types.ModuleCdc
	// function aliases
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	// error aliases
	ErrWrongOwner = types.ErrWrongOwner
)