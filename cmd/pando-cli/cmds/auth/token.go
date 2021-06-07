package auth

import (
	"fmt"

	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/spf13/cobra"
)

func NewTokenCmd() *cobra.Command {
	var (
		header bool
	)

	cmd := &cobra.Command{
		Use: "token",
		Short: "show current authorization token",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := cfg.GetAuthToken()
			if header {
				v = fmt.Sprintf(`Authorization:"Bearer %s"`, v)
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), v)
			return nil
		},
	}

	cmd.Flags().BoolVar(&header, "header", false, "as header")
	return cmd
}
