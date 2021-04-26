package main

import (
	"fmt"

	"github.com/fox-one/pando/handler/hc"
	"github.com/fox-one/pando/handler/node"
	"github.com/fox-one/pando/server"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/wire"
	"github.com/rs/cors"
)

var serverSet = wire.NewSet(
	node.New,
	provideRoute,
	provideServer,
)

func provideRoute(node *node.Server) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.StripSlashes)
	mux.Use(cors.AllowAll().Handler)
	mux.Use(logger.WithRequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.NewCompressor(5).Handler)

	mux.Mount("/hc", hc.Handle(version))
	mux.Mount("/node", node.Handler())

	return mux
}

func provideServer(mux *chi.Mux) *server.Server {
	return &server.Server{
		Addr:    fmt.Sprintf(":%d", _flag.port),
		Handler: mux,
	}
}
