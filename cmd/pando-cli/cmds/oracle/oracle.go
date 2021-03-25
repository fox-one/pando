package oracle

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "oracle",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewEditCmd())
	cmd.AddCommand(NewPokeCmd())
	cmd.AddCommand(NewRelyCmd())
	cmd.AddCommand(NewDenyCmd())

	return cmd
}
