package edit

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
		Use:  "edit <collateral id> <key> <value>",
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			values := []interface{}{types.UUID(args[0])}
			for _, v := range args[1:] {
				values = append(values, v)
			}

			memo, err := actions.InitProposal(core.ActionCatEdit, values...)
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.CNB, number.One, memo)
		},
	}

	return cmd
}
