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

func NewFoldCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "fold <collateral id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if strings.EqualFold(id, "all") {
				id = uuid.Zero.String()
			}

			memo, err := actions.Build(cmd, core.ActionCatFold, types.UUID(id))
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
