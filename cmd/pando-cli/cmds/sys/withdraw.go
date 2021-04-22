package sys

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "withdraw <asset> <amount> <opponent>",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionSysWithdraw,
				types.UUID(args[0]),
				types.Decimal(args[1]),
				types.UUID(args[2]))

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
