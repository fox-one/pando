package flip

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flip",
		Short: "execute flip actions",
	}

	cmd.AddCommand(NewKickCmd())
	cmd.AddCommand(NewBidCmd())
	cmd.AddCommand(NewDealCmd())
	cmd.AddCommand(NewListCmd())

	return cmd
}
