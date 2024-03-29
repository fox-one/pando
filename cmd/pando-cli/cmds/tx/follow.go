package tx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/cmd/pando-cli/internal/column"
	"github.com/fox-one/pando/cmd/pando-cli/internal/jq"
	api "github.com/fox-one/pando/handler/rpc/pando"
	"github.com/spf13/cobra"
	"github.com/twitchtv/twirp"
)

func NewFollowCmd() *cobra.Command {
	var loop bool

	cmd := &cobra.Command{
		Use:  "follow",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cfg.GetAuthToken()
			ctx := call.WithToken(cmd.Context(), token)

			id := args[0]

		loop:
			tx, err := call.RPC().FindTransaction(ctx, &api.Req_FindTransaction{Id: id})
			if err != nil {
				if terr, ok := err.(twirp.Error); ok && terr.Code() == twirp.NotFound && loop {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(time.Second):
						goto loop
					}
				}

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

	cmd.Flags().BoolVar(&loop, "loop", false, "polling until not 404")
	return cmd
}
