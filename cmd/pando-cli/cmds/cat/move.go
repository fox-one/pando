package cat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewMoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move <from id> <to id> <amount>",
		Short: "make a proposal to move supply between two collaterals",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from, to := args[0], args[1]
			amount := number.Decimal(args[2])

			values := []interface{}{
				core.ActionProposalMake,
				core.ActionCatMove,
				types.UUID(from),
				types.UUID(to),
				amount,
			}

			memo, err := actions.Build(cmd, values...)
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
