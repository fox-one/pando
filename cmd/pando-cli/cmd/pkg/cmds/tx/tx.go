package tx

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/tx/follow"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tx",
	}

	cmd.AddCommand(follow.NewCmd())
	return cmd
}
