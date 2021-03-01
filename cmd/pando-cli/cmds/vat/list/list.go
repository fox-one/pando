package list

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := call.R(cmd.Context()).Get("/api/vats")
			if err != nil {
				return err
			}

			var body struct {
				Vaults json.RawMessage `json:"vaults,omitempty"`
			}
			if err := call.UnmarshalResponse(r, &body); err != nil {
				return err
			}

			fields := []string{"id", "collateral_id", "ink", "art"}
			lines, err := jq.ParseObjects(body.Vaults, fields...)
			if err != nil {
				return err
			}

			cmd.Println(column.Print(lines))
			return nil
		},
	}

	return cmd
}
