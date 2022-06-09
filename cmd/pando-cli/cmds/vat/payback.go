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

func NewPaybackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payback",
		Short: "execute Payback action on a vault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			vatID := args[0]
			vat, err := call.RPC().FindVault(ctx, &api.Req_FindVault{Id: vatID})
			if err != nil {
				return err
			}

			cat, err := call.RPC().FindCollateral(ctx, &api.Req_FindCollateral{Id: vat.CollateralId})
			if err != nil {
				return err
			}

			debt := number.Decimal(args[1])
			memo, err := actions.Build(cmd, core.ActionVatPayback, types.UUID(vatID))
			if err != nil {
				return err
			}

			return pay.Request(ctx, cat.Dai, debt, memo)
		},
	}

	return cmd
}
