package auth

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/auth/login"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "auth",
	}

	cmd.AddCommand(login.NewCmd())

	return cmd
}
