package sys

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sys",
	}

	cmd.AddCommand(NewWithdrawCmd())
	return cmd
}
