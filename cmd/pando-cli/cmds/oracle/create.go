package oracle

import (
	"time"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var (
		price     string
		hop       int64
		threshold int64
		ts        int64
	)

	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionOracleCreate,
				types.UUID(id),
				types.Decimal(price),
				hop,
				threshold,
				ts,
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	cmd.Flags().StringVar(&price, "price", "0", "oracle price")
	cmd.Flags().Int64Var(&hop, "hop", 60*60, "poke delay")
	cmd.Flags().Int64Var(&threshold, "threshold", 0, "number of signatures required")
	cmd.Flags().Int64Var(&ts, "ts", time.Now().Unix(), "timestamp, default is time.Now()")

	return cmd
}
