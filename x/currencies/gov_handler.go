package currencies

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewGovHandler creates proposal type handler for Gov module.
func NewGovHandler(keeper Keeper) gov.Handler {
	return func(ctx sdk.Context, c govTypes.Content) error {
		if c.ProposalRoute() != GovRouterKey {
			return fmt.Errorf("invalid proposal route %q for module %q", c.ProposalRoute(), ModuleName)
		}

		switch p := c.(type) {
		case AddCurrencyProposal:
			return handleAddCurrencyProposal(ctx, keeper, p)
		default:
			return fmt.Errorf("unsupported proposal content type %q for module %q", c.ProposalType(), ModuleName)
		}
	}
}

// handleAddCurrencyProposal handles currency creation proposal.
func handleAddCurrencyProposal(ctx sdk.Context, keeper Keeper, p AddCurrencyProposal) error {
	logger := keeper.GetLogger(ctx)

	err := keeper.CreateCurrency(ctx, p.Denom, p.GetCurrencyParams())
	if err != nil {
		return sdkErrors.Wrapf(ErrGovInvalidProposal, "creating currency: %v", err)
	}

	logger.Info(fmt.Sprintf("proposal executed:\n%s", p.String()))

	return nil
}