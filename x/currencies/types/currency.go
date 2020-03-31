// Currency type implementation.
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Currency struct {
	CurrencyId sdk.Int `json:"currencyId" swaggertype:"string"`
	Symbol     string  `json:"symbol" example:"dfi"`                      // Denom
	Supply     sdk.Int `json:"supply" swaggertype:"string" example:"100"` // Total amount
	Decimals   int8    `json:"decimals" example:"0"`
}

// New currency
func NewCurrency(symbol string, supply sdk.Int, decimals int8) Currency {
	return Currency{
		Symbol:   symbol,
		Supply:   supply,
		Decimals: decimals,
	}
}

func (c Currency) String() string {
	return fmt.Sprintf("Currency: \n"+
		"\tSymbol:      %s\n"+
		"\tSupply:      %s\n"+
		"\tDecimals:    %d\n",
		c.Symbol, c.Supply.String(), c.Decimals)
}
