package main

import (
	"fmt"
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/api"
	"github.com/fox-one/pando/handler/hc"
	"github.com/fox-one/pando/handler/rpc"
	"github.com/fox-one/pando/server"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/wire"
	"github.com/rs/cors"
)

type (
	healthHandler http.Handler
)

var serverSet = wire.NewSet(
	api.New,
	rpc.New,
	provideHealth,
	provideRoute,
	provideServer,
)

func provideRoute(api *api.Server, rpc *rpc.Server, sessions core.Session, hc healthHandler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(cors.AllowAll().Handler)
	r.Use(logger.WithRequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.NewCompressor(5).Handler)

	r.Mount("/twirp", rpc.Handle(sessions))
	r.Mount("/api", api.Handler())
	r.Mount("/hc", hc)

	return r
}

func provideHealth(system *core.System) healthHandler {
	h := hc.Handle(system.Version)
	return healthHandler(h)
}

func provideServer(mux *chi.Mux) *server.Server {
	return &server.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}
}
