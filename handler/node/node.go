package node

import (
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/node/oracle"
	"github.com/fox-one/pando/handler/node/system"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pkg/property"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New(
	system *core.System,
	property property.Store,
	oracles core.OracleStore,
) *Server {
	return &Server{
		system:   system,
		property: property,
		oracles:  oracles,
	}
}

type Server struct {
	system   *core.System
	property property.Store
	oracles  core.OracleStore
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(render.WrapResponse(true))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Get("/info", system.HandleInfo(s.system))
	r.Get("/property", system.HandleProperty(s.property))

	r.Route("/oracles", func(r chi.Router) {
		r.Get("/requests", oracle.HandleScanRequests(s.oracles, s.system))
	})

	return r
}
