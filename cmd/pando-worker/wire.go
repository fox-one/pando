//+build wireinject

package main

import (
	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/google/wire"
)

func buildApp(cfg *config.Config) (app, error) {
	wire.Build(
		storeSet,
		serviceSet,
		notifierSet,
		parliamentSet,
		workerSet,
		serverSet,
		wire.Struct(new(app), "*"),
	)

	return app{}, nil
}
