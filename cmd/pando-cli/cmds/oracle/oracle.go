package oracle

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "oracle",
	}

	cmd.AddCommand(NewPokeCmd())
	cmd.AddCommand(NewStepCmd())

	return cmd
}
