package use

import (
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use {<host>}",
		Short: "use api host",
		Args:  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			cfg.SetApiHost(host)

			r, err := call.R(cmd.Context()).Get("/api/info")
			if err != nil {
				return err
			}

			var body struct {
				Members       []string `json:"members,omitempty"`
				Threshold     int      `json:"threshold,omitempty"`
				PublicKey     []byte   `json:"public_key,omitempty"`
				OauthClientID string   `json:"oauth_client_id,omitempty"`
			}

			if err := call.UnmarshalResponse(r, &body); err != nil {
				return err
			}

			cfg.SetGroupMembers(body.Members)
			cfg.SetGroupThreshold(body.Threshold)
			cfg.SetGroupVerify(body.PublicKey)
			cfg.SetAuthClient(body.OauthClientID)

			return cfg.Save()
		},
	}
}
