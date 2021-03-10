package user

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/param"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pkg/logger"
	"github.com/twitchtv/twirp"
)

type LoginRequest struct {
	// mixin oauth code
	Code string `json:"code,omitempty" valid:"required"`
}

type LoginResponse struct {
	// user mixin id
	ID string `json:"id,omitempty" format:"uuid"`
	// user name
	Name string `json:"name,omitempty"`
	// user avatar
	Avatar string `json:"avatar,omitempty"`
	// mixin oauth token
	Token string `json:"token,omitempty"`
	// mixin oauth scope
	Scope string `json:"scope,omitempty"`
}

// LoginByCode godoc
// @Summary login with mixin oauth code
// @Description
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body LoginRequest false "request login"
// @Success 200 {object} LoginResponse
// @Router /login [post]
func HandleOauth(
	userz core.UserService,
	sessions core.Session,
	notifier core.Notifier,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body LoginRequest
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

		render.JSON(w, LoginResponse{
			ID:     user.MixinID,
			Name:   user.Name,
			Avatar: user.Avatar,
			Token:  token,
			Scope:  extractScope(token),
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
