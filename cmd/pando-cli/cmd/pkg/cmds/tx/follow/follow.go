package follow

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/cfg"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "follow",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cfg.GetAuthToken()
			ctx := call.WithToken(cmd.Context(), token)

			follow := args[0]
			tx, err := call.RPC().FindTransaction(ctx, &api.Req_FindTransaction{Follow: follow})
			if err != nil {
				return err
			}

			d := json.NewEncoder(cmd.OutOrStdout())
			d.SetIndent("", "  ")
			return d.Encode(tx)
		},
	}

	return cmd
}
