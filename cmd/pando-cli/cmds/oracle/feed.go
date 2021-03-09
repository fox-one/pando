package oracle

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewFeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "feed <asset_id> <price>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, price := args[0], args[1]

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionOracleFeed,
				types.UUID(id),
				types.Decimal(price),
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
