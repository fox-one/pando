package vat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/deposit"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/generate"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/init"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/list"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/payback"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat/withdraw"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vat",
	}

	cmd.AddCommand(init.NewCmd())
	cmd.AddCommand(deposit.NewCmd())
	cmd.AddCommand(withdraw.NewCmd())
	cmd.AddCommand(payback.NewCmd())
	cmd.AddCommand(generate.NewCmd())
	cmd.AddCommand(list.NewCmd())

	return cmd
}
