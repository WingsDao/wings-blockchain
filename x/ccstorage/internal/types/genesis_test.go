// +build unit

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test genesis params validation.
func TestCCS_GenesisParams_Validate(t *testing.T) {
	t.Parallel()

	// ok
	{
		param := CurrencyParams{"xfi", 0}
		require.NoError(t, param.Validate())
	}

	// fail: invalid denom
	{
		param1 := CurrencyParams{"xfi1", 0}
		require.Error(t, param1.Validate())
	}
}

// Test genesis validation.
func TestCCS_Genesis_Validate(t *testing.T) {
	t.Parallel()

	state := GenesisState{}

	// ok: empty
	{
		require.NoError(t, state.Validate())
	}

	// ok: new 1
	{
		state.CurrenciesParams = append(state.CurrenciesParams, CurrencyParams{
			Denom:    "xfi",
			Decimals: 0,
		})
		require.NoError(t, state.Validate())
	}

	// ok: new 2
	{
		state.CurrenciesParams = append(state.CurrenciesParams, CurrencyParams{
			Denom:    "btc",
			Decimals: 8,
		})
		require.NoError(t, state.Validate())
	}

	// fail: duplicate
	{
		state.CurrenciesParams = append(state.CurrenciesParams, CurrencyParams{
			Denom:    "btc",
			Decimals: 4,
		})
		require.Error(t, state.Validate())
	}
}
