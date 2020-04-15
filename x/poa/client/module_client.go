package client

import (
	sdkClient "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/dfinance/dnode/x/poa/client/cli"
	"github.com/dfinance/dnode/x/poa/types"
)

// Return query commands for PoA module.
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "PoA commands for the validators module",
	}

	queryCmd.AddCommand(sdkClient.GetCommands(
		cli.GetValidator(types.ModuleName, cdc),
		cli.GetValidators(types.ModuleName, cdc),
		cli.GetMinMax(types.ModuleName, cdc),
	)...)

	return queryCmd
}

// Returns transactions commands for this module.
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "PoA transactions subcommands",
	}

	txCmd.AddCommand(sdkClient.PostCommands(
		cli.PostMsAddValidator(cdc),
		cli.PostMsRemoveValidator(cdc),
		cli.PostMsReplaceValidator(cdc),
	)...)

	return txCmd
}
