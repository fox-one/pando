package tx

import (
	"encoding/json"
	"fmt"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/spf13/cobra"
)

func NewFollowCmd() *cobra.Command {
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

			b, _ := json.Marshal(tx)
			lines, err := jq.ParseObject(b, "status", "msg", "action", "parameters")
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), column.Print(lines))
			return nil
		},
	}

	return cmd
}
