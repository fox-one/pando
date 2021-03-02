package node

import (
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/node/system"
	"github.com/fox-one/pando/handler/render"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New(system *core.System) *Server {
	return &Server{system: system}
}

type Server struct {
	system *core.System
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(render.WrapResponse(true))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Get("/info", system.HandleInfo(s.system))

	return r
}
