package cat

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := call.R(cmd.Context()).Get("/api/cats")
			if err != nil {
				return err
			}

			var body struct {
				Collaterals json.RawMessage `json:"collaterals,omitempty"`
			}

			if err := call.UnmarshalResponse(r, &body); err != nil {
				return err
			}

			fields := []string{"id", "name", "ink", "debt", "price"}
			for _, arg := range args {
				if !govalidator.IsIn(arg, fields...) {
					fields = append(fields, arg)
				}
			}

			lines, err := jq.ParseObjects(body.Collaterals, fields...)
			if err != nil {
				return err
			}

			cmd.Println(column.Print(lines))
			return nil
		},
	}

	return cmd
}
