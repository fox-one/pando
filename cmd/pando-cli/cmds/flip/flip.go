package flip

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/flip/bid"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/flip/deal"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/flip/kick"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/flip/list"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/flip/opt"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "flip",
	}

	cmd.AddCommand(kick.NewCmd())
	cmd.AddCommand(bid.NewCmd())
	cmd.AddCommand(deal.NewCmd())
	cmd.AddCommand(opt.NewCmd())
	cmd.AddCommand(list.NewCmd())

	return cmd
}
