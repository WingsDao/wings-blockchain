package cli

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/dfinance/dnode/x/oracle/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		client.PostCommands(
			GetCmdPostPrice(cdc),
			getCmdAddOracle(cdc),
			getCmdSetOracles(cdc),
			getCmdSetAsset(cdc),
			getCmdAddAsset(cdc),
		)...,
	)

	return cmd
}

// GetCmdPostPrice cli command for posting prices.
func GetCmdPostPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "postprice [from_key_or_address] [assetCode] [price] [expiry]",
		Short: "post the latest price for a particular asset",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			rawPrice := args[2]
			price, ok := sdk.NewIntFromString(rawPrice)
			if !ok {
				return fmt.Errorf("%s argument %q: wrong value for price", "price", args[2])
			}
			expiryInt, ok := sdk.NewIntFromString(args[3])
			if !ok {
				fmt.Printf("invalid expiry - %s \n", args[2])
				return nil
			}
			expiry := tmtime.Canonical(time.Unix(expiryInt.Int64(), 0))
			msg := types.NewMsgPostPrice(cliCtx.GetFromAddress(), args[1], price, expiry)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdAddOracle(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "add-oracle [nominee_key] [denom] [oracle_address]",
		Example: "dncli oracle add-oracle wallet1a7280dyzp487r7wghr99f6r3h2h2z4gk4d740m eth_usdt wallet1a7260dyzp487r7wghr99f6r3h2h2z4gk4d740k",
		Short:   "Create a new oracle",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			oracleAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return fmt.Errorf("%s argument %q: %w", "oracle_address", args[2], err)
			}

			msg := types.NewMsgAddOracle(cliCtx.GetFromAddress(), args[1], oracleAddr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdSetOracles(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "set-oracles [nominee_key] [denom] [oracle_addresses]",
		Example: "dncli oracle set-oracles wallet1a7280dyzp487r7wghr99f6r3h2h2z4gk4d740m eth_usdt wallet10ff6y8gm2re6awfwz5dvesar8jq02tx7vcvuxn,wallet1a7260dyzp487r7wghr99f6r3h2h2z4gk4d740k",
		Short:   "Sets a list of oracles for a denom",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			oracles, err := types.ParseOracles(args[2])
			if err != nil {
				return fmt.Errorf("%s argument %q: %w", "oracle_addresses", args[2], err)
			}

			msg := types.NewMsgSetOracles(cliCtx.GetFromAddress(), args[1], oracles)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdAddAsset(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "add-asset [nominee_key] [denom] [oracles]",
		Example: "dncli oracle add-asset wallet1a7280dyzp487r7wghr99f6r3h2h2z4gk4d740m eth_usdt wallet1a7260dyzp487r7wghr99f6r3h2h2z4gk4d740k",
		Short:   "Create a new asset",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			denom := args[1]
			if len(denom) == 0 {
				return fmt.Errorf("%s argument %q: empty", "denom", args[1])
			}

			oracles, err := types.ParseOracles(args[2])
			if err != nil {
				return fmt.Errorf("%s argument %q: %w", "oracles", args[2], err)
			}
			if len(oracles) == 0 {
				return fmt.Errorf("%s argument %q: empty slice", "oracles", args[2])
			}

			token := types.NewAsset(denom, oracles, true)
			if err := token.ValidateBasic(); err != nil {
				return err
			}

			msg := types.NewMsgAddAsset(cliCtx.GetFromAddress(), denom, token)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdSetAsset(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "set-asset [nominee_key] [denom] [oracles]",
		Example: "dncli oracle set-asset wallet1a7280dyzp487r7wghr99f6r3h2h2z4gk4d740m eth_usdt wallet1a7260dyzp487r7wghr99f6r3h2h2z4gk4d740k",
		Short:   "Create a set asset",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			denom := args[1]
			if len(denom) == 0 {
				return fmt.Errorf("%s argument %q: empty", "denom", args[1])
			}

			oracles, err := types.ParseOracles(args[2])
			if err != nil {
				return fmt.Errorf("%s argument %q: %v", "oracles", args[2], err)
			}
			if len(oracles) == 0 {
				return fmt.Errorf("%s argument %q: empty slice", "oracles", args[2])
			}

			token := types.NewAsset(denom, oracles, true)
			if err := token.ValidateBasic(); err != nil {
				return err
			}

			msg := types.NewMsgSetAsset(cliCtx.GetFromAddress(), denom, token)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
