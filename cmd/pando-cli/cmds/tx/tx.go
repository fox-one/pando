package tx

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tx",
	}

	cmd.AddCommand(NewFollowCmd())
	return cmd
}
