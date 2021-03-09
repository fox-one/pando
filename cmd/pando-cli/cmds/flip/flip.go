package flip

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "flip",
	}

	cmd.AddCommand(NewKickCmd())
	cmd.AddCommand(NewBidCmd())
	cmd.AddCommand(NewDealCmd())
	cmd.AddCommand(NewListCmd())

	return cmd
}
