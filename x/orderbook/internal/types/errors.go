package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInternal = sdkErrors.Register(ModuleName, 100, "internal")
	// HistoryItem not found.
	ErrWrongHistoryItem = sdkErrors.Register(ModuleName, 101, "wrong marketID / blockHeight")
)
