package auth

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "auth",
		Short: "Auth & login with auth code",
	}

	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewTokenCmd())

	return cmd
}
