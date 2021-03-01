package oracle

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/oracle/feed"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "oracle",
	}

	cmd.AddCommand(feed.NewCmd())
	return cmd
}
