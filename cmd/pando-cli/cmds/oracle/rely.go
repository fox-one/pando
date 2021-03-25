package oracle

import (
	"encoding/base64"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewRelyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "rely",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, publicKey := args[0], args[1]

			pk, err := base64.StdEncoding.DecodeString(publicKey)
			if err != nil {
				return err
			}

			memo, err := actions.Build(
				cmd,
				core.ActionProposalMake,
				core.ActionOracleRely,
				types.UUID(id),
				types.RawMessage(pk),
			)

			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), pay.DefaultAsset, number.One, memo)
		},
	}

	return cmd
}
