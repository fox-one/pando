package config

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/config/get"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/config/set"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/config/show"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config <command>",
	}

	cmd.AddCommand(show.NewCmd())
	cmd.AddCommand(get.NewCmd())
	cmd.AddCommand(set.NewCmd())

	return cmd
}
