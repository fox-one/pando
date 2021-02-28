package get

import (
	"fmt"

	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/cfg"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get {<key>}",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			_, err := fmt.Fprintln(cmd.OutOrStdout(), cfg.Get(key))
			return err
		},
	}

	return cmd
}
