package cat

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat",
		Short: "Create, supply, list, edit and fold collateral",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewSupplyCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewEditCmd())
	cmd.AddCommand(NewFoldCmd())
	cmd.AddCommand(NewMoveCmd())

	return cmd
}
