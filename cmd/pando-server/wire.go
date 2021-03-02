//+build wireinject

package main

import (
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/server"
	"github.com/google/wire"
)

func buildServer(cfg *config.Config) (*server.Server, error) {
	wire.Build(
		storeSet,
		serviceSet,
		notifierSet,
		sessionSet,
		serverSet,
	)

	return &server.Server{}, nil
}
