package login

import (
	"fmt"
	"net/url"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "login",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				clientID := cfg.GetAuthClient()
				if clientID == "" {
					return fmt.Errorf("oauth client id not set, run pd use {host} first")
				}

				return requestMixinOauth(clientID)
			}

			r, err := call.R(cmd.Context()).SetBody(map[string]interface{}{
				"code": args[0],
			}).Post("/login")
			if err != nil {
				return err
			}

			var body struct {
				Token string `json:"token,omitempty"`
			}

			if err := call.UnmarshalResponse(r, &body); err != nil {
				return err
			}

			u, err := mixin.UserMe(cmd.Context(), body.Token)
			if err != nil {
				return err
			}

			cmd.Printf("%s welcome!", u.FullName)

			cfg.SetAuthToken(body.Token)
			return cfg.Save()
		},
	}

	return cmd
}

func requestMixinOauth(clientID string) error {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("scope", "PROFILE:READ ASSETS:READ SNAPSHOTS:READ")
	q.Set("response_type", "code")
	q.Set("redirect_url", "https://mixin-oauth.firesbox.com/general-callback")

	u := &url.URL{
		Scheme:   "https",
		Host:     "mixin-oauth.firesbox.com",
		RawQuery: q.Encode(),
	}

	return browser.OpenURL(u.String())
}
