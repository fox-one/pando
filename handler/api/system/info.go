package system

import (
	"encoding/base64"
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/render"
)

func HandleInfo(system *core.System) http.HandlerFunc {
	view := render.H{
		"oauth_client_id": system.ClientID,
		"members":         system.Members,
		"threshold":       system.Threshold,
		"public_key":      base64.StdEncoding.EncodeToString(system.PublicKey),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, view)
	}
}
