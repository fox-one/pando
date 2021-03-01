package create

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create <gem> <dai> <name>",
		Args: cobra.ExactValidArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			memo, err := actions.MakeProposal(
				core.ActionCatCreate,
				types.UUID(args[0]),
				types.UUID(args[1]),
				args[2],
			)

			if err != nil {
				return err
			}

			return pay.Request(
				cmd.Context(),
				pay.CNB,
				number.One,
				memo,
			)
		},
	}

	return cmd
}
