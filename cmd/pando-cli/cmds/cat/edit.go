package cat

import (
	"strings"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/spf13/cobra"
)

func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit <collateral id> <key> <value>",
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if strings.EqualFold(id, "all") {
				id = uuid.Zero.String()
			}

			values := []interface{}{core.ActionProposalMake, core.ActionCatEdit, types.UUID(id)}
			for _, v := range args[1:] {
				values = append(values, v)
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
