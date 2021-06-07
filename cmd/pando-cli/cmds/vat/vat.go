package vat

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vat",
		Short: "manage vaults",
	}

	cmd.AddCommand(NewOpenCmd())
	cmd.AddCommand(NewDepositCmd())
	cmd.AddCommand(NewWithdrawCmd())
	cmd.AddCommand(NewPaybackCmd())
	cmd.AddCommand(NewGenerateCmd())
	cmd.AddCommand(NewListCmd())

	return cmd
}
