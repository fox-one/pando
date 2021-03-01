package vat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/deposit"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/generate"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/list"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/open"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/payback"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat/withdraw"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vat",
	}

	cmd.AddCommand(open.NewCmd())
	cmd.AddCommand(deposit.NewCmd())
	cmd.AddCommand(withdraw.NewCmd())
	cmd.AddCommand(payback.NewCmd())
	cmd.AddCommand(generate.NewCmd())
	cmd.AddCommand(list.NewCmd())

	return cmd
}
