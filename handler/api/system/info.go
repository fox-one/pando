package system

import (
	"encoding/base64"
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/render"
)

type InfoResponse struct {
	// oauth client id
	OauthClientID string `json:"oauth_client_id,omitempty" format:"uuid"`
	// multisig members
	Members []string `json:"members,omitempty"`
	// multisig threshold
	Threshold uint8  `json:"threshold,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
}

// ShowSystemInfo godoc
// @Summary Show system info
// @Description
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} InfoResponse
// @Router /info [get]
func HandleInfo(system *core.System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, InfoResponse{
			OauthClientID: system.ClientID,
			Members:       system.Members,
			Threshold:     system.Threshold,
			PublicKey:     base64.StdEncoding.EncodeToString(system.PublicKey),
		})
	}
}
