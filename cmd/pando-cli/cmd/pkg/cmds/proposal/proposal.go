package proposal

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/proposal/vote"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposal",
		Aliases: []string{"pp"},
	}

	cmd.AddCommand(vote.NewCmd())

	return cmd
}
