package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dnTypes "github.com/dfinance/dnode/helpers/types"
)

const (
	PriceBytesLimit = 8
)

// CurrentPrice struct that contains the metadata of a current price for a particular asset in the oracle module.
type CurrentPrice struct {
	// Asset code
	AssetCode dnTypes.AssetCode `json:"asset_code" yaml:"asset_code" example:"dfi"`
	// Price in integer format
	Price sdk.Int `json:"price" yaml:"price" swaggertype:"string" example:"1000"`
	// Timestamp Price createdAt
	ReceivedAt time.Time `json:"received_at" yaml:"received_at" format:"RFC 3339" example:"2020-03-27T13:45:15.293426Z"`
}

// String implement fmt.Stringer for the CurrentPrice type.
func (cp CurrentPrice) String() string {
	return strings.TrimSpace(fmt.Sprintf("AssetCode: %s\nPrice: %s\nReceivedAt: %s", cp.AssetCode, cp.Price, cp.ReceivedAt))
}

// PostedPrice struct represented a price for an asset posted by a specific oracle.
type PostedPrice struct {
	AssetCode     dnTypes.AssetCode `json:"asset_code" yaml:"asset_code" example:"dfi"`                                                                        // Denom
	OracleAddress sdk.AccAddress    `json:"oracle_address" yaml:"oracle_address" swaggertype:"string" example:"wallet13jyjuz3kkdvqw8u4qfkwd94emdl3vx394kn07h"` // Price source
	Price         sdk.Int           `json:"price" yaml:"price" swaggertype:"string" example:"1000"`
	ReceivedAt    time.Time         `json:"received_at" yaml:"received_at" format:"RFC 3339" example:"2020-03-27T13:45:15.293426Z"` // Timestamp Price createdAt
}

// String implement fmt.Stringer for the PostedPrice type.
func (pp PostedPrice) String() string {
	return strings.TrimSpace(
		fmt.Sprintf(
			"AssetCode: %s\nOracleAddress: %s\nPrice: %s\nReceivedAt: %s",
			pp.AssetCode,
			pp.OracleAddress,
			pp.Price,
			pp.ReceivedAt,
		),
	)
}

// PendingPriceAsset struct that contains the info about the asset which price is still to be determined.
type PendingPriceAsset struct {
	AssetCode string `json:"asset_code"`
}

// String implement fmt.Stringer for the PendingPriceAsset type.
func (a PendingPriceAsset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s`, a.AssetCode))
}
