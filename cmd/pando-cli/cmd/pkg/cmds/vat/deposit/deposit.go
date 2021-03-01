package deposit

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/cfg"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deposit <vault id> <deposit>",
		Args: cobra.ExactValidArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cfg.GetAuthToken()
			ctx := call.WithToken(cmd.Context(), token)

			user, err := mixin.UserMe(ctx, token)
			if err != nil {
				return err
			}
			// follow id
			follow := uuid.New()

			vatID := args[0]
			vat, err := call.RPC().FindVault(ctx, &api.Req_FindVault{Id: vatID})
			if err != nil {
				return err
			}

			cat, err := call.RPC().FindCollateral(ctx, &api.Req_FindCollateral{Id: vat.CollateralId})
			if err != nil {
				return err
			}

			dink := number.Decimal(args[1])
			memo, err := actions.Tx(
				core.ActionVatDeposit,
				types.UUID(user.UserID),
				types.UUID(follow),
				types.UUID(vatID),
			)
			if err != nil {
				return err
			}

			cmd.Println("tx follow id:", follow)
			return pay.Request(ctx, cat.Gem, dink, memo)
		},
	}

	return cmd
}
