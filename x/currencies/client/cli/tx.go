// Transaction commands for currencies CLI implementation.
package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	txBldrCtx "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"

	"github.com/dfinance/dnode/x/currencies/msgs"
)

// Destroy currency.
func PostDestroyCurrency(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy-currency [chainID] [symbol] [amount] [recipient]",
		Short: "destroy issued currency",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := txBldrCtx.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			accGetter := txBldrCtx.NewAccountRetriever(cliCtx)

			if err := accGetter.EnsureExists(cliCtx.FromAddress); err != nil {
				return fmt.Errorf("fromAddress: %w", err)
			}

			amount, isOk := sdk.NewIntFromString(args[2])
			if !isOk {
				return fmt.Errorf("%s argument %q is not a number, can't parse int", "amount", args[2])
			}

			msg := msgs.NewMsgDestroyCurrency(args[0], args[1], amount, cliCtx.GetFromAddress(), args[3])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			cliCtx.WithOutput(os.Stdout)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
