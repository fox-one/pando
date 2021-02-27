package acl

import (
	"net/http"

	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pando/handler/request"
	"github.com/fox-one/pkg/logger"
	"github.com/twitchtv/twirp"
)

func AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := request.UserFrom(r.Context())
		if !ok {
			render.Error(w, twirp.NewError(twirp.Unauthenticated, "authentication required"))
			logger.FromRequest(r).Debugln("api: authentication required")
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
