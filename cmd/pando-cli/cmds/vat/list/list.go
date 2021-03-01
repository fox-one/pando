package list

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cfg.GetAuthToken()
			ctx := call.WithToken(cmd.Context(), token)

			r, err := call.RPC().ListVaults(ctx, &api.Req_ListVaults{})
			if err != nil {
				return err
			}

			data, _ := json.Marshal(r.Vaults)

			fields := []string{"id", "collateral_id", "ink", "art"}
			lines, err := jq.ParseObjects(data, fields...)
			if err != nil {
				return err
			}

			cmd.Println(column.Print(lines))
			return nil
		},
	}

	return cmd
}
