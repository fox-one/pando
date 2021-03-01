package feed

import (
	"time"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "feed",
		Args: cobra.ExactValidArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, price := args[0], args[1]

			memo, err := actions.InitProposal(
				core.ActionOracleFeed,
				types.UUID(id),
				types.Decimal(price),
				time.Now().Unix(),
			)
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.CNB, number.One, memo)
		},
	}

	return cmd
}
