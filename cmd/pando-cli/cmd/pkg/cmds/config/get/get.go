package get

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/cfg"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get {<key>}",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			cmd.Println(cfg.Get(key))
			return nil
		},
	}

	return cmd
}
