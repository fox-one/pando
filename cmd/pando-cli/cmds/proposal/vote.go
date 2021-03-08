package proposal

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "vote",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			memo, err := actions.Build(cmd, core.ActionProposalVote, types.UUID(id))
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
