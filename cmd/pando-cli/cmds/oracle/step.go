package oracle

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func NewStepCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "step",
		Args: cobra.ExactValidArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, ts := args[0], args[1]

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionOracleFeed,
				types.UUID(id),
				cast.ToInt64(ts),
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
