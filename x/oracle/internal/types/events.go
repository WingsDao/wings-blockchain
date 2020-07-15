package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeAddAsset = ModuleName + ".add_asset"
	EventTypePrice    = ModuleName + ".price"
	//
	AttributeAssetCode  = "asset_code"
	AttributePrice      = "price"
	AttributeReceivedAt = "received_at"
)

// NewAssetAddedEvent creates an Event on asset creation.
func NewAssetAddedEvent(asset Asset) sdk.Event {
	return sdk.NewEvent(EventTypeAddAsset,
		sdk.NewAttribute(AttributeAssetCode, asset.AssetCode.String()),
	)
}

// NewPriceEvent creates an Event on price update.
func NewPriceEvent(price CurrentPrice) sdk.Event {
	return sdk.NewEvent(EventTypePrice,
		sdk.NewAttribute(AttributeAssetCode, price.AssetCode.String()),
		sdk.NewAttribute(AttributePrice, price.Price.String()),
		sdk.NewAttribute(AttributeReceivedAt, strconv.FormatInt(price.ReceivedAt.Unix(), 10)),
	)
}
