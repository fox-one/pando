package flip

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewBidCmd() *cobra.Command {
	var (
		lot string
	)

	cmd := &cobra.Command{
		Use:   "bid <flip id> <bid>",
		Short: "execute Bid action for an auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flipID := args[0]
			flip, err := call.RPC().FindFlip(cmd.Context(), &api.Req_FindFlip{Id: flipID})
			if err != nil {
				return err
			}

			cat, err := call.RPC().FindCollateral(cmd.Context(), &api.Req_FindCollateral{Id: flip.CollateralId})
			if err != nil {
				return err
			}

			bid := number.Decimal(args[1])
			lot := number.Decimal(lot)
			if lot.IsZero() {
				lot = number.Decimal(flip.Lot)
			}

			memo, err := actions.Build(cmd, core.ActionFlipBid, types.UUID(flipID), lot)
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), cat.Dai, bid, memo)
		},
	}

	cmd.Flags().StringVar(&lot, "lot", "0", "gem amount for return")
	return cmd
}
