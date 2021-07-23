package cat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewGainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gain <collateral id> <amount> <receipt>",
		Short: "execute Gain action on a collateral",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			amount := number.Decimal(args[1])
			receipt := args[2]

			memo, err := actions.Build(
				cmd, core.ActionProposalMake,
				core.ActionCatGain,
				types.UUID(id),
				amount,
				types.UUID(receipt),
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
