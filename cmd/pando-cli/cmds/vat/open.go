package vat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open <collateral id> <deposit> <generate>",
		Short: "open a new vault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			catID := args[0]
			cat, err := call.RPC().FindCollateral(ctx, &api.Req_FindCollateral{Id: catID})
			if err != nil {
				return err
			}

			dink := number.Decimal(args[1])
			debt := number.Decimal(args[2])
			memo, err := actions.Build(cmd, core.ActionVatOpen, types.UUID(catID), debt)
			if err != nil {
				return err
			}

			return pay.Request(ctx, cat.Gem, dink, memo)
		},
	}

	return cmd
}
