package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/param"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pkg/logger"
	"github.com/twitchtv/twirp"
)

func HandleOauth(
	userz core.UserService,
	sessions core.Session,
	notifier core.Notifier,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Code string `json:"code,omitempty" valid:"required"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, err)
			return
		}

		ctx := r.Context()
		token, err := userz.Auth(ctx, body.Code)
		if err != nil {
			render.Error(w, twirp.InvalidArgumentError("code", err.Error()))
			return
		}

		user, err := sessions.Login(ctx, token)
		if err != nil {
			render.Error(w, twirp.InvalidArgumentError("token", err.Error()))
			return
		}

		if err := notifier.Auth(ctx, user); err != nil {
			logger.FromContext(ctx).Errorf("api: cannot notify auth")
		}

		render.JSON(w, render.H{
			"id":     user.MixinID,
			"name":   user.Name,
			"avatar": user.Avatar,
			"token":  token,
			"scope":  extractScope(token),
		})
	}
}

func extractScope(token string) string {
	var claim struct {
		jwt.StandardClaims
		Scp string `json:"scp,omitempty"`
	}

	_, _ = jwt.ParseWithClaims(token, &claim, nil)
	return claim.Scp
}
