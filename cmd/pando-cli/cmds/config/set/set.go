package set

import (
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set {<key> | <value>}",
		Args: cobra.ExactValidArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			cfg.Set(key, value)
			return cfg.Save()
		},
	}

	return cmd
}
