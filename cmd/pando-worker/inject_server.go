package main

import (
	"fmt"
	"net/http"

	"github.com/fox-one/pando/server"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/wire"
	"github.com/rs/cors"
)

var serverSet = wire.NewSet(
	provideHandler,
	provideServer,
)

func provideHandler() http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.StripSlashes)
	mux.Use(cors.AllowAll().Handler)
	mux.Use(logger.WithRequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.NewCompressor(5).Handler)

	return mux
}

func provideServer(handler http.Handler) *server.Server {
	return &server.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
	}
}
