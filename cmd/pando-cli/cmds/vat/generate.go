package vat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generate <collateral id> <debt>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			vatID := args[0]
			debt := number.Decimal(args[1])

			memo, err := actions.Build(cmd, core.ActionVatGenerate, types.UUID(vatID), debt)
			if err != nil {
				return err
			}

			return pay.Request(ctx, pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
