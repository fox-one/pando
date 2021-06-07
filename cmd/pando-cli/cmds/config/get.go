package config

import (
	"fmt"

	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {<key>}",
		Short: "get local config by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			_, err := fmt.Fprintln(cmd.OutOrStdout(), cfg.Get(key))
			return err
		},
	}

	return cmd
}
