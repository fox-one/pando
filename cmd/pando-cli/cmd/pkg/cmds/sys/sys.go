package sys

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/sys/withdraw"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sys",
	}

	cmd.AddCommand(withdraw.NewCmd())
	return cmd
}
