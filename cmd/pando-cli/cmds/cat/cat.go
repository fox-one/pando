package cat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat/create"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat/edit"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat/fold"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat/list"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat/supply"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cat",
	}

	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(supply.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(edit.NewCmd())
	cmd.AddCommand(fold.NewCmd())

	return cmd
}
