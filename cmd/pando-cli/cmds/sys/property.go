package sys

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewPropertyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "property <key> <value>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionSysProperty,
				key,
				value,
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
