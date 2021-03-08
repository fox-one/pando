package auth

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "auth",
	}

	cmd.AddCommand(NewLoginCmd())

	return cmd
}
