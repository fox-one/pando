package cat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <gem> <dai> <name>",
		Short: "make a proposal to create new collateral, admin required",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionCatCreate,
				types.UUID(args[0]),
				types.UUID(args[1]),
				args[2],
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
