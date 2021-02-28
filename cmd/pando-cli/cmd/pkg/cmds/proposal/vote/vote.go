package vote

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "vote",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			memo, err := actions.Member(core.ActionProposalVote, types.UUID(id))

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
