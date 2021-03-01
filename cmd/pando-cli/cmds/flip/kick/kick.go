package kick

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "kick <vault id> <bid>",
		Args: cobra.ExactArgs(2),
		RunE: run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	token := cfg.GetAuthToken()

	user, err := mixin.UserMe(cmd.Context(), token)
	if err != nil {
		return err
	}
	// follow id
	follow := uuid.New()

	vatID := args[0]
	bid := number.Decimal(args[1])

	vat, err := call.RPC().FindVault(cmd.Context(), &api.Req_FindVault{Id: vatID})
	if err != nil {
		return err
	}

	cat, err := call.RPC().FindCollateral(cmd.Context(), &api.Req_FindCollateral{Id: vat.CollateralId})
	if err != nil {
		return err
	}

	memo, err := actions.Tx(
		core.ActionFlipKick,
		types.UUID(user.UserID),
		types.UUID(follow),
		types.UUID(vatID),
	)
	if err != nil {
		return err
	}

	cmd.Println("tx follow id:", follow)
	return pay.Request(cmd.Context(), cat.Dai, bid, memo)
}
