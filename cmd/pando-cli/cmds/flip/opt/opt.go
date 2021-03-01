package opt

import (
	"time"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "opt",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := call.RPC().ReadFlipOption(cmd.Context(), &api.Req_ReadFlipOption{})
			if err != nil {
				return err
			}

			beg := number.Decimal(opt.Beg)
			ttl := time.Duration(opt.Ttl) * time.Second
			tau := time.Duration(opt.Tau) * time.Second

			if len(args) != 3 {
				cmd.Printf("Beg:%s,TTL:%s,Tau:%s\n", beg, ttl.String(), tau.String())
				return nil
			}

			beg = number.Decimal(args[0])
			ttl, _ = time.ParseDuration(args[1])
			tau, _ = time.ParseDuration(args[2])

			memo, err := actions.MakeProposal(
				core.ActionFlipOpt,
				beg,
				int64(ttl.Seconds()),
				int64(tau.Seconds()),
			)
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
