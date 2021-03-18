package proposal

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposal",
		Aliases: []string{"pp"},
	}

	cmd.AddCommand(NewVoteCmd())
	cmd.AddCommand(NewShoutCmd())

	return cmd
}
