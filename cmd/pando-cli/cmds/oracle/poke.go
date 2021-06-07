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

func NewPokeCmd() *cobra.Command {
	var ts int64

	cmd := &cobra.Command{
		Use:   "poke <asset_id> <price>",
		Short: "make a proposal to execute poke action on an oracle asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, price := args[0], args[1]

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionOraclePoke,
				types.UUID(id),
				types.Decimal(price),
				ts,
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	cmd.Flags().Int64Var(&ts, "ts", time.Now().Unix(), "timestamp")
	return cmd
}
