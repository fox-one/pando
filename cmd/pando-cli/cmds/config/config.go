package config

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config <command>",
	}

	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewSetCmd())

	return cmd
}
