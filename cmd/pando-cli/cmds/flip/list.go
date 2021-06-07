package flip

import (
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var limit = 20

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list auctions",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := call.R(cmd.Context()).Get("/api/flips")
			if err != nil {
				return err
			}

			var body struct {
				Flips json.RawMessage `json:"flips,omitempty"`
			}

			if err := call.UnmarshalResponse(r, &body); err != nil {
				return err
			}

			fields := []string{"id", "collateral_id", "bid", "lot", "tab"}
			lines, err := jq.ParseObjects(body.Flips, fields...)
			if err != nil {
				return err
			}

			cmd.Println(column.Print(lines))
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "page limit")
	return cmd
}
