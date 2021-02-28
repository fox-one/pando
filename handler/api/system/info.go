package system

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/render"
)

func HandleInfo(system *core.System) http.HandlerFunc {
	verifyKey := system.PrivateKey.Public().(ed25519.PublicKey)
	view := render.H{
		"oauth_client_id": system.ClientID,
		"members":         system.MemberIDs(),
		"threshold":       system.Threshold,
		"public_key":      base64.StdEncoding.EncodeToString(verifyKey),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, view)
	}
}
