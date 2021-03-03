package follow

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "follow",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cfg.GetAuthToken()
			ctx := call.WithToken(cmd.Context(), token)

			id := args[0]
			tx, err := call.RPC().FindTransaction(ctx, &api.Req_FindTransaction{Id: id})
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
