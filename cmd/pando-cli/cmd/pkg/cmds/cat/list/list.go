package list

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/jq"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/spf13/cobra"
)

var defaultFields = []string{
	"name",
	"debt",
	"line",
	"price",
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := call.RPC().ListCollaterals(cmd.Context(), &api.Req_ListCollaterals{})
			if err != nil {
				return err
			}

			data, _ := json.Marshal(r.Collaterals)

			fields := []string{"id"}
			if len(args) > 0 {
				fields = append(fields, args...)
			} else {
				fields = append(fields, defaultFields...)
			}

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
